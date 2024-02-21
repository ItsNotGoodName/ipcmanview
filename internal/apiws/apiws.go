package apiws

import (
	"context"
	"errors"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Logger(conn *websocket.Conn) zerolog.Logger {
	return log.With().
		Str("package", "apiws").
		Str("remote", conn.RemoteAddr().String()).
		Logger()
}

func Flush(ctx context.Context, vistor Visitor, writeC chan<- []byte) error {
	defer close(writeC)

	data, err := vistor.Visit(ctx)
	if err != nil {
		if errors.Is(err, ErrVisitorEmpty) {
			return nil
		}
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case writeC <- data:
	}

	return nil
}

type Signal struct {
	C chan struct{}
}

func NewSignal() Signal {
	return Signal{
		C: make(chan struct{}, 1),
	}
}

func (s Signal) Queue() {
	select {
	case s.C <- struct{}{}:
	default:
	}
}

var ErrVisitorEmpty = errors.New("visitor empty")

type Visitor interface {
	Visit(ctx context.Context) ([]byte, error)
	HasMore() bool
}

type Vistors struct {
	visitors []Visitor
}

func NewVisitors(visitors ...Visitor) Vistors {
	return Vistors{
		visitors: visitors,
	}
}

func (c Vistors) Visit(ctx context.Context) ([]byte, error) {
	for _, v := range c.visitors {
		data, err := v.Visit(ctx)
		if err != nil {
			if errors.Is(err, ErrVisitorEmpty) {
				continue
			}
			return nil, err
		}

		return data, nil
	}

	return nil, ErrVisitorEmpty
}

func (c Vistors) HasMore() bool {
	for _, v := range c.visitors {
		if v.HasMore() {
			return true
		}
	}
	return false
}

type BufferVisitor struct {
	buffer chan []byte
}

func NewBufferVisitor(count int) *BufferVisitor {
	return &BufferVisitor{
		buffer: make(chan []byte, count),
	}
}

func (v *BufferVisitor) Push(data []byte) bool {
	select {
	case v.buffer <- data:
		return true
	default:
		return false
	}
}

func (v *BufferVisitor) HasMore() bool {
	return len(v.buffer) > 0
}

func (v *BufferVisitor) Visit(ctx context.Context) ([]byte, error) {
	select {
	case event := <-v.buffer:
		return event, nil
	default:
		return nil, ErrVisitorEmpty
	}
}

type OnceVisitor struct {
	Done    bool
	hasMore bool
	data    []byte
}

func NewOnceVisitor() *OnceVisitor {
	return &OnceVisitor{
		data: []byte{},
	}
}

func (v *OnceVisitor) Set(data []byte) {
	v.data = data
	v.hasMore = true
	v.Done = true
}

func (v *OnceVisitor) HasMore() bool {
	return v.hasMore
}

func (v *OnceVisitor) Visit(ctx context.Context) ([]byte, error) {
	if !v.hasMore {
		return nil, ErrVisitorEmpty
	}
	v.hasMore = false
	return v.data, nil
}

func Check(visitor Visitor, sig Signal) {
	if visitor.HasMore() {
		sig.Queue()
	}
}
