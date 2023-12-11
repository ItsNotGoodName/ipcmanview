package files

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

type DahuaFile struct {
	ID       int64
	CameraID int64
}

func fromDahuaFileName(fileName string) (DahuaFile, error) {
	dotSplit := strings.Split(fileName, ".")
	if len(dotSplit) != 2 {
		return DahuaFile{}, fmt.Errorf("invalid length: %d", len(dotSplit))
	}

	dashSplit := strings.Split(dotSplit[0], "-")
	dashSplitLen := len(dashSplit)
	if dashSplitLen < 2 {
		return DahuaFile{}, fmt.Errorf("invalid length: %d", dashSplitLen)
	}

	cameraID, err := strconv.ParseInt(dashSplit[dashSplitLen-2], 10, 64)
	if err != nil {
		return DahuaFile{}, err
	}

	id, err := strconv.ParseInt(dashSplit[dashSplitLen-1], 10, 64)
	if err != nil {
		return DahuaFile{}, err
	}

	return DahuaFile{
		ID:       id,
		CameraID: cameraID,
	}, nil
}

func toDahuaFileName(file repo.DahuaFile) string {
	return fmt.Sprintf("%s-%d-%d.%s", file.StartTime.UTC().Format("2006-01-02-15-04-05"), file.CameraID, file.ID, file.Type)
}

type DahuaFileStore struct {
	dir string
}

func NewDahuaFileStore(dir string) DahuaFileStore {
	return DahuaFileStore{
		dir: dir,
	}
}

func (s DahuaFileStore) filePath(file repo.DahuaFile) string {
	return path.Join(s.dir, toDahuaFileName(file))
}

func (s DahuaFileStore) Exists(ctx context.Context, file repo.DahuaFile) (bool, error) {
	if _, err := os.Stat(s.filePath(file)); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

func (s DahuaFileStore) Save(ctx context.Context, file repo.DahuaFile, r io.ReadCloser) error {
	filePath := s.filePath(file)
	filePathSwap := filePath + ".swap"
	f, err := os.OpenFile(filePathSwap, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, r); err != nil {
		return err
	}

	return os.Rename(filePathSwap, filePath)
}

func (s DahuaFileStore) Get(ctx context.Context, file repo.DahuaFile) (io.ReadCloser, error) {
	return os.Open(s.filePath(file))
}

func (s DahuaFileStore) Remove(ctx context.Context, file repo.DahuaFile) error {
	err := os.Remove(s.filePath(file))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

func (s DahuaFileStore) Size(ctx context.Context) (int64, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return 0, err
	}

	dirSize := int64(0)
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return 0, err
		}
		if info.Mode().IsRegular() {
			dirSize += info.Size()
		}
	}

	return dirSize, nil
}

func (s DahuaFileStore) Trim(ctx context.Context, size int64, minAge time.Time) (int, error) {
	currentSize, err := s.Size(ctx)
	if err != nil {
		return 0, err
	}

	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return 0, err
	}
	infos := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return 0, err
		}
		infos = append(infos, info)
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].ModTime().Before(infos[j].ModTime())
	})

	var count int
	for i := range infos {
		if currentSize < size {
			break
		}
		if infos[i].ModTime().After(minAge) {
			continue
		}

		if err := os.Remove(path.Join(s.dir, infos[i].Name())); err != nil {
			return 0, err
		}

		currentSize -= infos[i].Size()
		count += 1
	}

	return count, nil
}
