package mediafilefind

import (
	"context"
	"strconv"
	"time"

	"github.com/ItsNotGoodName/ipcmango/pkg/dahua"
)

func Create(ctx context.Context, gen dahua.GenRPC) (int64, error) {
	rpc, err := gen.RPC(ctx)
	if err != nil {
		return 0, err
	}

	res, err := dahua.Send[any](ctx, rpc.Method("mediaFileFind.factory.create"))

	return res.Result.Integer(), err
}

func FindFile(ctx context.Context, gen dahua.GenRPC, object int64, condition Condition) (bool, error) {
	rpc, err := gen.RPC(ctx)
	if err != nil {
		return false, err
	}

	res, err := dahua.Send[any](ctx, rpc.
		Method("mediaFileFind.findFile").
		Params(dahua.JSON{"condition": condition}).
		Object(object))

	return res.Result.Bool(), err
}

type Condition struct {
	Channel   int             `json:"Channel"`
	Dirs      []string        `json:"Dirs"`
	Types     []string        `json:"Types"`
	Order     ConditionOrder  `json:"Order"`
	Redundant string          `json:"Redundant"`
	Events    []string        `json:"Events"`
	StartTime dahua.Timestamp `json:"StartTime"`
	EndTime   dahua.Timestamp `json:"EndTime"`
	Flags     []string        `json:"Flags"`
}

type ConditionOrder = string

const (
	ConditionOrderAscent  ConditionOrder = "Ascent"
	ConditionOrderDescent ConditionOrder = "Descent"
)

func NewCondtion(startTime dahua.Timestamp, endTime dahua.Timestamp) Condition {
	return Condition{
		Channel:   0,
		Dirs:      []string{},
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

func FindNextFile(ctx context.Context, gen dahua.GenRPC, object int64, count int) (FindNextFileResult, error) {
	rpc, err := gen.RPC(ctx)
	if err != nil {
		return FindNextFileResult{}, err
	}

	res, err := dahua.Send[FindNextFileResult](ctx, rpc.
		Method("mediaFileFind.findNextFile").
		Params(dahua.JSON{"count": count}).
		Object(object))

	return res.Params, err
}

type FindNextFileResult struct {
	Found int                `json:"found"`
	Infos []FindNextFileInfo `json:"infos"`
}

type FindNextFileInfo struct {
	Channel     int             `json:"Channel"`
	StartTime   dahua.Timestamp `json:"StartTime"`
	EndTime     dahua.Timestamp `json:"EndTime"`
	Length      int             `json:"Length"`
	Type        string          `json:"Type"`
	FilePath    string          `json:"FilePath"`
	Duration    int             `json:"Duration"`
	Disk        int             `json:"Disk"`
	VideoStream string          `json:"VideoStream"`
	Flags       []string        `json:"Flags"`
	Events      []string        `json:"Events"`
	Cluster     int             `json:"Cluster"`
	Partition   int             `json:"Partition"`
	PicIndex    int             `json:"PicIndex"`
	Repeat      int             `json:"Repeat"`
	WorkDir     string          `json:"WorkDir"`
	WorkDirSN   int             `json:"WorkDirSN"`
}

// UniqueTime returns StartTime and EndTime that are unique.
//
//	Dahua cameras can only handle timestamps that are precise up to the second, the leftover microseconds (.000_000) are used to create a unique time.
//	Unique means it won't conflict with other media files on the camera.
//	A seed can optionally be passed to make the media file not conflict with other cameras.
func (f FindNextFileInfo) UniqueTime(affixSeed int, cameraLocation *time.Location) (time.Time, time.Time, error) {
	startTime, err := f.StartTime.Parse(cameraLocation)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endTime, err := f.EndTime.Parse(cameraLocation)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	var prefixSeed int

	{
		tags := dahua.ExtractFilePathTags(f.FilePath)
		if len(tags) >= 4 {
			tag1Seed, _ := strconv.Atoi(tags[2])
			tag2Seed, _ := strconv.Atoi(tags[3])
			prefixSeed += tag1Seed + tag2Seed
		}
	}

	for _, c := range f.Type {
		prefixSeed += int(c)
	}

	seed := (time.Duration((prefixSeed % 999)) * time.Millisecond) + (time.Duration((affixSeed % 999)) * time.Microsecond)

	return startTime.Add(seed), endTime.Add(seed), nil
}

var defaultTagSeed = []string{"0", "0"}

func GetCount(ctx context.Context, gen dahua.GenRPC, object int64) (int, error) {
	rpc, err := gen.RPC(ctx)
	if err != nil {
		return 0, err
	}

	res, err := dahua.Send[struct {
		Count int `json:"count"`
	}](ctx, rpc.
		Method("mediaFileFind.getCount").
		Object(object))

	return res.Params.Count, err
}

func Close(ctx context.Context, gen dahua.GenRPC, object int64) (bool, error) {
	rpc, err := gen.RPC(ctx)
	if err != nil {
		return false, err
	}

	res, err := dahua.Send[any](ctx, rpc.
		Method("mediaFileFind.close").
		Object(object))

	return res.Result.Bool(), err
}

func Destroy(ctx context.Context, gen dahua.GenRPC, object int64) (bool, error) {
	rpc, err := gen.RPC(ctx)
	if err != nil {
		return false, err
	}

	res, err := dahua.Send[any](ctx, rpc.
		Method("mediaFileFind.destroy").
		Object(object))

	return res.Result.Bool(), err
}
