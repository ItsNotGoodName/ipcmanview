package dahua

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/google/uuid"
	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"golang.org/x/crypto/ssh"
)

var ErrFileServiceConflict = fmt.Errorf("file service conflict")

const FileEchoRoute = "/v1/dahua/:id/files/*"

func FileURI(deviceID int64, filePath string) string {
	return fmt.Sprintf("/v1/dahua/%d/files/%s", deviceID, filePath)
}

func FileFTPReadCloser(ctx context.Context, db repo.DB, dahuaFile models.DahuaFile) (io.ReadCloser, error) {
	u, err := url.Parse(dahuaFile.FilePath)
	if err != nil {
		return nil, err
	}

	dest, err := db.GetDahuaStorageDestinationByServerAddressAndStorage(ctx, repo.GetDahuaStorageDestinationByServerAddressAndStorageParams{
		ServerAddress: u.Hostname(),
		Storage:       models.StorageFTP,
	})
	if err != nil {
		return nil, err
	}

	c, err := ftp.Dial(core.Address(dest.ServerAddress, int(dest.Port)), ftp.DialWithContext(ctx))
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

	dest, err := db.GetDahuaStorageDestinationByServerAddressAndStorage(ctx, repo.GetDahuaStorageDestinationByServerAddressAndStorageParams{
		ServerAddress: u.Hostname(),
		Storage:       models.StorageSFTP,
	})
	if err != nil {
		return nil, err
	}

	conn, err := ssh.Dial("tcp", core.Address(dest.ServerAddress, int(dest.Port)), &ssh.ClientConfig{
		User: dest.Username,
		Auth: []ssh.AuthMethod{ssh.Password(dest.Password)},
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

	username := "/" + dest.Username
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, dahuarpc.LoadFileURL(client.Conn.Url, filePath), nil)
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

// FileLocalDownload downloads file from device and saves it to the afero file system.
func FileLocalDownload(ctx context.Context, db repo.DB, fs afero.Fs, client Client, filePath, fileName string) error {
	rd, err := FileLocalReadCloser(ctx, client, filePath)
	if err != nil {
		return err
	}
	defer rd.Close()

	wt, err := fs.Create(fileName)
	if err != nil {
		return err
	}
	defer wt.Close()

	log.Info().Int64("device-id", client.Conn.ID).Str("file-path", filePath).Msg("Downloading...")

	if _, err := io.Copy(wt, rd); err != nil {
		return err
	}

	return nil
}

func FileLocalDownloadByFilter(ctx context.Context, db repo.DB, fs afero.Fs, store *Store, filter repo.DahuaFileFilter) (int, error) {
	filter.Storage = []models.Storage{models.StorageLocal}

	downloaded := 0
	for cursor := ""; ; {
		files, err := db.CursorListDahuaFile(ctx, repo.CursorListDahuaFileParams{
			PerPage:         100,
			Cursor:          cursor,
			DahuaFileFilter: filter,
		})
		if err != nil {
			return downloaded, err
		}
		cursor = files.Cursor

		for _, file := range files.Data {
			aferoFile, err := db.GetDahuaAferoFileByFileID(ctx, sql.NullInt64{Valid: true, Int64: file.ID})
			aferoFileExists, err := SyncAferoFile(ctx, db, fs, aferoFile, err)
			if err != nil {
				return downloaded, err
			} else if aferoFileExists {
				continue
			}

			aferoFile, err = db.CreateDahuaAferoFile(ctx, repo.CreateDahuaAferoFileParams{
				FileID: sql.NullInt64{
					Int64: file.ID,
					Valid: true,
				},
				Name:      uuid.NewString() + "." + file.Type,
				CreatedAt: types.NewTime(time.Now()),
			})
			if err != nil {
				return downloaded, err
			}

			device, err := db.GetDahuaDevice(ctx, file.DeviceID)
			if err != nil {
				return downloaded, err
			}
			client := store.Client(ctx, device.Convert().DahuaConn)

			if err := FileLocalDownload(ctx, db, fs, client, file.FilePath, aferoFile.Name); err != nil {
				return downloaded, err
			}

			downloaded++
		}

		if !files.HasMore {
			break
		}
	}

	return downloaded, nil
}

type fileServiceReq struct {
	filter repo.DahuaFileFilter
	resC   chan<- fileServiceRes
}

type fileServiceRes struct {
	downloaded int
	err        error
}

func NewFileService(db repo.DB, fs afero.Fs, store *Store) FileService {
	return FileService{
		db:       db,
		fs:       fs,
		store:    store,
		requestC: make(chan fileServiceReq),
	}
}

// FileService handles downloading files.
type FileService struct {
	db       repo.DB
	fs       afero.Fs
	store    *Store
	requestC chan fileServiceReq
}

func (s FileService) String() string {
	return "dahua.FileService"
}

func (s FileService) Serve(ctx context.Context) error {
	return sutureext.SanitizeError(ctx, s.serve(ctx))
}

func (s FileService) serve(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case request := <-s.requestC:
			downloaded, err := FileLocalDownloadByFilter(ctx, s.db, s.fs, s.store, request.filter)
			request.resC <- fileServiceRes{
				downloaded: downloaded,
				err:        err,
			}
		}
	}
}

func (s FileService) Run(ctx context.Context, filter repo.DahuaFileFilter) (int, error) {
	resC := make(chan fileServiceRes, 1)
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case s.requestC <- fileServiceReq{
		filter: filter,
		resC:   resC,
	}:
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case res := <-resC:
			return res.downloaded, res.err
		}
	default:
		return 0, ErrFileServiceConflict
	}
}
