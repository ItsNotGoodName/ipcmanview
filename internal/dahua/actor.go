package dahua

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ItsNotGoodName/ipcmango/internal/core"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua/auth"
)

type CameraActor struct {
	rpcCC  chan chan rpcPayload
	closeC chan *dahua.Conn
}

func CameraActorNew(ctx context.Context, cam core.DahuaCamera) CameraActor {
	rpcCC := make(chan chan rpcPayload)
	closeC := make(chan *dahua.Conn)

	go CameraActorStart(ctx, cam, rpcCC, closeC)

	return CameraActor{
		rpcCC:  rpcCC,
		closeC: closeC,
	}
}

func (c CameraActor) RPC(ctx context.Context) (dahua.RequestBuilder, error) {
	select {
	case <-ctx.Done():
		return dahua.RequestBuilder{}, ctx.Err()
	case rpcC, ok := <-c.rpcCC:
		if !ok {
			return dahua.RequestBuilder{}, fmt.Errorf("camera actor closed")
		}

		select {
		case <-ctx.Done():
			return dahua.RequestBuilder{}, ctx.Err()
		case rpc, ok := <-rpcC:
			if !ok {
				return dahua.RequestBuilder{}, fmt.Errorf("camera actor did not reply")
			}
			return rpc.rpc, rpc.err
		}
	}
}

func (c CameraActor) Close(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case conn, ok := <-c.closeC:
		if ok && conn.State == dahua.StateLogin {
			auth.Logout(ctx, conn)
		}

		return nil
	}
}

type rpcPayload struct {
	rpc dahua.RequestBuilder
	err error
}

func CameraActorStart(ctx context.Context, cam core.DahuaCamera, rpcC chan chan rpcPayload, closeC chan *dahua.Conn) {
	defer close(rpcC)
	defer close(closeC)

	conn := dahua.NewConn(http.DefaultClient, dahua.NewCamera(cam.Address))

	for {
		rpcClientC := make(chan rpcPayload)

		select {
		case <-ctx.Done():
			return
		case closeC <- conn:
			return
		case rpcC <- rpcClientC:
			rpc, err := rpc(ctx, conn, &cam)
			select {
			case <-ctx.Done():
				close(rpcClientC)
				return
			case rpcClientC <- rpcPayload{rpc: rpc, err: err}:
				close(rpcClientC)
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
