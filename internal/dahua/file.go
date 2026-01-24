package dahua

import (
	"context"
	"io"
	"net"
	"net/url"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/jlaffaye/ftp"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func OpenFileFTP(ctx context.Context, db *sqlx.DB, filePath string) (io.ReadCloser, int64, error) {
	urL, err := url.Parse(filePath)
	if err != nil {
		return nil, 0, err
	}

	var dest StorageDestination
	err = db.GetContext(ctx, &dest, `
		SELECT * FROM dahua_storage_destinations
		WHERE server_address = ? AND storage = ?
	`, urL.Host, StorageFTP)
	if err != nil {
		return nil, 0, err
	}

	c, err := ftp.Dial(core.Address(dest.Server_Address, int(dest.Port)), ftp.DialWithContext(ctx))
	if err != nil {
		return nil, 0, err
	}

	err = c.Login(dest.Username, dest.Password)
	if err != nil {
		return nil, 0, err
	}

	username := "/" + dest.Username
	path, _ := strings.CutPrefix(urL.Path, username)

	contentLength, err := c.FileSize(path)
	if err != nil {
		c.Quit()
		return nil, 0, err
	}

	rd, err := c.Retr(path)
	if err != nil {
		c.Quit()
		return nil, 0, err
	}

	return core.MultiReadCloser{
		Reader:  rd,
		Closers: []func() error{rd.Close, c.Quit},
	}, contentLength, nil
}

func OpenFileSFTP(ctx context.Context, db *sqlx.DB, filePath string) (io.ReadCloser, int64, error) {
	urL, err := url.Parse(filePath)
	if err != nil {
		return nil, 0, err
	}

	var dest StorageDestination
	err = db.GetContext(ctx, &dest, `
		SELECT * FROM dahua_storage_destinations
		WHERE server_address = ? AND storage = ?
	`, urL.Host, StorageFTP)
	if err != nil {
		return nil, 0, err
	}

	conn, err := ssh.Dial("tcp", core.Address(dest.Server_Address, int(dest.Port)), &ssh.ClientConfig{
		User: dest.Username,
		Auth: []ssh.AuthMethod{ssh.Password(dest.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// TODO: check public key
			return nil
		},
	})
	if err != nil {
		return nil, 0, err
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, 0, err
	}

	username := "/" + dest.Username
	path, _ := strings.CutPrefix(urL.Path, username)

	var contentLength int64
	if stat, err := client.Stat(path); err == nil {
		contentLength = stat.Size()
	}

	rd, err := client.Open(path)
	if err != nil {
		client.Close()
		return nil, 0, err
	}

	return core.MultiReadCloser{
		Reader:  rd,
		Closers: []func() error{rd.Close, client.Close},
	}, contentLength, nil
}

func OpenFileLocal(ctx context.Context, client StoreClient, filePath string) (io.ReadCloser, int64, error) {
	v, err := client.File.Do(ctx, dahuarpc.LoadFileURL(client.URL, filePath), dahuarpc.Cookie(client.RPC.Session(ctx)))
	if err != nil {
		return nil, 0, err
	}
	return v, v.ContentLength, nil
}

type File struct {
	ID           string
	Device_ID    int64
	Channel      int
	Start_Time   types.Time
	End_Time     types.Time
	Length       int64
	Type         string
	File_Path    string
	Duration     int64
	Disk         int64
	Video_Stream string
	Flags        types.Slice[string]
	Events       types.Slice[string]
	Cluster      int64
	Partition    int64
	Pic_Index    int64
	Repeat       int64
	Work_Dir     string
	Work_Dir_Sn  bool
	Storage      Storage
	Updated_At   types.Time
}
