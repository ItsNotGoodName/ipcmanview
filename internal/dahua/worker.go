package dahua

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/auth"
	"github.com/ItsNotGoodName/ipcmanview/pkg/qes"
	"github.com/jackc/pgx/v5"
	"github.com/thejerf/suture/v4"
)

var ErrWorkerClosed = fmt.Errorf("worker is closed")
var ErrWorkerRestart = fmt.Errorf("worker is restarting")

type workerRPCRequest struct {
	ctx context.Context
	res chan<- workerRPCResponse
}

type workerRPCResponse struct {
	rpc dahuarpc.RequestBuilder
	err error
}

// WorkerRPC is responsible for genrating RPC requests for the given camera.
type WorkerRPC struct {
	camC     <-chan Camera
	rpcC     <-chan workerRPCRequest
	restartC <-chan struct{}
}

func (w WorkerRPC) String() string {
	return "dahua.WorkerRPC"
}

func NewWorkerRPC(camC <-chan Camera, rpcC <-chan workerRPCRequest, restartC <-chan struct{}) WorkerRPC {
	return WorkerRPC{
		camC:     camC,
		rpcC:     rpcC,
		restartC: restartC,
	}
}

func (w WorkerRPC) Serve(ctx context.Context) error {
	drainChannel(w.restartC)

	for {
		if err := w.serve(ctx); !errors.Is(err, ErrWorkerRestart) {
			return err
		}
	}
}

func (w WorkerRPC) serve(ctx context.Context) error {
	cam, ok := <-w.camC
	if !ok {
		return suture.ErrDoNotRestart
	}

	authConn := auth.NewConn(dahuarpc.NewConn(http.DefaultClient, cam.Address), cam.Username, cam.Password)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		authConn.Logout(ctx)
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-w.restartC:
			// Wait for RPCs that are in progress to complete.
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(1 * time.Second):
			}
			return ErrWorkerRestart
		case req := <-w.rpcC:
			rpc, err := authConn.RPC(req.ctx)
			req.res <- workerRPCResponse{rpc, err}
		}
	}
}

// WorkerCamera fetches camera from the database and caches it.
// The camera is always available through the camera channel until it is closed.
// If there is no camera for the given id, the entire supervisor tree above it is terminated.
type WorkerCamera struct {
	cameraID   int64
	db         qes.Querier
	refetchC   <-chan struct{}
	dependents []chan<- struct{}
	camC       chan<- Camera
}

func (w WorkerCamera) String() string {
	return "dahua.WorkerCamera"
}

func NewWorkerCamera(cameraID int64, db qes.Querier, camC chan<- Camera, refetchC chan struct{}, dependents []chan<- struct{}) WorkerCamera {
	return WorkerCamera{
		cameraID:   cameraID,
		db:         db,
		camC:       camC,
		dependents: dependents,
		refetchC:   refetchC,
	}
}

func (w WorkerCamera) Serve(ctx context.Context) error {
	cam, err := w.fetch(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case w.camC <- cam:
		case <-w.refetchC:
			newCam, err := w.fetch(ctx)
			if err != nil {
				return err
			}

			// Restart only if new camera is different
			if cam.Different(newCam) {
				for i := range w.dependents {
					select {
					case w.dependents[i] <- struct{}{}:
					default:
					}
				}
			}

			cam = newCam
		}
	}
}

func (w WorkerCamera) fetch(ctx context.Context) (Camera, error) {
	cam, err := DB.CameraGet(ctx, w.db, w.cameraID)
	if errors.Is(err, pgx.ErrNoRows) {
		close(w.camC)
		return Camera{}, errors.Join(suture.ErrTerminateSupervisorTree, err)
	}

	return NewCamera(cam), err
}

// WorkerDone closes the done channel when the context is canceled.
type WorkerDone struct {
	cameraID int64
	doneC    chan<- struct{}
}

func (w WorkerDone) String() string {
	return "dahua.WorkerDone"
}

func NewWorkerDone(cameraID int64, doneC chan<- struct{}) WorkerDone {
	return WorkerDone{
		cameraID: cameraID,
		doneC:    doneC,
	}
}

func (w WorkerDone) Serve(ctx context.Context) error {
	<-ctx.Done()

	close(w.doneC)

	return nil
}

// WorkerEvent subscribes to camera events and inserts them into the database.
type WorkerEvent struct {
	cameraID int64
	db       qes.Querier
	camC     <-chan Camera
	restartC <-chan struct{}
}

func (w WorkerEvent) String() string {
	return "dahua.WorkerEvent"
}

func NewWorkerEvent(cameraID int64, db qes.Querier, camC <-chan Camera, restartC <-chan struct{}) WorkerEvent {
	return WorkerEvent{
		cameraID: cameraID,
		db:       db,
		camC:     camC,
		restartC: restartC,
	}
}

func (w WorkerEvent) Serve(ctx context.Context) error {
	drainChannel(w.restartC)

	restartC := make(chan struct{}, 1)
	for {
		restartCtx, cancel := context.WithCancel(ctx)

		go func() {
			// Wait for restart channel then set restart flag before canceling the restart context
			select {
			case <-ctx.Done():
				return
			case <-w.restartC:
				select {
				case <-ctx.Done():
					cancel()
				case restartC <- struct{}{}:
					cancel()
				}
			}
		}()

		// Run and return err if it is not context.Canceled
		err := w.serve(restartCtx)
		if !errors.Is(err, context.Canceled) {
			cancel()
			return err
		}

		// Continue if it is a restart else return error
		select {
		case <-restartC:
		default:
			cancel()
			return err
		}
		cancel()
	}
}

func (w WorkerEvent) serve(ctx context.Context) error {
	cam, ok := <-w.camC
	if !ok {
		return suture.ErrDoNotRestart
	}

	c := dahuacgi.NewConn(cam.Address, cam.Username, cam.Password)

	em, err := dahuacgi.EventManagerGet(ctx, c, 0)
	if err != nil {
		return err
	}
	defer em.Close()

	rd := em.Reader()

	for {
		if err := rd.Poll(); err != nil {
			return err
		}

		evt, err := rd.ReadEvent()
		if err != nil {
			return err
		}

		_, err = DB.CameraEventCreate(ctx, w.db, w.cameraID, evt)
		if err != nil {
			return err
		}
	}
}
