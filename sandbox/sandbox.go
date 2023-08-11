package sandbox

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/qes"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

type Sandbox struct {
	*suture.Supervisor
	db qes.Querier
}

func NewSandbox(db qes.Querier) Sandbox {
	super := suture.New("sandbox", suture.Spec{
		EventHook: sutureext.EventHook(),
	})

	return Sandbox{
		Supervisor: super,
		db:         db,
	}
}

func (s Sandbox) Serve(ctx context.Context) error {
	// --------------------- Seed
	username, _ := os.LookupEnv("IPC_USERNAME")
	password, _ := os.LookupEnv("IPC_PASSWORD")
	ips, _ := os.LookupEnv("IPC_IPS")

	for _, ip := range strings.Split(ips, ",") {
		cam := models.DahuaCamera{
			Address:  ip,
			Username: username,
			Password: password,
			Location: models.Location{Location: time.Local},
		}

		// Force create
		cam, err := dahua.DB.CameraCreate(ctx, s.db, cam)
		if err != nil {
			log.Err(err).Msg("Already exists")
		}
	}

	// ----------------------------------------------------------------------------

	dahuaSuper := dahua.NewSupervisor(s.db)
	s.Supervisor.Add(dahuaSuper)

	dahuaScanSuper := dahua.NewScanSupervisor(s.db, dahuaSuper, 5)
	s.Supervisor.Add(dahuaScanSuper)

	go func() {
		cams, err := dahua.DB.CameraList(ctx, s.db)
		if err != nil {
			log.Err(err).Msg("")
			return
		}

		for _, cam := range cams {
			err = dahua.DB.ScanCursorReset(ctx, s.db, cam.ID)
			if err != nil {
				log.Err(err).Msg("")
				return
			}

			scanCam, err := dahua.DB.ScanCursorGet(ctx, s.db, cam.ID)
			if err != nil {
				log.Err(err).Msg("")
				return
			}

			queueTask := dahua.NewScanTaskQuick(scanCam)

			queueTask2, err := dahua.NewScanTaskFull(scanCam)
			if err != nil {
				log.Err(err).Msg("")
			}

			err = dahua.DB.ScanQueueTaskCreate(ctx, s.db, queueTask)
			if err != nil {
				log.Err(err).Msg("")
			}

			err = dahua.DB.ScanQueueTaskCreate(ctx, s.db, queueTask2)
			if err != nil {
				log.Err(err).Msg("")
			}
		}

		dahuaScanSuper.Scan()
		dahuaScanSuper.Scan()
		dahuaScanSuper.Scan()
		dahuaScanSuper.Scan()
		dahuaScanSuper.Scan()
		dahuaScanSuper.Scan()
		dahuaScanSuper.Scan()
		dahuaScanSuper.Scan()
		dahuaScanSuper.Scan()
		dahuaScanSuper.Scan()
		dahuaScanSuper.Scan()
		dahuaScanSuper.Scan()
		dahuaScanSuper.Scan()
	}()

	// ----------------------------------------------------------------------------

	return s.Supervisor.Serve(ctx)
}

var tesMu *sync.Mutex = &sync.Mutex{}

func scan(ctx context.Context, db qes.Querier, cam models.DahuaCamera, super *dahua.ScanSupervisor) {
}

func print(data ...any) {
	if len(data) > 1 && data[1] != nil {
		log.Debug().Err(data[1].(error)).Msg("")
		return
	}
	log.Debug().Any("data", data[0]).Msg("")
}
