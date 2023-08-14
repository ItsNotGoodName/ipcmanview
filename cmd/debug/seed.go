package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/qes"
	"github.com/rs/zerolog/log"
)

func seed(ctx context.Context, db qes.Querier) {
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
		cam, err := dahua.DB.CameraCreate(ctx, db, cam)
		if err != nil {
			log.Err(err).Msg("Already exists")
		}
	}
}
