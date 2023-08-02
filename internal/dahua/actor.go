package dahua

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/ipcmango/internal/core"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua/auth"
)

var ErrActorClosed = fmt.Errorf("dahua actor is closed")

const actorLogoutTimeout = 15 * time.Second

type ActorHandle struct {
	Camera core.DahuaCamera
	rpcC   chan<- actorRPCRequest
	stopC  <-chan *dahua.Conn
	doneC  <-chan struct{}
}

func NewActorHandle(cam core.DahuaCamera) ActorHandle {
	rpcC := make(chan actorRPCRequest)
	stopC := make(chan *dahua.Conn)
	doneC := make(chan struct{})

	go actorStart(cam, rpcC, stopC, doneC)

	return ActorHandle{
		Camera: cam,
		rpcC:   rpcC,
		stopC:  stopC,
		doneC:  doneC,
	}
}

func actorStart(cam core.DahuaCamera, rpcC <-chan actorRPCRequest, stopC chan<- *dahua.Conn, doneC chan<- struct{}) {
	defer close(doneC)

	conn := actorNewConn(cam)

	for {
		select {
		case stopC <- conn:
			return
		case req := <-rpcC:
			rpc, err := actorNewRPC(req.ctx, conn, &cam)

			select {
			case <-req.ctx.Done():
			case req.res <- actorRPCResponse{rpc: rpc, err: err}:
			}
		}
	}
}

// Close stops the actor and cleans up any resources.
func (c ActorHandle) Close(ctx context.Context) {
	select {
	case <-ctx.Done():
	case <-c.doneC:
	case conn := <-c.stopC:
		if conn.State() == dahua.StateLogin {
			ctx, cancel := context.WithTimeout(ctx, actorLogoutTimeout)
			auth.Logout(ctx, conn)
			cancel()
		}
	}
}

func (c ActorHandle) RPC(ctx context.Context) (dahua.RequestBuilder, error) {
	res := make(chan actorRPCResponse)
	select {
	case <-ctx.Done():
		return dahua.RequestBuilder{}, ctx.Err()
	case <-c.doneC:
		return dahua.RequestBuilder{}, ErrActorClosed
	case c.rpcC <- actorRPCRequest{ctx: ctx, res: res}:
		rpc := <-res
		if rpc.err != nil {
			return dahua.RequestBuilder{}, rpc.err
		}

		return rpc.rpc, nil
	}
}

type actorRPCRequest struct {
	ctx context.Context
	res chan<- actorRPCResponse
}

type actorRPCResponse struct {
	rpc dahua.RequestBuilder
	err error
}

func actorNewConn(cam core.DahuaCamera) *dahua.Conn {
	return dahua.NewConn(http.DefaultClient, cam.Address)
}

func actorNewRPC(ctx context.Context, conn *dahua.Conn, cam *core.DahuaCamera) (dahua.RequestBuilder, error) {
	switch conn.State() {
	case dahua.StateLogout:
		err := auth.Login(ctx, conn, cam.Username, cam.Password)
		if err != nil {
			return dahua.RequestBuilder{}, err
		}

		return conn.RPC(ctx)
	case dahua.StateLogin:
		if err := auth.KeepAlive(ctx, conn); err != nil {
			if conn.State() == dahua.StateLogout {
				return actorNewRPC(ctx, conn, cam)
			}

			return dahua.RequestBuilder{}, err
		}

		return conn.RPC(ctx)
	case dahua.StateError:
		return dahua.RequestBuilder{}, conn.Error
	}

	panic("unhandled connection state")
}
