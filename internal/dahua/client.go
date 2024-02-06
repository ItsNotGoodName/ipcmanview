package dahua

import (
	"context"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/ptz"
)

func NewClient(conn Conn) Client {
	rpcHTTPClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	cgiHTTPClient := http.Client{}
	fileHTTPClient := http.Client{}

	clientRPC := dahuarpc.NewClient(rpcHTTPClient, conn.URL, conn.Username, conn.Password)
	clientPTZ := ptz.NewClient(clientRPC)
	clientCGI := dahuacgi.NewClient(cgiHTTPClient, conn.URL, conn.Username, conn.Password)
	clientFile := dahuarpc.NewFileClient(&fileHTTPClient, 10)

	return Client{
		Conn: conn,
		RPC:  clientRPC,
		PTZ:  clientPTZ,
		CGI:  clientCGI,
		File: clientFile,
	}
}

type Client struct {
	Conn Conn
	RPC  dahuarpc.Client
	PTZ  ptz.Client
	CGI  dahuacgi.Client
	File dahuarpc.FileClient
}

func (c Client) Close(ctx context.Context) error {
	c.File.Close()
	return c.RPC.Close(ctx)
}
