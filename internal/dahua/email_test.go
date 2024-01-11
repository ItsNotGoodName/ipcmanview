package dahua

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEmailContent(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want EmailContent
	}{
		{
			name: "",
			args: args{`Alarm Event: Intrusion
		Alarm Input Channel: 1
		Alarm Start Time(D/M/Y H:M:S): 10/01/2024 15:48:47
		Alarm Device Name: CAM-5
		Alarm Name: Front Entrance
		IP Address: 192.168.60.15`},
			want: EmailContent{
				AlarmEvent:        "Intrusion",
				AlarmInputChannel: 1,
				AlarmDeviceName:   "CAM-5",
				AlarmName:         "Front Entrance",
				IPAddress:         "192.168.60.15",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseEmailContent(tt.args.text)
			assert.Equal(t, tt.want, got)
		})
	}
}
