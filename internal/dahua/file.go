package dahua

import (
	"context"
	"io"
	"net"
	"net/url"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"golang.org/x/crypto/ssh"
)

func FileFTPReadCloser(ctx context.Context, db repo.DB, fileFilePath string) (io.ReadCloser, error) {
	u, err := url.Parse(fileFilePath)
	if err != nil {
		return nil, err
	}

	dest, err := db.DahuaGetStorageDestinationByServerAddressAndStorage(ctx, repo.DahuaGetStorageDestinationByServerAddressAndStorageParams{
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

	dest, err := db.DahuaGetStorageDestinationByServerAddressAndStorage(ctx, repo.DahuaGetStorageDestinationByServerAddressAndStorageParams{
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
	return client.File.Do(ctx, dahuarpc.LoadFileURL(client.Conn.URL, filePath), dahuarpc.Cookie(client.RPC.Session(ctx)))
}

// FileLocalDownload downloads file from device and saves it to the afero file system.
func FileLocalDownload(ctx context.Context, db repo.DB, afs afero.Fs, client Client, fileID int64, fileFilePath, fileType string) error {
	rd, err := FileLocalReadCloser(ctx, client, fileFilePath)
	if err != nil {
		return err
	}
	defer rd.Close()

	aferoFile, err := CreateAferoFile(ctx, db, afs, AferoForeignKeys{FileID: fileID}, NewAferoFileName(fileType))
	if err != nil {
		return err
	}
	defer aferoFile.Close()

	log.Info().Int64("device-id", client.Conn.ID).Str("file-path", fileFilePath).Msg("Downloading...")

	if _, err := io.Copy(aferoFile, rd); err != nil {
		return err
	}

	return ReadyAferoFile(ctx, db, aferoFile.ID, aferoFile.File)
}

// func FileLocalDownloadByFilter(ctx context.Context, db repo.DB, afs afero.Fs, store *Store, filter repo.DahuaFileFilter) (int, error) {
// 	filter.Storage = []models.Storage{models.StorageLocal}
//
// 	downloaded := 0
// 	for cursor := ""; ; {
// 		files, err := db.CursorListDahuaFile(ctx, repo.CursorListDahuaFileParams{
// 			PerPage:         100,
// 			Cursor:          cursor,
// 			DahuaFileFilter: filter,
// 		})
// 		if err != nil {
// 			return downloaded, err
// 		}
// 		cursor = files.Cursor
//
// 		for _, file := range files.Data {
// 			aferoFile, err := db.DahuaGetAferoFileByFileID(ctx, sql.NullInt64{Valid: true, Int64: file.ID})
// 			if err != nil {
// 				err = SyncAferoFile(ctx, db, afs, aferoFile.ID, aferoFile.Name)
// 			}
//
// 			if repo.IsNotFound(err) {
// 				// File does not exist
// 			} else if err != nil {
// 				// File error
// 				return downloaded, err
// 			} else {
// 				// File exists
// 				continue
// 			}
//
// 			conn, err := GetConn(ctx, db, file.DeviceID)
// 			if err != nil {
// 				return downloaded, err
// 			}
// 			client := store.Client(ctx, conn)
//
// 			if err := FileLocalDownload(ctx, db, afs, client, file.ID, file.FilePath, file.Type); err != nil {
// 				return downloaded, err
// 			}
//
// 			downloaded++
// 		}
//
// 		if !files.HasMore {
// 			break
// 		}
// 	}
//
// 	return downloaded, nil
// }
