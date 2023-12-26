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
	SharedDevices
	Params bool   `help:"Set RPC params by reading from stdin as JSON."`
	Method string `help:"Set RPC method."`
	Object int64  `help:"Set RPC object."`
}

func (c *CmdRPC) Run(ctx *Context) error {
	var params json.RawMessage
	if c.Params {
		err := json.NewDecoder(os.Stdin).Decode(&params)
		if err != nil {
			return err
		}
	}

	db, err := c.useDB(ctx)
	if err != nil {
		return err
	}

	devices, err := c.useDevices(ctx, db)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	for _, device := range devices {
		wg.Add(1)
		go func(device models.DahuaDeviceConn) {
			conn := dahua.NewClient(device.DahuaConn)

			res, err := func() (string, error) {
				res, err := dahuarpc.SendRaw[json.RawMessage](ctx, conn.RPC, dahuarpc.
					New(c.Method).
					Params(params).
					Object(c.Object))
				if err != nil {
					return "", err
				}

				b, err := json.MarshalIndent(res, "", "  ")
				if err != nil {
					return "", err
				}

				return string(b), nil
			}()
			prefix := fmt.Sprintf("id=%d name=%s", device.DahuaDevice.ID, device.Name)
			if err != nil {
				fmt.Println(prefix, err)
			} else {
				fmt.Println(prefix, res)
			}

			conn.RPC.Close(context.Background())
			wg.Done()
		}(device)
	}
	wg.Wait()

	return nil
}
