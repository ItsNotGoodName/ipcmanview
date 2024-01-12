package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/ffmpeg"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pagination"
	"github.com/google/uuid"
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

		thumbnail, err := db.UpsertDahuaFileThumbnail(ctx, repo.UpsertDahuaFileThumbnailParams{
			FileID: file.ID,
			Width:  int64(width),
			Height: int64(height),
		})
		if err != nil {
			return err
		}

		aferoFile, err := db.CreateDahuaAferoFile(ctx, repo.CreateDahuaAferoFileParams{
			FileThumbnailID: sql.NullInt64{
				Int64: thumbnail.ID,
				Valid: true,
			},
			Name:      uuid.NewString() + "." + extension,
			CreatedAt: types.NewTime(time.Now()),
		})
		if err != nil {
			return err
		}

		f, err := fs.Create(aferoFile.Name)
		if err != nil {
			return err
		}

		switch file.Type {
		case models.DahuaFileTypeDAV:
			if err := ffmpeg.VideoSnapshot(ctx, inputPath, extension, f, ffmpeg.VideoSnapshotConfig{
				Width:    int(thumbnail.Width),
				Height:   int(thumbnail.Height),
				Position: 5 * time.Second,
			}); err != nil {
				return errors.Join(fs.Remove(f.Name()), err)
			}
		case models.DahuaFileTypeJPG:
			if err := ffmpeg.ImageSnapshot(ctx, inputPath, extension, f, ffmpeg.ImageSnapshotConfig{
				Width:  int(thumbnail.Width),
				Height: int(thumbnail.Height),
			}); err != nil {
				return errors.Join(fs.Remove(f.Name()), err)
			}
		}
	}

	return nil
}
