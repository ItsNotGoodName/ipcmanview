package files

import (
	"testing"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_DahuaFileName(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name  string
		args  models.DahuaFile
		want  string
		want2 DahuaFile
	}{
		{args: models.DahuaFile{ID: 1, DeviceID: 3, StartTime: now, Type: "jpg"}, want: now.UTC().Format("2006-01-02-15-04-05") + "-3-1.jpg", want2: DahuaFile{ID: 1, DeviceID: 3}},
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
