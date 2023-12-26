package core

import (
	"reflect"
	"testing"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

func TestStorageFromFilePath(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want models.Storage
	}{
		{name: "", args: args{filePath: "/some/file.jpg"}, want: models.StorageLocal},
		{name: "", args: args{filePath: "sftp://some/file.jpg"}, want: models.StorageSFTP},
		{name: "", args: args{filePath: "ftp://some/file.jpg"}, want: models.StorageFTP},
		{name: "", args: args{filePath: "nfs://some/file.jpg"}, want: models.StorageNFS},
		{name: "", args: args{filePath: "smb://some/file.jpg"}, want: models.StorageSMB},
		{name: "", args: args{filePath: "../some/file.jpg"}, want: models.StorageLocal},
		{name: "", args: args{filePath: "unknow://some/file.jpg"}, want: models.StorageLocal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StorageFromFilePath(tt.args.filePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StorageFromFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
