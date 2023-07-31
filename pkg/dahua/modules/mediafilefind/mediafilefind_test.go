package mediafilefind

import (
	"testing"
	"time"

	"github.com/ItsNotGoodName/ipcmango/pkg/dahua"
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
				StartTime: dahua.NewTimestamp(startTime, time.Local),
				EndTime:   dahua.NewTimestamp(endTime, time.Local),
				Type:      "jpg",
			},
			second: FindNextFileInfo{
				FilePath:  "/mnt/sd/2023-04-09/001/jpg/07/12/04[M][0@0][0][1].jpg",
				StartTime: dahua.NewTimestamp(startTime, time.Local),
				EndTime:   dahua.NewTimestamp(endTime, time.Local),
				Type:      "jpg",
			},
		},
		{
			first: FindNextFileInfo{
				FilePath:  "/mnt/sd/2023-04-09/001/jpg/07/12/04[M][0@0][0][].jpg",
				StartTime: dahua.NewTimestamp(startTime, time.Local),
				EndTime:   dahua.NewTimestamp(endTime, time.Local),
				Type:      "dav",
			},
			second: FindNextFileInfo{
				FilePath:  "/mnt/sd/2023-04-09/001/jpg/07/12/04[M][0@0][0][].jpg",
				StartTime: dahua.NewTimestamp(startTime, time.Local),
				EndTime:   dahua.NewTimestamp(endTime, time.Local),
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
				StartTime: dahua.NewTimestamp(startTime, time.Local),
				EndTime:   dahua.NewTimestamp(endTime, time.Local),
				Type:      "jpg",
			},
			second: FindNextFileInfo{
				FilePath:  "/mnt/sd/2023-04-09/001/jpg/07/12/04[M][0@0][0][].jpg",
				StartTime: dahua.NewTimestamp(startTime, time.Local),
				EndTime:   dahua.NewTimestamp(endTime, time.Local),
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
