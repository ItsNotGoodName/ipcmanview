package mediafilefind

import (
	"context"
	"errors"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahua"
)

type Stream struct {
	object int64
	count  int
	closed bool
}

func NewStream(ctx context.Context, gen dahua.GenRPC, condtion Condition) (*Stream, error) {
	object, err := Create(ctx, gen)
	if err != nil {
		return nil, err
	}

	var closed bool
	ok, err := FindFile(ctx, gen, object, condtion)
	if err != nil {
		var resErr *dahua.ErrResponse
		if !errors.As(err, &resErr) {
			return nil, err
		}

		if resErr.Type != dahua.ErrResponseTypeNoData {
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

func (s *Stream) Next(ctx context.Context, gen dahua.GenRPC) ([]FindNextFileInfo, error) {
	if s.closed {
		return nil, nil
	}

	files, err := FindNextFile(ctx, gen, s.object, s.count)
	if err != nil {
		s.Close(gen)
		return nil, err
	}

	if files.Infos == nil {
		s.Close(gen)
		return nil, nil
	}

	if files.Found < s.count {
		s.Close(gen)
	}

	return files.Infos, nil
}

func (s *Stream) Close(gen dahua.GenRPC) {
	if s.closed {
		return
	}

	s.closed = true

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := Close(ctx, gen, s.object); err != nil {
		return
	}

	Destroy(ctx, gen, s.object)
}
