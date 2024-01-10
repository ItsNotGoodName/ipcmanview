package files

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

func DahuaFileName(startTime time.Time, id int64, typ string) string {
	return fmt.Sprintf("%s.%d.%s", startTime.UTC().Format("2006-01-02_15-04-05"), id, typ)
}

func DahuaFileIDFromFileName(fileName string) (int64, error) {
	s := strings.Split(fileName, ".")
	if len(s) != 3 {
		return 0, fmt.Errorf("invalid file name: %s", fileName)
	}

	return strconv.ParseInt(s[1], 10, 64)
}

func NewDahuaFileStore(dir string) DahuaFileStore {
	return DahuaFileStore{
		dir: dir,
	}
}

type DahuaFileStore struct {
	dir string
}

func (s DahuaFileStore) FilePath(startTime time.Time, id int64, typ string) string {
	return filepath.Join(s.dir, DahuaFileName(startTime, id, typ))
}

func (s DahuaFileStore) filePath(file models.DahuaFile) string {
	return filepath.Join(s.dir, DahuaFileName(file.StartTime, file.ID, file.Type))
}

func (s DahuaFileStore) Exists(ctx context.Context, file models.DahuaFile) (bool, error) {
	if _, err := os.Stat(s.filePath(file)); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

func (s DahuaFileStore) List() ([]int64, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		id, err := DahuaFileIDFromFileName(info.Name())
		if err != nil {
			continue
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (s DahuaFileStore) Save(ctx context.Context, file models.DahuaFile, r io.Reader) error {
	filePath := s.filePath(file)
	filePathSwap := filePath + ".swap"
	f, err := os.OpenFile(filePathSwap, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return err
	}
	f.Close()

	return os.Rename(filePathSwap, filePath)
}

func (s DahuaFileStore) Get(ctx context.Context, file models.DahuaFile) (io.ReadCloser, error) {
	return os.Open(s.filePath(file))
}

func (s DahuaFileStore) Remove(ctx context.Context, file models.DahuaFile) error {
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

		if err := os.Remove(filepath.Join(s.dir, infos[i].Name())); err != nil {
			return 0, err
		}

		currentSize -= infos[i].Size()
		count += 1
	}

	return count, nil
}
