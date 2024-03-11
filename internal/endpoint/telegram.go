package endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
)

func BuildTelegram(urL string) (Sender, error) {
	_, after, _ := strings.Cut(urL, "://")

	paths := strings.Split(after, "/")
	if len(paths) < 2 {
		return nil, fmt.Errorf("no chat id supplied")
	}

	return NewTelegram(paths[0], paths[1], paths[2:]...), nil
}

func NewTelegram(token, chatID string, chatIDs ...string) Telegram {
	return Telegram{
		token:   token,
		chatIDs: append([]string{chatID}, chatIDs...),
	}
}

type Telegram struct {
	token   string
	chatIDs []string
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

	var errs []error
	if len(images) == 0 {
		// Send message
		for _, chatID := range t.chatIDs {
			if err := t.sendMessage(ctx, body, chatID); err != nil {
				errs = append(errs, err)
			}
		}
	} else {
		// Send images
		imagesLength := len(images)
		if imagesLength > 10 {
			imagesLength = 10
		}
		for _, chatID := range t.chatIDs {
			// TODO: use sendMediaGroup when more than 1 attachment

			// Send with 1 image
			if err := t.sendPhoto(ctx, chatID, body, images[0].Name, images[0].Data); err != nil {
				errs = append(errs, err)
				continue
			}

			// Send rest of images
			for i := 1; i < imagesLength; i++ {
				if err := t.sendPhoto(ctx, chatID, "", images[i].Name, images[i].Data); err != nil {
					errs = append(errs, err)
					break
				}
			}
		}
	}

	return errors.Join(errs...)
}

type telegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

func (t Telegram) sendMessage(ctx context.Context, text string, chatID string) error {
	if text == "" {
		return nil
	}

	// Create and send request
	if len(text) > 4096 {
		text = text[:4096]
	}

	// Create request
	values := url.Values{"chat_id": {chatID}, "text": {text}}
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

func (t Telegram) sendPhoto(ctx context.Context, chatID string, caption, name string, data []byte) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Photo
	w, err := writer.CreateFormFile("photo", name)
	if err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
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
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.telegram.org/bot"+t.token+"/sendPhoto?chat_id="+chatID, body)
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
