package gorise

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
)

func BuildTelegram(cfg Config) (Sender, error) {
	paths := cfg.Paths()
	if len(paths) > 2 {
		return nil, fmt.Errorf("multiple chat ids are not supported")
	}
	if len(paths) != 2 {
		return nil, fmt.Errorf("no chat id supplied")
	}

	return NewTelegram(paths[0], paths[1]), nil
}

func NewTelegram(token, chatID string) Telegram {
	return Telegram{
		token:  token,
		chatID: chatID,
	}
}

// https://github.com/caronc/apprise/wiki/Notify_telegram
type Telegram struct {
	token  string
	chatID string
}

func (t Telegram) Send(ctx context.Context, msg Message) error {
	body := core.First(msg.Body, msg.Title)

	// Get images
	var images []Attachment
	for _, attachment := range msg.Attachments {
		if attachment.IsImage() {
			images = append(images, attachment)
		}
	}

	if len(images) == 0 {
		// Send message
		if err := t.sendMessage(ctx, body); err != nil {
			return err
		}
	} else {
		// Send images
		imagesLength := len(images)
		if imagesLength > 10 {
			imagesLength = 10
		}
		// TODO: use sendMediaGroup when more than 1 attachment

		// Send with 1 image
		if err := t.sendPhoto(ctx, body, images[0].Name, images[0].Reader); err != nil {
			return err
		}

		// Send rest of images
		for i := 1; i < imagesLength; i++ {
			if err := t.sendPhoto(ctx, "", images[i].Name, images[i].Reader); err != nil {
				return err
			}
		}
	}

	return nil
}

type telegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

func (t Telegram) sendMessage(ctx context.Context, text string) error {
	if text == "" {
		return nil
	}

	// Create and send request
	if len(text) > 4096 {
		text = text[:4096]
	}

	// Create request
	values := url.Values{"chat_id": {t.chatID}, "text": {text}}
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.telegram.org/bot"+t.token+"/sendMessage", strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse response
	res := &telegramResponse{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err
	}
	if !res.OK {
		return errors.New(res.Description)
	}

	return nil
}

func (t Telegram) sendPhoto(ctx context.Context, caption, name string, data io.Reader) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Photo
	w, err := writer.CreateFormFile("photo", name)
	if err != nil {
		return err
	}

	if _, err := io.Copy(body, data); err != nil {
		return err
	}

	// Caption
	if caption != "" {
		w, err = writer.CreateFormField("caption")
		if err != nil {
			return err
		}
		if len(caption) > 1024 {
			caption = caption[:1024]
		}
		w.Write([]byte(caption))
	}

	// Close
	if err := writer.Close(); err != nil {
		return err
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.telegram.org/bot"+t.token+"/sendPhoto?chat_id="+t.chatID, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse response
	res := &telegramResponse{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err
	}
	if !res.OK {
		return errors.New(res.Description)
	}

	return nil
}
