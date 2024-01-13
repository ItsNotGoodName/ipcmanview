package dahua

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net"
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

func FileFTPReadCloser(ctx context.Context, db repo.DB, fileFilePath string) (io.ReadCloser, error) {
	u, err := url.Parse(fileFilePath)
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

func FileSFTPReadCloser(ctx context.Context, db repo.DB, fileFilePath string) (io.ReadCloser, error) {
	u, err := url.Parse(fileFilePath)
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
	return client.File.Do(ctx,
		dahuarpc.LoadFileURL(client.Conn.Url, filePath),
		dahuarpc.Cookie(client.RPC.Session(ctx)))
}

// FileLocalDownload downloads file from device and saves it to the afero file system.
func FileLocalDownload(ctx context.Context, db repo.DB, afs afero.Fs, client Client, fileID int64, fileFilePath, fileType string) error {
	aferoFile, err := db.CreateDahuaAferoFile(ctx, repo.CreateDahuaAferoFileParams{
		FileID: sql.NullInt64{
			Int64: fileID,
			Valid: true,
		},
		Name:      uuid.NewString() + "." + fileType,
		CreatedAt: types.NewTime(time.Now()),
	})
	if err != nil {
		return err
	}

	rd, err := FileLocalReadCloser(ctx, client, fileFilePath)
	if err != nil {
		return err
	}
	defer rd.Close()

	wt, err := afs.Create(aferoFile.Name)
	if err != nil {
		return err
	}
	defer wt.Close()

	log.Info().Int64("device-id", client.Conn.ID).Str("file-path", fileFilePath).Msg("Downloading...")

	if _, err := io.Copy(wt, rd); err != nil {
		return err
	}

	return nil
}

func FileLocalDownloadByFilter(ctx context.Context, db repo.DB, afs afero.Fs, store *Store, filter repo.DahuaFileFilter) (int, error) {
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
			aferoFileExists, err := SyncAferoFile(ctx, db, afs, aferoFile, err)
			if err != nil {
				return downloaded, err
			} else if aferoFileExists {
				continue
			}

			device, err := db.GetDahuaDevice(ctx, file.DeviceID)
			if err != nil {
				return downloaded, err
			}
			client := store.Client(ctx, device.Convert().DahuaConn)

			if err := FileLocalDownload(ctx, db, afs, client, file.ID, file.FilePath, file.Type); err != nil {
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
	ids  []int64
	resC chan<- fileServiceRes
}

type fileServiceRes struct {
	downloaded int
	err        error
}

type fileServiceFilterReq struct {
	filter repo.DahuaFileFilter
	resC   chan<- fileServiceFilterRes
}

type fileServiceFilterRes struct {
	downloaded int
	err        error
}

func NewFileService(db repo.DB, fs afero.Fs, store *Store) FileService {
	return FileService{
		db:        db,
		fs:        fs,
		store:     store,
		req:       make(chan fileServiceReq),
		filterReq: make(chan fileServiceFilterReq),
	}
}

// FileService handles downloading files.
type FileService struct {
	db        repo.DB
	fs        afero.Fs
	store     *Store
	req       chan fileServiceReq
	filterReq chan fileServiceFilterReq
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
		case req := <-s.filterReq:
			downloaded, err := FileLocalDownloadByFilter(ctx, s.db, s.fs, s.store, req.filter)
			req.resC <- fileServiceFilterRes{
				downloaded: downloaded,
				err:        err,
			}
		case req := <-s.req:
			downloaded := 0
			var firstErr error
			for _, id := range req.ids {
				err := s.download(ctx, id)
				if err != nil {
					if firstErr == nil {
						firstErr = err
					}
					log.Err(err).Msg("Failed to download file")
				} else {
					downloaded++
				}
			}

			req.resC <- fileServiceRes{
				downloaded: downloaded,
				err:        firstErr,
			}
		}
	}
}

func (s FileService) download(ctx context.Context, id int64) error {
	file, err := s.db.GetDahuaFile(ctx, id)
	if err != nil {
		return err
	}

	device, err := s.db.GetDahuaDevice(ctx, file.DeviceID)
	if err != nil {
		return err
	}
	client := s.store.Client(ctx, device.Convert().DahuaConn)

	return FileLocalDownload(ctx, s.db, s.fs, client, file.ID, file.FilePath, file.Type)
}

func (s FileService) DownloadByFilter(ctx context.Context, filter repo.DahuaFileFilter) (int, error) {
	resC := make(chan fileServiceFilterRes, 1)
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case s.filterReq <- fileServiceFilterReq{
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

func (s FileService) Download(ctx context.Context, fileIDs ...int64) (int, error) {
	resC := make(chan fileServiceRes, 1)
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case s.req <- fileServiceReq{
		ids:  fileIDs,
		resC: resC,
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
