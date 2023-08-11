package sandbox

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

func Sandbox(ctx context.Context, pool *pgxpool.Pool) {
	// --------------------- Seed
	username, _ := os.LookupEnv("IPC_USERNAME")
	password, _ := os.LookupEnv("IPC_PASSWORD")
	ip, _ := os.LookupEnv("IPC_IP")

	cam := models.DahuaCamera{
		Address:  ip,
		Username: username,
		Password: password,
		Location: models.Location{Location: time.Local},
	}

	// Force create
	cam, err := dahua.DB.CameraCreate(ctx, pool, cam)
	if err != nil {
		cam, err = dahua.DB.CameraGetByAddress(ctx, pool, ip)
		if err != nil {
			panic(err)
		}
	}

	// ----------------------------------------------------------------------------

	super := suture.New("root", suture.Spec{
		EventHook: sutureext.EventHook(),
	})

	dahuaSuper := dahua.NewSupervisor(pool, &event.Bus{})
	super.Add(dahuaSuper)
	dahuaScanSuper := dahua.NewScanSupervisor(pool, dahuaSuper, 5)
	super.Add(dahuaScanSuper)

	go func() {
		scan(ctx, pool, cam, dahuaScanSuper)
	}()

	super.Serve(ctx)
}

var tesMu *sync.Mutex = &sync.Mutex{}

func scan(ctx context.Context, pool *pgxpool.Pool, cam models.DahuaCamera, super *dahua.ScanSupervisor) {
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

	queueTask := dahua.NewScanTaskQuick(scanCam)

	queueTask2, err := dahua.NewScanTaskFull(scanCam)
	if err != nil {
		log.Err(err).Msg("")
	}

	err = dahua.DB.ScanQueueTaskCreate(ctx, pool, queueTask)
	if err != nil {
		log.Err(err).Msg("")
	}

	err = dahua.DB.ScanQueueTaskCreate(ctx, pool, queueTask2)
	if err != nil {
		log.Err(err).Msg("")
	}

	super.Scan()
	super.Scan()
	super.Scan()
	super.Scan()
	super.Scan()
	super.Scan()
	super.Scan()
	super.Scan()
	super.Scan()
	super.Scan()

	// for {
	// 	err = dahua.DB.ScanQueueTaskNext(ctx, pool, func(ctx context.Context, scanCursorLock dahua.ScanCursorLock, queueTask models.DahuaScanQueueTask) error {
	// 		return dahua.ScanTaskQueueExecute(ctx, pool, worker, scanCursorLock, queueTask)
	// 	})
	// 	time.Sleep(5 * time.Second)
	// }

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
