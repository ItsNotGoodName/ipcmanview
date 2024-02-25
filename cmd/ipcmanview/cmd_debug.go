package main

import "fmt"

type CmdDebug struct {
	Shared
}

func (c *CmdDebug) Run(ctx *Context) error {
	db, err := c.useDB(ctx)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, "CREATE VIRTUAL TABLE IF NOT EXISTS dahua_email_search USING fts5(title, body)")
	if err != nil {
		return err
	}

	res := db.QueryRowContext(ctx, "select sqlite_version()")
	var s any
	if err := res.Scan(&s); err != nil {
		return err
	}
	if err := res.Err(); err != nil {
		return err
	}

	fmt.Println(s)

	return nil
}
