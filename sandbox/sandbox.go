package sandbox

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmango/internal/dahua"
	"github.com/ItsNotGoodName/ipcmango/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

func Sandbox(ctx context.Context, pool *pgxpool.Pool) {
	super := suture.NewSimple("root")

	username, _ := os.LookupEnv("IPC_USERNAME")
	password, _ := os.LookupEnv("IPC_PASSWORD")
	ip, _ := os.LookupEnv("IPC_IP")

	cam := models.DahuaCamera{
		Address:  ip,
		Username: username,
		Password: password,
		Location: models.Location{Location: time.Local},
	}

	cam, err := dahua.DB.CameraCreate(ctx, pool, cam)
	if err != nil {
		cam, err = dahua.DB.CameraGetByAddress(ctx, pool, ip)
		if err != nil {
			panic(err)
		}
	}

	worker := dahua.NewWorker(pool, cam.ID)

	super.Add(worker)

	// err = dahua.DB.ScanQueueTaskClear(ctx, pool)
	// if err != nil {
	// 	panic(err)
	// }

	go scan(ctx, pool, worker, cam)

	super.Serve(ctx)
}

var tesMu *sync.Mutex = &sync.Mutex{}

func scan(ctx context.Context, pool *pgxpool.Pool, worker dahua.Worker, cam models.DahuaCamera) {
	worker.Queue(&dahua.TestJob{DB: pool})

	err := dahua.DB.ScanCursorReset(ctx, pool, cam.ID)
	if err != nil {
		log.Err(err).Msg("")
		return
	}

	scanCam, err := dahua.DB.ScanCursorGet(ctx, pool, cam.ID)
	if err != nil {
		log.Err(err).Msg("")
		return
	}

	// queueTask := dahua.NewScanTaskQuick(scanCam)

	queueTask2, err := dahua.NewScanTaskFull(scanCam)
	if err != nil {
		log.Err(err).Msg("")
		return
	}

	// err = dahua.DB.ScanQueueTaskCreate(ctx, pool, queueTask)
	// if err != nil {
	// 	log.Err(err).Msg("")
	// 	return
	// }

	err = dahua.DB.ScanQueueTaskCreate(ctx, pool, queueTask2)
	if err != nil {
		log.Err(err).Msg("")
		return
	}

	for {
		err = dahua.DB.ScanQueueTaskGetAndLock(ctx, pool, scanCam.CameraID, func(ctx context.Context, queueTask models.DahuaScanQueueTask) error {
			return dahua.ScanTaskQueueExecute(ctx, pool, worker, queueTask)
		})
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Err(err).Msg("")
	}
}

func print(data ...any) {
	if len(data) > 1 && data[1] != nil {
		log.Debug().Err(data[1].(error)).Msg("")
		return
	}
	log.Debug().Any("data", data[0]).Msg("")
}
