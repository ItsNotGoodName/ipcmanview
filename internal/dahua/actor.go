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

const logoutTimeout = 15 * time.Second

type Actor struct {
	ID      int64
	rpcC    chan<- rpcRequest
	closeC  <-chan *dahua.Conn
	updateC chan<- core.DahuaCamera
	doneC   <-chan struct{}
}

func StartActor(cam core.DahuaCamera) Actor {
	rpcC := make(chan rpcRequest)
	closeC := make(chan *dahua.Conn)
	updateC := make(chan core.DahuaCamera)
	doneC := make(chan struct{})

	go startActor(cam, rpcC, closeC, updateC, doneC)

	return Actor{
		ID:      cam.ID,
		rpcC:    rpcC,
		closeC:  closeC,
		updateC: updateC,
		doneC:   doneC,
	}
}

func startActor(cam core.DahuaCamera, rpcC <-chan rpcRequest, closeC chan<- *dahua.Conn, updateC <-chan core.DahuaCamera, doneC chan<- struct{}) {
	defer close(doneC)

	conn := newConn(cam)

	for {
		select {
		case closeC <- conn:
			return
		case newCam := <-updateC:
			if newCam.Equal(cam) {
				continue
			}
			cam = newCam

			if conn.State == dahua.StateLogin {
				ctx, cancel := context.WithTimeout(context.Background(), logoutTimeout)
				auth.Logout(ctx, conn)
				cancel()
			}

			conn = newConn(cam)
		case req := <-rpcC:
			rpc, err := newRPC(req.ctx, conn, &cam)

			select {
			case <-req.ctx.Done():
			case req.res <- rpcResponse{rpc: rpc, err: err}:
			}
		}
	}
}

type rpcRequest struct {
	ctx context.Context
	res chan<- rpcResponse
}

type rpcResponse struct {
	rpc dahua.RequestBuilder
	err error
}

func newConn(cam core.DahuaCamera) *dahua.Conn {
	return dahua.NewConn(http.DefaultClient, dahua.NewCamera(cam.Address))
}

func newRPC(ctx context.Context, conn *dahua.Conn, cam *core.DahuaCamera) (dahua.RequestBuilder, error) {
	switch conn.State {
	case dahua.StateLogout:
		err := auth.Login(ctx, conn, cam.Username, cam.Password)
		if err != nil {
			return dahua.RequestBuilder{}, err
		}

		return conn.RPC(ctx)
	case dahua.StateLogin:
		if err := auth.KeepAlive(ctx, conn); err != nil {
			if conn.State == dahua.StateLogout {
				return newRPC(ctx, conn, cam)
			}

			return dahua.RequestBuilder{}, err
		}

		return conn.RPC(ctx)
	case dahua.StateError:
		return dahua.RequestBuilder{}, conn.Error
	}

	panic("unhandled connection state")
}

func (c Actor) Update(ctx context.Context, cam core.DahuaCamera) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.doneC:
		return ErrActorClosed
	case c.updateC <- cam:
		return nil
	}
}

func (c Actor) RPC(ctx context.Context) (dahua.RequestBuilder, error) {
	res := make(chan rpcResponse)
	select {
	case <-ctx.Done():
		return dahua.RequestBuilder{}, ctx.Err()
	case <-c.doneC:
		return dahua.RequestBuilder{}, ErrActorClosed
	case c.rpcC <- rpcRequest{ctx: ctx, res: res}:
		rpc := <-res
		return rpc.rpc, rpc.err
	}
}

func (c Actor) Close(ctx context.Context) {
	select {
	case <-ctx.Done():
	case <-c.doneC:
	case conn := <-c.closeC:
		if conn.State == dahua.StateLogin {
			ctx, cancel := context.WithTimeout(ctx, logoutTimeout)
			auth.Logout(ctx, conn)
			cancel()
		}
	}
}
