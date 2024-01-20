package main

import (
	"fmt"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/ffmpeg"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pagination"
)

type CmdDebug struct {
	Shared
	SharedDevices
}

func (c *CmdDebug) Run(ctx *Context) error {
	db, err := c.useDB(ctx)
	if err != nil {
		return err
	}

	afs, err := c.useDahuaAFS()
	if err != nil {
		return err
	}

	if err := db.OrphanDeleteDahuaThumbnail(ctx, types.NewTime(time.Now())); err != nil {
		return err
	}

	baseURL := "http://192.168.20.30:8080"

	files, err := db.ListDahuaFile(ctx, repo.ListDahuaFileParams{
		Page: pagination.Page{
			Page:    1,
			PerPage: 100,
		},
		DahuaFileFilter: repo.DahuaFileFilter{},
	})
	if err != nil {
		return err
	}

	for _, file := range files.Data {
		fmt.Println("creating", file.FilePath)

		width, height := 480, 320
		extension := "jpg"
		inputPath := baseURL + api.FileURI(file.DeviceID, file.FilePath)
		videoPosition := 5 * time.Second

		if err := func() error {
			thumbnail, err := dahua.CreateThumbnail(ctx, db, afs,
				dahua.ThumbnailForeignKeys{FileID: file.ID},
				int64(width),
				int64(height),
				dahua.NewAferoFileName(extension))
			if err != nil {
				return nil
			}
			defer thumbnail.Close()

			switch file.Type {
			case models.DahuaFileTypeDAV:
				err := ffmpeg.VideoSnapshot(ctx, inputPath, extension, thumbnail, ffmpeg.VideoSnapshotConfig{
					Width:    int(thumbnail.Width),
					Height:   int(thumbnail.Height),
					Position: videoPosition,
				})
				if err != nil {
					return err
				}
			case models.DahuaFileTypeJPG:
				err := ffmpeg.ImageSnapshot(ctx, inputPath, extension, thumbnail, ffmpeg.ImageSnapshotConfig{
					Width:  int(thumbnail.Width),
					Height: int(thumbnail.Height),
				})
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("invalid file type: %s", file.Type)
			}

			return dahua.ReadyAferoFile(ctx, db, thumbnail.AferoFile.ID, thumbnail.File)
		}(); err != nil {
			return err
		}
	}

	return nil
}
