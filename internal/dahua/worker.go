package dahua

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/auth"
	"github.com/ItsNotGoodName/ipcmanview/pkg/qes"
	"github.com/jackc/pgx/v5"
	"github.com/thejerf/suture/v4"
)

var ErrWorkerClosed = fmt.Errorf("worker is closed")

// Worker is responsible for genrating RPC requests for the given cameraID.
type Worker struct {
	CameraID int64
	db       qes.Querier
	rpcC     chan workerRPCRequest
	restartC chan struct{}
	doneC    <-chan struct{}
}

func NewWorker(db qes.Querier, cameraID int64, doneC <-chan struct{}) Worker {
	return Worker{
		CameraID: cameraID,
		db:       db,
		rpcC:     make(chan workerRPCRequest),
		restartC: make(chan struct{}),
		doneC:    doneC,
	}
}

func (w *Worker) String() string {
	return fmt.Sprintf("dahua.Worker@id-%d", w.CameraID)
}

func (w Worker) Serve(ctx context.Context) error {
	cam, err := DB.CameraGet(ctx, w.db, w.CameraID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Camera was deleted
			return errors.Join(suture.ErrDoNotRestart, err)
		}
		return err
	}

	authConn := auth.NewConn(dahuarpc.NewConn(http.DefaultClient, cam.Address), cam.Username, cam.Password)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		authConn.Logout(ctx)
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-w.restartC:
			return fmt.Errorf("restarting")
		case req := <-w.rpcC:
			rpc, err := authConn.RPC(req.ctx)
			req.res <- workerRPCResponse{rpc, err}
		}
	}
}

func (w Worker) Restart(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-w.doneC:
		return ErrWorkerClosed
	case w.restartC <- struct{}{}:
		return nil
	}
}

func (w Worker) RPC(ctx context.Context) (dahuarpc.RequestBuilder, error) {
	resC := make(chan workerRPCResponse)
	select {
	case <-ctx.Done():
		return dahuarpc.RequestBuilder{}, ctx.Err()
	case <-w.doneC:
		return dahuarpc.RequestBuilder{}, ErrWorkerClosed
	case w.rpcC <- workerRPCRequest{ctx, resC}:
		res := <-resC
		if res.err != nil {
			return dahuarpc.RequestBuilder{}, res.err
		}

		return res.rpc, nil
	}
}

type workerRPCRequest struct {
	ctx context.Context
	res chan<- workerRPCResponse
}

type workerRPCResponse struct {
	rpc dahuarpc.RequestBuilder
	err error
}
