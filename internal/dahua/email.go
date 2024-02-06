package dahua

import (
	"context"
	"strconv"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/ssq"
	sq "github.com/Masterminds/squirrel"
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

type ListLatestEmailsResult struct {
	repo.DahuaEmailMessage
	AttachmentCount int64
}

func ListLatestEmails(ctx context.Context, db repo.DB, count int) ([]ListLatestEmailsResult, error) {
	sb := sq.
		Select("dahua_email_messages.*", "COUNT(dahua_email_attachments.id) AS attachment_count").
		From("dahua_email_messages").
		LeftJoin("dahua_email_attachments ON dahua_email_attachments.message_id = dahua_email_messages.id").
		OrderBy("created_at DESC").
		GroupBy("dahua_email_messages.id").
		Limit(uint64(count))

	var res []ListLatestEmailsResult
	err := ssq.Query(ctx, db, &res, repo.DahuaSelectFilter(ctx, sb, "dahua_email_messages.device_id", models.DahuaPermissionLevelAdmin))
	if err != nil {
		return nil, err
	}

	return res, nil
}
