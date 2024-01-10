package main

import (
	"fmt"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
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

	dahuaFileStore, err := c.useDahuaFileStore()
	if err != nil {
		return err
	}

	ids, err := dahuaFileStore.List()
	if err != nil {
		return err
	}

	for _, id := range ids {
		dbFile, err := db.GetDahuaFile(ctx, id)
		if err != nil {
			if repo.IsNotFound(err) {
				continue
			}
			return err
		}
		file := dbFile.Convert()

		if file.Type != models.DahuaFileTypeDAV {
			continue
		}

		exists, err := dahuaFileStore.Exists(ctx, file)
		if err != nil {
			return err
		}
		if !exists {
			continue
		}

		err = dahua.FileDAVToJPG(ctx, dahuaFileStore, file)
		if err != nil {
			return err
		}
	}

	fmt.Println(ids)

	return nil
}
