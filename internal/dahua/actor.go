package dahua

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ItsNotGoodName/ipcmango/internal/core"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua/auth"
)

var ErrCameraActorClosed = fmt.Errorf("camera actor closed")

type CameraActor struct {
	ID     int64
	rpcC   chan rpcRequest
	closeC chan *dahua.Conn
	doneC  chan struct{}
}

func NewCameraActor(cam core.DahuaCamera) CameraActor {
	rpcC := make(chan rpcRequest)
	closeC := make(chan *dahua.Conn)
	doneC := make(chan struct{})

	go CameraActorStart(cam, rpcC, closeC, doneC)

	return CameraActor{
		ID:     cam.ID,
		rpcC:   rpcC,
		closeC: closeC,
		doneC:  doneC,
	}
}

func (c CameraActor) RPC(ctx context.Context) (dahua.RequestBuilder, error) {
	res := make(chan rpcResponse)
	select {
	case <-ctx.Done():
		return dahua.RequestBuilder{}, ctx.Err()
	case <-c.doneC:
		return dahua.RequestBuilder{}, ErrCameraActorClosed
	case c.rpcC <- rpcRequest{ctx: ctx, res: res}:
		rpc := <-res
		return rpc.rpc, rpc.err
	}
}

func (c CameraActor) Close(ctx context.Context) {
	select {
	case <-ctx.Done():
	case <-c.doneC:
	case conn := <-c.closeC:
		if conn.State == dahua.StateLogin {
			auth.Logout(ctx, conn)
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

func CameraActorStart(cam core.DahuaCamera, rpcC <-chan rpcRequest, stopC chan<- *dahua.Conn, doneC chan<- struct{}) {
	defer close(doneC)

	conn := dahua.NewConn(http.DefaultClient, dahua.NewCamera(cam.Address))

	for {
		select {
		case stopC <- conn:
			return
		case req := <-rpcC:
			rpc, err := rpc(req.ctx, conn, &cam)

			select {
			case <-req.ctx.Done():
			case req.res <- rpcResponse{rpc: rpc, err: err}:
			}
		}
	}
}

func rpc(ctx context.Context, conn *dahua.Conn, cam *core.DahuaCamera) (dahua.RequestBuilder, error) {
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
				return rpc(ctx, conn, cam)
			}

			return dahua.RequestBuilder{}, err
		}

		return conn.RPC(ctx)
	case dahua.StateError:
		return dahua.RequestBuilder{}, conn.Error
	}

	panic("unhandled connection state")
}
