package dahua

import (
	"context"
	"strconv"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
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

type Email struct {
	Message     repo.DahuaEmailMessage
	Attachments []repo.DahuaEmailAttachment
}

func CreateEmail(ctx context.Context, db sqlite.DB, arg repo.DahuaCreateEmailMessageParams, args ...repo.DahuaCreateEmailAttachmentParams) (Email, error) {
	tx, err := db.BeginTx(ctx, true)
	if err != nil {
		return Email{}, err
	}
	defer tx.Rollback()

	msg, err := db.C().DahuaCreateEmailMessage(ctx, arg)
	if err != nil {
		return Email{}, err
	}

	atts := make([]repo.DahuaEmailAttachment, 0, len(args))
	for _, a := range args {
		a.MessageID = msg.ID
		att, err := db.C().DahuaCreateEmailAttachment(ctx, a)
		if err != nil {
			return Email{}, err
		}
		atts = append(atts, att)
	}

	return Email{
		Message:     msg,
		Attachments: atts,
	}, tx.Commit()
}
