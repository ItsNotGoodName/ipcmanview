package dahua

import (
	"strconv"
	"strings"
)

type EmailContent struct {
	AlarmEvent        string
	AlarmInputChannel int
	AlarmDeviceName   string
	AlarmName         string
	IPAddress         string
}

func ParseEmailContent(text string) EmailContent {
	var content EmailContent
	for _, line := range strings.Split(text, "\n") {
		kv := strings.SplitN(line, ":", 2)
		if len(kv) != 2 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		switch key {
		case "Alarm Event":
			content.AlarmEvent = value
		case "Alarm Input Channel":
			channel, _ := strconv.Atoi(value)
			content.AlarmInputChannel = channel
		case "Alarm Device Name":
			content.AlarmDeviceName = value
		case "Alarm Name":
			content.AlarmName = value
		case "IP Address":
			content.IPAddress = value
		default:
		}
	}

	return content
}
