package dahuacgi

import (
	_ "embed"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var eventReaderFixture string = `--myboundary
Content-Type: text/plain
Content-Length: 147

Code=VideoMotion;action=Start;index=0;data={
	"LocaleTime":	"2023-08-11 21:37:30",
	"UTC":	1691815050,
	"SmartMotionEnable":	true
	}
	
	"Name":	"IPC",
--myboundary
Content-Type: text/plain
Content-Length: 146

Code=VideoMotion;action=Stop;index=0;test=doesnotexist;data={
	"LocaleTime":	"2023-08-11 21:37:46",
	"UTC":	1691815066,
	"Name":	"IPC",
	"SmartMotionEnable":	true
}
`

func Test_EventReader(t *testing.T) {
	eventReader := NewEventReader(strings.NewReader(eventReaderFixture), DefaultEventBoundary)

	output := []Event{
		{
			ContentType:   "text/plain",
			ContentLength: 147,
			Code:          "VideoMotion",
			Action:        "Start",
			Index:         0,
		},
		{
			ContentType:   "text/plain",
			ContentLength: 146,
			Code:          "VideoMotion",
			Action:        "Stop",
			Index:         0,
		},
	}

	for i := 0; ; i++ {
		err := eventReader.Poll()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
		}

		event, err := eventReader.ReadEvent()
		if !assert.NoError(t, err) {
			return
		}

		if !assert.Equal(t, output[i].ContentType, event.ContentType) {
			return
		}
		if !assert.Equal(t, output[i].ContentLength, event.ContentLength) {
			return
		}
		if !assert.Equal(t, output[i].Code, event.Code) {
			return
		}
		if !assert.Equal(t, output[i].Action, event.Action) {
			return
		}
		if !assert.Equal(t, output[i].Index, event.Index) {
			return
		}
	}
}
