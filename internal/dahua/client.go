package dahua

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/ptz"
)

func clientDialTimeout(duration time.Duration) func(network, addr string) (net.Conn, error) {
	return func(network, addr string) (net.Conn, error) {
		return net.DialTimeout(network, addr, duration)
	}
}

func NewClient(conn Conn) Client {
	httpClient := http.Client{
		Transport: &http.Transport{
			Dial: clientDialTimeout(5 * time.Second),
		},
	}

	clientRPC := dahuarpc.NewClient(&httpClient, conn.URL, conn.Username, conn.Password)
	clientPTZ := ptz.NewClient(clientRPC)
	clientCGI := dahuacgi.NewClient(httpClient, conn.URL, conn.Username, conn.Password)
	clientFile := dahuarpc.NewFileClient(&httpClient, 10)

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
