package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	webdahua "github.com/ItsNotGoodName/ipcmanview/internal/web/dahua"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type CmdRPC struct {
	Shared
	ID     int64  `help:"Run on one camera by ID."`
	All    bool   `help:"Run on all cameras."`
	Method string `help:"Set method."`
	Params bool   `help:"Set params by reading from stdin as JSON."`
	Object int64  `help:"Set object."`
	Seq    int    `help:"Set seq."`
}

func (c *CmdRPC) Run(ctx *Context) error {
	var params json.RawMessage
	if c.Params {
		err := json.NewDecoder(os.Stdin).Decode(&params)
		if err != nil {
			return err
		}
	}

	db, err := useDB(c.DBPath)
	if err != nil {
		return err
	}

	var names []string
	var conns []dahua.Conn
	if c.All {
		dbCameras, err := db.ListDahuaCamera(ctx)
		if err != nil {
			return err
		}

		for i, dbCamera := range webdahua.ConvertListDahuaCameraRows(dbCameras) {
			conns = append(conns, dahua.NewConn(dbCamera))
			names = append(names, dbCameras[i].Name)
		}
	} else {
		dbCamera, err := db.GetDahuaCamera(ctx, c.ID)
		if err != nil {
			return err
		}

		conns = append(conns, dahua.NewConn(webdahua.ConvertGetDahuaCameraRow(dbCamera)))
		names = append(names, dbCamera.Name)
	}
	defer func() {
		wg := sync.WaitGroup{}
		for _, c := range conns {
			wg.Add(1)
			go func(c dahua.Conn) {
				c.RPC.Close(ctx)
				wg.Done()
			}(c)
		}
		wg.Wait()
	}()

	wg := sync.WaitGroup{}
	for i, conn := range conns {
		wg.Add(1)
		go func(i int, conn dahua.Conn) {
			res, err := func() (string, error) {
				rpc, err := conn.RPC.RPC(ctx)
				if err != nil {
					return "", err
				}

				res, err := dahuarpc.
					SendRaw[json.RawMessage](ctx, rpc.
					Method(c.Method).
					Params(params).
					Object(c.Object).
					Seq(c.Seq))
				if err != nil {
					return "", err
				}

				b, err := json.MarshalIndent(res, "", "  ")
				if err != nil {
					return "", err
				}

				return string(b), nil
			}()
			prefix := fmt.Sprintf("id=%d, camera=%s", conn.Camera.ID, names[i])
			if err != nil {
				fmt.Println(prefix, err)
			} else {
				fmt.Println(prefix, res)
			}

			wg.Done()
		}(i, conn)
	}
	wg.Wait()

	return nil
}
