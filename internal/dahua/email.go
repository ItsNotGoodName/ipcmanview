package dahua

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
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
	Content  []byte
}

func CreateEmail(ctx context.Context, arg CreateEmailParams) (int64, error) {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return 0, err
	}

	// Create message and attachments
	res, err := createEmail(ctx, &arg)
	if err != nil {
		return 0, err
	}

	// Save attachment files
	for i := range arg.Attachments {
		err := createEmailAttachmentFile(ctx, createEmailAttachmentFileParams{
			AttachmentID: res.AttachmentIDs[i],
			FileID:       res.FileIDs[i],
			FileName:     arg.Attachments[i].FileName,
			Content:      arg.Attachments[i].Content,
		})
		if err != nil {
			return 0, err
		}
	}

	// Publish email created event
	if err := system.CreateEvent(ctx, app.DB, action.DahuaEmailCreated.Create(res.MessageID)); err != nil {
		return 0, err
	}
	app.Hub.DahuaEmailCreated(bus.DahuaEmailCreated{
		DeviceID:  arg.DeviceID,
		MessageID: res.MessageID,
	})

	return res.MessageID, nil
}

type createEmailResult struct {
	MessageID     int64
	AttachmentIDs []int64
	FileIDs       []int64
}

func createEmail(ctx context.Context, arg *CreateEmailParams) (createEmailResult, error) {
	tx, err := app.DB.BeginTx(ctx, true)
	if err != nil {
		return createEmailResult{}, err
	}
	defer tx.Rollback()

	now := types.NewTime(time.Now())
	date := types.NewTime(arg.Date)

	msgID, err := app.DB.C().DahuaCreateEmailMessage(ctx, repo.DahuaCreateEmailMessageParams{
		DeviceID:          arg.DeviceID,
		Date:              date,
		From:              arg.From,
		To:                types.NewStringSlice(arg.To),
		Subject:           arg.Subject,
		Text:              arg.Text,
		AlarmEvent:        arg.AlarmEvent,
		AlarmInputChannel: int64(arg.AlarmInputChannel),
		AlarmName:         arg.AlarmName,
		CreatedAt:         now,
	})
	if err != nil {
		return createEmailResult{}, err
	}

	attachmentIDs := make([]int64, len(arg.Attachments))
	fileIDs := make([]int64, len(arg.Attachments))
	for i, v := range arg.Attachments {
		att, err := app.DB.C().DahuaCreateEmailAttachment(ctx, repo.DahuaCreateEmailAttachmentParams{
			MessageID: msgID,
			FileName:  v.FileName,
		})
		if err != nil {
			return createEmailResult{}, err
		}
		attachmentIDs[i] = att

		fileID, err := app.DB.C().DahuaCreateFile(ctx, repo.DahuaCreateFileParams{
			DeviceID:  arg.DeviceID,
			Channel:   0,
			StartTime: date,
			EndTime:   date,
			Length:    int64(len(v.Content)),
			Type:      models.DahuaFileType_JPG,
			FilePath:  fmt.Sprintf("ipcmanview+email://%d", att),
			Storage:   models.StorageLocal,
			Source:    models.DahuaFileSource_Email,
			UpdatedAt: now,
		})
		if err != nil {
			return createEmailResult{}, err
		}
		fileIDs[i] = fileID
	}

	if err := tx.Commit(); err != nil {
		return createEmailResult{}, err
	}

	return createEmailResult{
		MessageID:     msgID,
		AttachmentIDs: attachmentIDs,
		FileIDs:       fileIDs,
	}, nil
}

type createEmailAttachmentFileParams struct {
	AttachmentID int64
	FileID       int64
	FileName     string
	Content      []byte
}

func createEmailAttachmentFile(ctx context.Context, arg createEmailAttachmentFileParams) error {
	aferoFile, err := createAferoFile(
		ctx,
		aferoForeignKeys{EmailAttachmentID: arg.AttachmentID, FileID: arg.FileID},
		newAferoFileName(parseFileExtension(arg.FileName, http.DetectContentType(arg.Content))),
	)
	if err != nil {
		return err
	}
	defer aferoFile.Close()

	if _, err = aferoFile.Write(arg.Content); err != nil {
		return err
	}

	return aferoFile.Ready(ctx)
}
