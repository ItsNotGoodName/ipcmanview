package mediafilefind

import (
	"reflect"
	"testing"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/stretchr/testify/assert"
)

func TestFindNextFileInfo_UniqueTime(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(time.Duration(5))

	type fields struct {
		first  FindNextFileInfo
		second FindNextFileInfo
	}

	notEqual := []fields{
		{
			first: FindNextFileInfo{
				FilePath:  "/mnt/sd/2023-04-09/001/jpg/07/12/04[M][0@0][0][].jpg",
				StartTime: dahuarpc.NewTimestamp(startTime, time.Local),
				EndTime:   dahuarpc.NewTimestamp(endTime, time.Local),
				Type:      "jpg",
			},
			second: FindNextFileInfo{
				FilePath:  "/mnt/sd/2023-04-09/001/jpg/07/12/04[M][0@0][0][1].jpg",
				StartTime: dahuarpc.NewTimestamp(startTime, time.Local),
				EndTime:   dahuarpc.NewTimestamp(endTime, time.Local),
				Type:      "jpg",
			},
		},
		{
			first: FindNextFileInfo{
				FilePath:  "/mnt/sd/2023-04-09/001/jpg/07/12/04[M][0@0][0][].jpg",
				StartTime: dahuarpc.NewTimestamp(startTime, time.Local),
				EndTime:   dahuarpc.NewTimestamp(endTime, time.Local),
				Type:      "dav",
			},
			second: FindNextFileInfo{
				FilePath:  "/mnt/sd/2023-04-09/001/jpg/07/12/04[M][0@0][0][].jpg",
				StartTime: dahuarpc.NewTimestamp(startTime, time.Local),
				EndTime:   dahuarpc.NewTimestamp(endTime, time.Local),
				Type:      "jpg",
			},
		},
	}

	for _, field := range notEqual {
		firstStartTime, firstEndTime, err := field.first.UniqueTime(0, time.Local)
		assert.NoError(t, err)

		secondStartTime, secondEndTime, err := field.second.UniqueTime(0, time.Local)
		assert.NoError(t, err)

		assert.NotEqual(t, firstStartTime, secondStartTime)
		assert.NotEqual(t, firstEndTime, secondEndTime)
	}

	equal := []fields{
		{
			first: FindNextFileInfo{
				FilePath:  "/mnt/sd/2023-04-09/001/jpg/07/12/04[M][0@0][0][].jpg",
				StartTime: dahuarpc.NewTimestamp(startTime, time.Local),
				EndTime:   dahuarpc.NewTimestamp(endTime, time.Local),
				Type:      "jpg",
			},
			second: FindNextFileInfo{
				FilePath:  "/mnt/sd/2023-04-09/001/jpg/07/12/04[M][0@0][0][].jpg",
				StartTime: dahuarpc.NewTimestamp(startTime, time.Local),
				EndTime:   dahuarpc.NewTimestamp(endTime, time.Local),
				Type:      "jpg",
			},
		},
	}

	for _, field := range equal {
		{

			firstStartTime, firstEndTime, err := field.first.UniqueTime(0, time.Local)
			assert.NoError(t, err)

			secondStartTime, secondEndTime, err := field.second.UniqueTime(0, time.Local)
			assert.NoError(t, err)

			assert.Equal(t, firstStartTime, secondStartTime)
			assert.Equal(t, firstEndTime, secondEndTime)
		}

		{
			firstStartTime, firstEndTime, err := field.first.UniqueTime(0, time.Local)
			assert.NoError(t, err)

			// Seed
			secondStartTime, secondEndTime, err := field.second.UniqueTime(1, time.Local)
			assert.NoError(t, err)

			assert.NotEqual(t, firstStartTime, secondStartTime)
			assert.NotEqual(t, firstEndTime, secondEndTime)
		}
	}
}

func TestFindNextFileInfo_CleanEvents(t *testing.T) {
	type fields struct {
		Channel     int
		StartTime   dahuarpc.Timestamp
		EndTime     dahuarpc.Timestamp
		Length      int
		Type        string
		FilePath    string
		Duration    int
		Disk        int
		VideoStream string
		Flags       []string
		Events      []string
		Cluster     int
		Partition   int
		PicIndex    int
		Repeat      int
		WorkDir     string
		WorkDirSN   int
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			fields: fields{
				Events: []string{"CrossRegionDetection", ""},
			},
			want: []string{"CrossRegionDetection"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FindNextFileInfo{
				Channel:     tt.fields.Channel,
				StartTime:   tt.fields.StartTime,
				EndTime:     tt.fields.EndTime,
				Length:      tt.fields.Length,
				Type:        tt.fields.Type,
				FilePath:    tt.fields.FilePath,
				Duration:    tt.fields.Duration,
				Disk:        tt.fields.Disk,
				VideoStream: tt.fields.VideoStream,
				Flags:       tt.fields.Flags,
				Events:      tt.fields.Events,
				Cluster:     tt.fields.Cluster,
				Partition:   tt.fields.Partition,
				PicIndex:    tt.fields.PicIndex,
				Repeat:      tt.fields.Repeat,
				WorkDir:     tt.fields.WorkDir,
				WorkDirSN:   tt.fields.WorkDirSN,
			}
			if got := f.CleanEvents(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindNextFileInfo.CleanEvents() = %v, want %v", got, tt.want)
			}
		})
	}
}
