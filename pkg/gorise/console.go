package gorise

import (
	"context"
	"fmt"
)

func BuildConsole(cfg Config) (Sender, error) {
	return NewConsole(), nil
}

func NewConsole() Console {
	return Console{}
}

type Console struct {
}

func (Console) Send(ctx context.Context, msg Message) error {
	fmt.Println(msg.Text())
	return nil
}
