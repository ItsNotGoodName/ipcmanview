package core

import (
	"testing"
)

func TestDahuaCameraNew(t *testing.T) {
	type args struct {
		r DahuaCameraCreate
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args:    args{r: DahuaCameraCreate{Address: "localhost:80"}},
			wantErr: false,
		},
		{
			args:    args{r: DahuaCameraCreate{Address: "localhost"}},
			wantErr: false,
		},
		{
			args:    args{r: DahuaCameraCreate{Address: "localhost/"}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDahuaCamera(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("CameraNew() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
