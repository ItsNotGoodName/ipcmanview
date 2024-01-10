package dahua

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/files"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func FileDAVToJPG(ctx context.Context, fileStore files.DahuaFileStore, file models.DahuaFile) error {
	if file.Type != models.DahuaFileTypeDAV {
		return fmt.Errorf("invalid type: %s", file.Type)
	}

	inputFilePath := fileStore.FilePath(file.StartTime, file.ID, file.Type)
	outputFilePath := fileStore.FilePath(file.StartTime, file.ID, models.DahuaFileTypeJPG)

	// ffmpeg -n -i file:2024-01-08T04:25:01Z.115614.dav -ss 00:00:06.000 -vframes 1 output.jpg
	output, err := exec.Command("ffmpeg", "-n", "-i", "file:"+inputFilePath, "-ss", "00:00:06.000", "-vframes", "1", outputFilePath).CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return err
	}

	return nil
}

func FileURL(urL string, deviceID int64, filePath string) string {
	return fmt.Sprintf("%s/v1/dahua/%d/files/%s", urL, deviceID, filePath)
}

func FileFTPReadCloser(ctx context.Context, db repo.DB, dahuaFile models.DahuaFile) (io.ReadCloser, error) {
	storage := core.StorageFromFilePath(dahuaFile.FilePath)

	u, err := url.Parse(dahuaFile.FilePath)
	if err != nil {
		return nil, err
	}

	dbDest, err := db.GetDahuaStorageDestinationByServerAddressAndStorage(ctx, repo.GetDahuaStorageDestinationByServerAddressAndStorageParams{
		ServerAddress: u.Hostname(),
		Storage:       storage,
	})
	if err != nil {
		return nil, err
	}
	dest := dbDest.Convert()

	c, err := ftp.Dial(StorageDestinationURL(dest), ftp.DialWithContext(ctx))
	if err != nil {
		return nil, err
	}

	err = c.Login(dest.Username, dest.Password)
	if err != nil {
		return nil, err
	}

	username := "/" + dest.Username
	path, _ := strings.CutPrefix(u.Path, username)

	rd, err := c.Retr(path)
	if err != nil {
		c.Quit()
		return nil, err
	}

	return core.MultiReadCloser{
		Reader:  rd,
		Closers: []func() error{rd.Close, c.Quit},
	}, nil
}

func FileSFTPReadCloser(ctx context.Context, db repo.DB, dahuaFile models.DahuaFile) (io.ReadCloser, error) {
	u, err := url.Parse(dahuaFile.FilePath)
	if err != nil {
		return nil, err
	}

	dbCred, err := db.GetDahuaStorageDestinationByServerAddressAndStorage(ctx, repo.GetDahuaStorageDestinationByServerAddressAndStorageParams{
		ServerAddress: u.Hostname(),
		Storage:       models.StorageSFTP,
	})
	if err != nil {
		return nil, err
	}
	cred := dbCred.Convert()

	conn, err := ssh.Dial("tcp", StorageDestinationURL(cred), &ssh.ClientConfig{
		User: cred.Username,
		Auth: []ssh.AuthMethod{ssh.Password(cred.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// TODO: check public key
			return nil
		},
	})
	if err != nil {
		return nil, err
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}

	username := "/" + cred.Username
	path, _ := strings.CutPrefix(u.Path, username)

	rd, err := client.Open(path)
	if err != nil {
		client.Close()
		return nil, err
	}

	return core.MultiReadCloser{
		Reader:  rd,
		Closers: []func() error{rd.Close, client.Close},
	}, nil
}

func FileLocalReadCloser(ctx context.Context, client Client, filePath string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, dahuarpc.LoadFileURL(client.Conn.Address, filePath), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Cookie", dahuarpc.Cookie(client.RPC.Session(ctx)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
