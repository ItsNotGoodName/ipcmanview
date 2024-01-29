package dahuasmtp

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

type Email struct {
	Message     repo.DahuaEmailMessage
	Attachments []repo.DahuaEmailAttachment
}

func CreateDahuaEmail(ctx context.Context, db repo.DB, arg repo.DahuaCreateEmailMessageParams, args ...repo.DahuaCreateEmailAttachmentParams) (Email, error) {
	tx, err := db.BeginTx(ctx, true)
	if err != nil {
		return Email{}, err
	}
	defer tx.Rollback()

	msg, err := db.DahuaCreateEmailMessage(ctx, arg)
	if err != nil {
		return Email{}, err
	}

	atts := make([]repo.DahuaEmailAttachment, 0, len(args))
	for _, a := range args {
		a.MessageID = msg.ID
		att, err := db.DahuaCreateEmailAttachment(ctx, a)
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
