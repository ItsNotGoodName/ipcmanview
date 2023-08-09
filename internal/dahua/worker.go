package dahua

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/ipcmango/pkg/dahua"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua/auth"
	"github.com/ItsNotGoodName/ipcmango/pkg/qes"
	"github.com/jackc/pgx/v5"
	"github.com/thejerf/suture/v4"
)

var ErrWorkerClosed = fmt.Errorf("dahua worker is closed")

type Worker struct {
	cameraID  int64
	db        qes.Querier
	rpcC      chan workerRPCRequest
	jobC      chan WorkerJob
	doneC     chan struct{}
	restartC  chan struct{}
	shutdownC chan struct{}
}

func NewWorker(db qes.Querier, cameraID int64) Worker {
	return Worker{
		cameraID:  cameraID,
		db:        db,
		rpcC:      make(chan workerRPCRequest),
		jobC:      make(chan WorkerJob, 1),
		doneC:     make(chan struct{}),
		restartC:  make(chan struct{}),
		shutdownC: make(chan struct{}),
	}
}

func (w Worker) Serve(ctx context.Context) error {
	cam, err := DB.CameraGet(ctx, w.db, w.cameraID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Die when camera does not exist
			return errors.Join(suture.ErrDoNotRestart, err)
		}
		return err
	}

	authConn := auth.NewConn(dahua.NewConn(http.DefaultClient, cam.Address), cam.Username, cam.Password)
	defer func() {
		// Do not leave open connections on the camera
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		authConn.Logout(ctx)
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			// Safe to close channel because we will not start up again
			close(w.doneC)
			return nil
		case <-w.shutdownC:
			return suture.ErrDoNotRestart
		case <-w.restartC:
			return fmt.Errorf("restarting")
		case req := <-w.rpcC:
			rpc, err := authConn.RPC(req.ctx)
			req.res <- workerRPCResponse{rpc, err}
		case job := <-w.jobC:
			if err := job.Execute(ctx, &authConn, cam); err != nil {
				return err
			}
		}
	}
}

func (w Worker) Shutdown(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-w.doneC:
		return nil
	case w.restartC <- struct{}{}:
		return nil
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

func (w Worker) Queue(job WorkerJob) {
	for {
		select {
		case <-w.jobC:
		case w.jobC <- job:
			return
		}
	}
}

func (w Worker) RPC(ctx context.Context) (dahua.RequestBuilder, error) {
	resC := make(chan workerRPCResponse)
	select {
	case <-ctx.Done():
		return dahua.RequestBuilder{}, ctx.Err()
	case <-w.doneC:
		return dahua.RequestBuilder{}, ErrWorkerClosed
	case w.rpcC <- workerRPCRequest{ctx, resC}:
		res := <-resC
		if res.err != nil {
			return dahua.RequestBuilder{}, res.err
		}

		return res.rpc, nil
	}
}

type workerRPCRequest struct {
	ctx context.Context
	res chan<- workerRPCResponse
}

type workerRPCResponse struct {
	rpc dahua.RequestBuilder
	err error
}
