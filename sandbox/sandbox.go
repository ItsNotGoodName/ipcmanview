package sandbox

import (
	"context"
	"os"
	"time"

	"github.com/ItsNotGoodName/ipcmango/internal/core"
	"github.com/ItsNotGoodName/ipcmango/internal/dahua"
	"github.com/ItsNotGoodName/ipcmango/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Sandbox(ctx context.Context, pool *pgxpool.Pool) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return
	}
	defer conn.Release()

	username, _ := os.LookupEnv("IPC_USERNAME")
	password, _ := os.LookupEnv("IPC_PASSWORD")
	ip, _ := os.LookupEnv("IPC_IP")

	cam := core.DahuaCamera{
		Address:  ip,
		Username: username,
		Password: password,
	}

	cam, err = dahua.DB.CameraCreate(ctx, conn, cam)
	if err != nil {
		print(cam, err)
		cam, err = dahua.DB.CameraGetByAddress(ctx, conn, ip)
		print(cam, err)
	}

	c := dahua.NewActorHandle(cam)
	defer c.Close(ctx)

	print(dahua.Scan(ctx, conn, c, models.DahuaScanCamera{
		ID:       cam.ID,
		Location: time.Local,
	}, dahua.ScanPeriod{
		Start: time.Now().Add(-24 * time.Hour * 30),
		End:   time.Now(),
	}))
}
