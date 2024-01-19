package mediafilefind

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

func Create(ctx context.Context, c dahuarpc.Conn) (int64, error) {
	res, err := dahuarpc.Send[any](ctx, c, dahuarpc.
		New("mediaFileFind.factory.create"))

	return res.Result.Integer(), err
}

func FindFile(ctx context.Context, c dahuarpc.Conn, object int64, condition Condition) (bool, error) {
	res, err := dahuarpc.Send[any](ctx, c, dahuarpc.
		New("mediaFileFind.findFile").
		Params(struct {
			Condition Condition `json:"condition"`
		}{
			Condition: condition,
		}).
		Object(object))

	return res.Result.Bool(), err
}

type Condition struct {
	Channel   int                `json:"Channel"`
	Dirs      []string           `json:"Dirs"`
	Types     []string           `json:"Types"`
	Order     ConditionOrder     `json:"Order"`
	Redundant string             `json:"Redundant"`
	Events    []string           `json:"Events"`
	StartTime dahuarpc.Timestamp `json:"StartTime"`
	EndTime   dahuarpc.Timestamp `json:"EndTime"`
	Flags     []string           `json:"Flags"`
}

type ConditionOrder = string

const (
	ConditionOrderAscent  ConditionOrder = "Ascent"
	ConditionOrderDescent ConditionOrder = "Descent"
)

func NewCondtion(startTime dahuarpc.Timestamp, endTime dahuarpc.Timestamp) Condition {
	return Condition{
		Channel:   0,
		Dirs:      nil,
		Types:     []string{"dav", "jpg"},
		Order:     ConditionOrderAscent,
		Redundant: "Exclusion",
		Events:    nil,
		StartTime: startTime,
		EndTime:   endTime,
		Flags:     []string{"Timing", "Event", "Event", "Manual"},
	}
}

func (c Condition) Video() Condition {
	c.Types = []string{"dav"}
	return c
}

func (c Condition) Picture() Condition {
	c.Types = []string{"jpg"}
	c.Flags = []string{"Timing", "Event", "Event"}
	return c
}

func FindNextFile(ctx context.Context, c dahuarpc.Conn, object int64, count int) (FindNextFileResult, error) {
	res, err := dahuarpc.Send[FindNextFileResult](ctx, c, dahuarpc.
		New("mediaFileFind.findNextFile").
		Params(struct {
			Count int `json:"count"`
		}{
			Count: count,
		}).
		Object(object))

	return res.Params, err
}

type FindNextFileResult struct {
	Found int                `json:"found"`
	Infos []FindNextFileInfo `json:"infos"`
}

type FindNextFileInfo struct {
	Channel     int                `json:"Channel"`
	StartTime   dahuarpc.Timestamp `json:"StartTime"`
	EndTime     dahuarpc.Timestamp `json:"EndTime"`
	Length      int                `json:"Length"`
	Type        string             `json:"Type"`
	FilePath    string             `json:"FilePath"`
	Duration    int                `json:"Duration"`
	Disk        int                `json:"Disk"`
	VideoStream string             `json:"VideoStream"`
	Flags       []string           `json:"Flags"`
	Events      []string           `json:"Events"`
	Cluster     int                `json:"Cluster"`
	Partition   int                `json:"Partition"`
	PicIndex    int                `json:"PicIndex"`
	Repeat      int                `json:"Repeat"`
	// WorkDir is the working directory (e.g. /mnt/dvr/mmc0p2_0).
	WorkDir string `json:"WorkDir"`
	// WorkDirSN is indicates that the WorkDir's name is the Serial Number (e.g ftp://192.168.20.30/test/share/XXXXXXXXXXXXXXX).
	WorkDirSN int `json:"WorkDirSN"`
}

// UniqueTime returns start and end times that do not conflict with other files.
//
// Dahua devices can only handle timestamps that are precise up to the second (e.g. "2006-01-02 15:04:05"), the leftover microseconds (.000_000) are used create a unique time.
// A seed can optionally be passed to not conflict with other devices.
func (f FindNextFileInfo) UniqueTime(seed int, deviceLocation *time.Location) (time.Time, time.Time, error) {
	startTime, err := f.StartTime.Parse(deviceLocation)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endTime, err := f.EndTime.Parse(deviceLocation)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	var prefixSeed int

	{
		tags := dahuarpc.ExtractFilePathTags(f.FilePath)
		if len(tags) >= 4 {
			tag1Seed, _ := strconv.Atoi(tags[2])
			tag2Seed, _ := strconv.Atoi(tags[3])
			prefixSeed += tag1Seed + tag2Seed
		}
	}

	for _, c := range f.Type {
		prefixSeed += int(c)
	}

	finalSeed := (time.Duration((prefixSeed % 999)) * time.Millisecond) + (time.Duration((seed % 999)) * time.Microsecond)

	return startTime.Add(finalSeed), endTime.Add(finalSeed), nil
}

// Local checks if the file is stored directly on disk which allows it to be loaded through RPC_Loadfile.
func (f FindNextFileInfo) Local() bool {
	return strings.HasPrefix(f.FilePath, "/")
}

func GetCount(ctx context.Context, c dahuarpc.Conn, object int64) (int, error) {
	res, err := dahuarpc.Send[struct {
		Count int `json:"count"`
	}](ctx, c, dahuarpc.
		New("mediaFileFind.getCount").
		Object(object))

	return res.Params.Count, err
}

func Close(ctx context.Context, c dahuarpc.Conn, object int64) (bool, error) {
	res, err := dahuarpc.Send[any](ctx, c, dahuarpc.
		New("mediaFileFind.close").
		Object(object))

	return res.Result.Bool(), err
}

func Destroy(ctx context.Context, c dahuarpc.Conn, object int64) (bool, error) {
	res, err := dahuarpc.Send[any](ctx, c, dahuarpc.
		New("mediaFileFind.destroy").
		Object(object))

	return res.Result.Bool(), err
}
