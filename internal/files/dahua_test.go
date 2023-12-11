package files

import (
	"testing"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/stretchr/testify/assert"
)

func Test_DahuaFileName(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name  string
		args  repo.DahuaFile
		want  string
		want2 DahuaFile
	}{
		{args: repo.DahuaFile{ID: 1, CameraID: 3, StartTime: types.NewTime(now), Type: "jpg"}, want: now.UTC().Format("2006-01-02-15-04-05") + "-3-1.jpg", want2: DahuaFile{ID: 1, CameraID: 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toDahuaFileName(tt.args)
			if !assert.Equal(t, tt.want, got, "") {
				return
			}
			got2, err := fromDahuaFileName(got)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tt.want2, got2, "")
		})
	}
}
