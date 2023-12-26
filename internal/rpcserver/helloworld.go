package rpcserver

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type HelloWorld struct {
}

func (s *HelloWorld) Hello(ctx context.Context, req *rpc.HelloReq) (*rpc.HelloResp, error) {
	return &rpc.HelloResp{Text: "Hello " + req.Subject, CurrentTime: timestamppb.New(req.CurrentTime.AsTime().Add(-24 * time.Hour))}, nil
}
