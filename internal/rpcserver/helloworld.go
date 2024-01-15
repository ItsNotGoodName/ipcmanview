package rpcserver

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type HelloWorld struct {
}

var i int32 = 0

func (s *HelloWorld) Hello(ctx context.Context, req *rpc.HelloReq) (*rpc.HelloResp, error) {
	time.Sleep(250 * time.Millisecond)

	if f := atomic.AddInt32(&i, 1); f%5 == 0 {
		return nil, fmt.Errorf("random error: %d", f)
	}

	return &rpc.HelloResp{Text: "Hello " + req.Subject, CurrentTime: timestamppb.New(req.CurrentTime.AsTime().Add(-24 * time.Hour))}, nil
}
