package dahuasmtp

import (
	"context"
	"errors"
	"io"
	"net/mail"
	"slices"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/emersion/go-smtp"
	"github.com/jhillyerd/enmime"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// The Backend implements SMTP server methods.
type Backend struct {
	log zerolog.Logger
}

func NewBackend() *Backend {
	log := log.With().Str("package", "dahuasmtp").Logger()
	return &Backend{
		log: log,
	}
}

func (b *Backend) NewSession(state *smtp.Conn) (smtp.Session, error) {
	address := state.Conn().RemoteAddr().String()
	log := b.log.With().Str("address", address).Logger()

	// log.Debug().Msg("NewSession")

	return &session{
		log:     log,
		address: address,
	}, nil
}

// A Session is returned after EHLO.
type session struct {
	log     zerolog.Logger
	address string
	from    string
	to      string
}

func (s *session) AuthPlain(username, password string) error {
	// s.log.Debug().Str("username", username).Str("password", password).Msg("AuthPlain")

	return nil
}

func (s *session) Mail(from string, opts *smtp.MailOptions) error {
	// s.log.Debug().Str("from", from).Msg("Mail")

	s.from = from

	return nil
}

func (s *session) Rcpt(to string, opts *smtp.RcptOptions) error {
	// s.log.Debug().Str("to", to).Msg("Rcpt")

	s.to = to

	return nil
}

func (s *session) Data(r io.Reader) error {
	// s.log.Debug().Msg("Data")

	ctx := context.Background()
	log := s.log.With().Logger()

	// Get device by IP
	host, _ := core.SplitAddress(s.address)
	device, err := dahua.GetDevice(ctx, dahua.GetDeviceFilter{
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

	// Read
	e, err := enmime.ReadEnvelope(r)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read envelope")
		return err
	}

	// Parse to
	to := []string{s.to}
	if addresses, err := e.AddressList("To"); err == nil {
		for _, t := range addresses {
			to = append(to, t.Address)
		}
	} else {
		log.Warn().Err(err).Msg("Failed to get 'To' from address list")
	}
	to = slices.Compact(to)

	// Parse date
	date, err := e.Date()
	if err != nil && errors.Is(err, mail.ErrHeaderNotPresent) {
		log.Warn().Err(err).Str("date", e.GetHeader("Date")).Msg("Failed to parse date")
	}

	// Create email
	body := dahua.ParseEmailContent(e.Text)
	attachments := make([]dahua.CreateEmailParamsAttachment, 0, len(e.Attachments))
	for _, a := range e.Attachments {
		attachments = append(attachments, dahua.CreateEmailParamsAttachment{
			FileName: a.FileName,
			Data:     a.Content,
		})
	}
	arg := dahua.CreateEmailParams{
		DeviceID:          device.ID,
		Date:              date,
		From:              s.from,
		To:                to,
		Subject:           e.GetHeader("Subject"),
		Text:              e.Text,
		AlarmEvent:        body.AlarmEvent,
		AlarmInputChannel: body.AlarmInputChannel,
		AlarmName:         body.AlarmName,
		Attachments:       attachments,
	}
	messageID, err := dahua.CreateEmail(ctx, arg)
	if err != nil {
		log.Err(err).Msg("Failed to create email")
		return err
	}
	log.Info().Int64("message_id", messageID).Msg("Created email")

	return nil
}

func (s *session) Reset() {
	// s.log.Debug().Msg("Reset")
}

func (s *session) Logout() error {
	// s.log.Debug().Msg("Logout")

	return nil
}
