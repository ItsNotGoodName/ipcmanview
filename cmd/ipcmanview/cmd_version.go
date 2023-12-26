package main

import (
	"fmt"

	"github.com/ItsNotGoodName/ipcmanview/internal/build"
)

type CmdVersion struct {
}

func (c *CmdVersion) Run(ctx *Context) error {
	fmt.Println(build.Current.Version)
	return nil
}
