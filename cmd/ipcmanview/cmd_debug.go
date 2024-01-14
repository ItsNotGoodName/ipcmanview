package main

import (
	"fmt"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/ffmpeg"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
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

	if err := db.OrphanDeleteDahuaThumbnail(ctx); err != nil {
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
		inputPath := baseURL + dahua.FileURI(file.DeviceID, file.FilePath)
		videoPosition := 5 * time.Second

		if err := func() error {
			thumbnail, err := dahua.CreateThumbnail(ctx, db, dahua.ThumbnailForeignKeys{FileID: file.ID}, int64(width), int64(height))
			if err != nil {
				return nil
			}

			aferoFile, err := dahua.CreateAferoFile(ctx, db, afs, dahua.AferoForeignKeys{ThumbnailID: thumbnail.ID}, dahua.NewAferoFileName(extension))
			if err != nil {
				return err
			}
			defer aferoFile.Close()

			switch file.Type {
			case models.DahuaFileTypeDAV:
				err := ffmpeg.VideoSnapshot(ctx, inputPath, extension, aferoFile, ffmpeg.VideoSnapshotConfig{
					Width:    int(thumbnail.Width),
					Height:   int(thumbnail.Height),
					Position: videoPosition,
				})
				if err != nil {
					return err
				}
			case models.DahuaFileTypeJPG:
				err := ffmpeg.ImageSnapshot(ctx, inputPath, extension, aferoFile, ffmpeg.ImageSnapshotConfig{
					Width:  int(thumbnail.Width),
					Height: int(thumbnail.Height),
				})
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("invalid file type: %s", file.Type)
			}

			return dahua.ReadyAferoFile(ctx, db, aferoFile.ID, aferoFile.File)
		}(); err != nil {
			return err
		}
	}

	return nil
}
