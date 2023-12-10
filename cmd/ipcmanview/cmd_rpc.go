package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type CmdRPC struct {
	Shared
	SharedCameras
	Method string `help:"Set RPC method."`
	Params bool   `help:"Set RPC params by reading from stdin as JSON."`
	Object int64  `help:"Set RPC object."`
	Seq    int    `help:"Set RPC seq."`
}

func (c *CmdRPC) Run(ctx *Context) error {
	var params json.RawMessage
	if c.Params {
		err := json.NewDecoder(os.Stdin).Decode(&params)
		if err != nil {
			return err
		}
	}

	db, err := useDB(ctx, c.DBPath)
	if err != nil {
		return err
	}

	cameras, err := c.useCameras(ctx, db)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	for _, camera := range cameras {
		wg.Add(1)
		go func(camera models.DahuaCameraInfo) {
			conn := dahua.NewConn(camera.DahuaCamera)

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
			prefix := fmt.Sprintf("id=%d name=%s", camera.ID, camera.Name)
			if err != nil {
				fmt.Println(prefix, err)
			} else {
				fmt.Println(prefix, res)
			}

			conn.RPC.Close(context.Background())
			wg.Done()
		}(camera)
	}
	wg.Wait()

	return nil
}
