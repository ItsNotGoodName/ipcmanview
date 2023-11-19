package mediafilefind

import (
	"context"
	"errors"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type Stream struct {
	object int64
	count  int
	closed bool
}

func NewStream(ctx context.Context, c dahuarpc.Client, condtion Condition) (*Stream, error) {
	object, err := Create(ctx, c)
	if err != nil {
		return nil, err
	}

	var closed bool
	ok, err := FindFile(ctx, c, object, condtion)
	if err != nil {
		var resErr *dahuarpc.ResponseError
		if !errors.As(err, &resErr) {
			return nil, err
		}

		if resErr.Type != dahuarpc.ErrResponseTypeNoData {
			return nil, err
		}

		closed = true
	} else {
		closed = !ok
	}

	return &Stream{
		object: object,
		count:  64,
		closed: closed,
	}, nil
}

func (s *Stream) Next(ctx context.Context, c dahuarpc.Client) ([]FindNextFileInfo, error) {
	if s.closed {
		return nil, nil
	}

	files, err := FindNextFile(ctx, c, s.object, s.count)
	if err != nil {
		s.Close(c)
		return nil, err
	}

	if files.Infos == nil {
		s.Close(c)
		return nil, nil
	}

	if files.Found < s.count {
		s.Close(c)
	}

	return files.Infos, nil
}

func (s *Stream) Close(c dahuarpc.Client) {
	if s.closed {
		return
	}

	s.closed = true

	// TODO: find another way to close stream when context was canceled.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := Close(ctx, c, s.object); err != nil {
		return
	}

	Destroy(ctx, c, s.object)
}
