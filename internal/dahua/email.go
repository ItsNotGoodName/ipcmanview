package dahua

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/system"
	"github.com/ItsNotGoodName/ipcmanview/internal/system/action"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
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

type CreateEmailParams struct {
	DeviceID          int64
	Date              time.Time
	From              string
	To                []string
	Subject           string
	Text              string
	AlarmEvent        string
	AlarmInputChannel int
	AlarmName         string
	Attachments       []CreateEmailParamsAttachment
}

type CreateEmailParamsAttachment struct {
	FileName string
	Data     []byte
}

func CreateEmail(ctx context.Context, arg CreateEmailParams) (int64, error) {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return 0, err
	}

	// Create message and attachments
	msg, atts, err := createEmail(ctx, &arg)
	if err != nil {
		return 0, err
	}

	// Save attachment files
	for i := range atts {
		err := createEmailAttachmentFile(ctx, atts[i].ID, arg.Attachments[i].FileName, arg.Attachments[i].Data)
		if err != nil {
			return 0, err
		}
	}

	// Publish email created event
	if err := system.CreateEvent(ctx, app.DB, action.DahuaEmailCreated.Create(msg.ID)); err != nil {
		return 0, err
	}
	app.Hub.DahuaEmailCreated(bus.DahuaEmailCreated{
		DeviceID:  arg.DeviceID,
		MessageID: msg.ID,
	})

	return msg.ID, nil
}

func createEmail(ctx context.Context, arg *CreateEmailParams) (repo.DahuaEmailMessage, []repo.DahuaEmailAttachment, error) {
	tx, err := app.DB.BeginTx(ctx, true)
	if err != nil {
		return repo.DahuaEmailMessage{}, nil, err
	}
	defer tx.Rollback()

	msg, err := app.DB.C().DahuaCreateEmailMessage(ctx, repo.DahuaCreateEmailMessageParams{
		DeviceID:          arg.DeviceID,
		Date:              types.NewTime(arg.Date),
		From:              arg.From,
		To:                types.NewStringSlice(arg.To),
		Subject:           arg.Subject,
		Text:              arg.Text,
		AlarmEvent:        arg.AlarmEvent,
		AlarmInputChannel: int64(arg.AlarmInputChannel),
		AlarmName:         arg.AlarmName,
		CreatedAt:         types.NewTime(time.Now()),
	})
	if err != nil {
		return repo.DahuaEmailMessage{}, nil, err
	}

	atts := make([]repo.DahuaEmailAttachment, 0, len(arg.Attachments))
	for _, a := range arg.Attachments {
		att, err := app.DB.C().DahuaCreateEmailAttachment(ctx, repo.DahuaCreateEmailAttachmentParams{
			MessageID: msg.ID,
			FileName:  a.FileName,
		})
		if err != nil {
			return repo.DahuaEmailMessage{}, nil, err
		}
		atts = append(atts, att)
	}

	if err := tx.Commit(); err != nil {
		return repo.DahuaEmailMessage{}, nil, err
	}

	return msg, atts, nil
}

func createEmailAttachmentFile(ctx context.Context, attachmentID int64, fileName string, content []byte) error {
	aferoFile, err := createAferoFile(
		ctx,
		aferoForeignKeys{EmailAttachmentID: attachmentID},
		newAferoFileName(parseFileExtension(fileName, http.DetectContentType(content))),
	)
	if err != nil {
		return err
	}
	defer aferoFile.Close()

	if _, err = aferoFile.Write(content); err != nil {
		return err
	}

	return aferoFile.Ready(ctx)
}
