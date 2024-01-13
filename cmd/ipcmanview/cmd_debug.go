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

	fs, err := c.useDahuaFileFS()
	if err != nil {
		return err
	}

	baseURL := "http://192.168.20.30:8080"

	files, err := db.ListDahuaFile(ctx, repo.ListDahuaFileParams{
		Page: pagination.Page{
			Page:    1,
			PerPage: 100,
		},
		DahuaFileFilter: repo.DahuaFileFilter{Type: []string{"dav"}},
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
			thumbnail, err := dahua.UpsertFileThumbnail(ctx, db, file.ID, int64(width), int64(height))
			if err != nil {
				return nil
			}

			aferoFile, err := dahua.AferoCreateFileThumbnail(ctx, db, fs, thumbnail.ID, dahua.NewAferoFileName(extension))
			if err != nil {
				return err
			}
			defer aferoFile.Close()

			switch file.Type {
			case models.DahuaFileTypeDAV:
				return ffmpeg.VideoSnapshot(ctx, inputPath, extension, aferoFile, ffmpeg.VideoSnapshotConfig{
					Width:    int(thumbnail.Width),
					Height:   int(thumbnail.Height),
					Position: videoPosition,
				})
			case models.DahuaFileTypeJPG:
				return ffmpeg.ImageSnapshot(ctx, inputPath, extension, aferoFile, ffmpeg.ImageSnapshotConfig{
					Width:  int(thumbnail.Width),
					Height: int(thumbnail.Height),
				})
			default:
				return fmt.Errorf("invalid file type: %s", file.Type)
			}
		}(); err != nil {
			return err
		}
	}

	return nil
}
