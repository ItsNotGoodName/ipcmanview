package dahuasmtp

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/mail"
	"slices"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/event/action"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/emersion/go-smtp"
	"github.com/jhillyerd/enmime"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

type App struct {
	db  sqlite.DB
	bus *event.Bus
	afs afero.Fs
}

func NewApp(db sqlite.DB, bus *event.Bus, afs afero.Fs) App {
	return App{
		db:  db,
		bus: bus,
		afs: afs,
	}
}

// The Backend implements SMTP server methods.
type Backend struct {
	app App
	log zerolog.Logger
}

func NewBackend(app App) *Backend {
	log := log.With().Str("package", "dahuasmtp").Logger()
	return &Backend{
		app: app,
		log: log,
	}
}

func (b *Backend) NewSession(state *smtp.Conn) (smtp.Session, error) {
	address := state.Conn().RemoteAddr().String()
	log := b.log.With().Str("address", address).Logger()

	// log.Debug().Msg("NewSession")

	return &session{
		App:     b.app,
		log:     log,
		address: address,
	}, nil
}

// A Session is returned after EHLO.
type session struct {
	App
	authenticated bool
	log           zerolog.Logger
	address       string
	from          string
	to            string
}

func (s *session) AuthPlain(username, password string) error {
	// s.log.Debug().Str("username", username).Str("password", password).Msg("AuthPlain")

	// err := s.app.AuthSMTPLogin(context.Background(), username, password)
	// if err != nil {
	// 	return smtp.ErrAuthFailed
	// }

	s.authenticated = true

	return nil
}

func (s *session) Mail(from string, opts *smtp.MailOptions) error {
	// s.log.Debug().Str("from", from).Msg("Mail")

	if !s.authenticated {
		return smtp.ErrAuthRequired
	}

	s.from = from

	return nil
}

func (s *session) Rcpt(to string, opts *smtp.RcptOptions) error {
	// s.log.Debug().Str("to", to).Msg("Rcpt")

	if !s.authenticated {
		return smtp.ErrAuthRequired
	}

	s.to = to

	return nil
}

func (s *session) Data(r io.Reader) error {
	// s.log.Debug().Msg("Data")

	ctx := context.Background()
	log := s.log.With().Logger()

	if !s.authenticated {
		return smtp.ErrAuthRequired
	}

	e, err := enmime.ReadEnvelope(r)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read envelope")
		return err
	}

	to := []string{s.to}
	if addresses, err := e.AddressList("To"); err == nil {
		for _, t := range addresses {
			to = append(to, t.Address)
		}
	} else {
		log.Warn().Err(err).Msg("Failed to get 'To' from address list")
	}
	to = slices.Compact(to)

	date, err := e.Date()
	if err != nil && errors.Is(err, mail.ErrHeaderNotPresent) {
		log.Warn().Err(err).Str("date", e.GetHeader("Date")).Msg("Failed to parse date")
	}

	host, _ := core.SplitAddress(s.address)

	device, err := dahua.GetDevice(ctx, s.db, dahua.GetDeviceFilter{
		IP: host,
	})
	if err != nil {
		if core.IsNotFound(err) {
			return err
		}
		log.Err(err).Msg("Failed to get device")
		return err
	}
	log = log.With().Str("device", device.Name).Logger()

	body := dahua.ParseEmailContent(e.Text)
	arg := repo.DahuaCreateEmailMessageParams{
		DeviceID:          device.ID,
		Date:              types.NewTime(date),
		From:              s.from,
		To:                types.NewStringSlice(to),
		Subject:           e.GetHeader("Subject"),
		Text:              e.Text,
		AlarmEvent:        body.AlarmEvent,
		AlarmInputChannel: int64(body.AlarmInputChannel),
		AlarmName:         body.AlarmName,
		CreatedAt:         types.NewTime(time.Now()),
	}

	args := make([]repo.DahuaCreateEmailAttachmentParams, 0, len(e.Attachments))
	for _, a := range e.Attachments {
		args = append(args, repo.DahuaCreateEmailAttachmentParams{
			FileName: a.FileName,
		})
	}

	email, err := dahua.CreateEmail(ctx, s.db, arg, args...)
	if err != nil {
		log.Err(err).Msg("Failed to create email")
		return err
	}
	log = log.With().Int64("message-id", email.Message.ID).Logger()

	for i, attachment := range e.Attachments {
		if err := func() error {
			aferoFile, err := dahua.CreateAferoFile(
				ctx,
				s.db,
				s.afs,
				dahua.AferoForeignKeys{EmailAttachmentID: email.Attachments[i].ID},
				dahua.NewAferoFileName(parseFileExtension(attachment.FileName, http.DetectContentType(attachment.Content))),
			)
			if err != nil {
				log.Err(err).Msg("Failed to create file")
				return err
			}
			defer aferoFile.Close()

			if _, err = aferoFile.Write(attachment.Content); err != nil {
				log.Err(err).Msg("Failed to write file")
				return err
			}

			if err := dahua.ReadyAferoFile(ctx, s.db, aferoFile.ID, aferoFile.File); err != nil {
				log.Err(err).Msg("Failed to ready file")
			}

			return nil
		}(); err != nil {
			return err
		}
	}

	if err = event.CreateEvent(ctx, s.db, action.DahuaEmailCreated.Create(email.Message.ID)); err != nil {
		return err
	}
	s.bus.DahuaEmailCreated(event.DahuaEmailCreated{
		DeviceID: email.Message.DeviceID,
		EmailID:  email.Message.ID,
	})

	log.Info().Msg("Created email")

	return nil
}

func (s *session) Reset() {
	// s.log.Debug().Msg("Reset")
}

func (s *session) Logout() error {
	// s.log.Debug().Msg("Logout")

	return nil
}
