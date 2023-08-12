package dahuacgi

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func AudioInputChannelCount(ctx context.Context, cgi Gen) (int, error) {
	method := "devAudioInput.cgi?action=getCollect"

	table, err := OKTable(cgi.CGIGet(ctx, method))
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(table.Get("result"))
}

func AudioOutputChannelCount(ctx context.Context, cgi Gen) (int, error) {
	method := "devAudioOutput.cgi?action=getCollect"

	table, err := OKTable(cgi.CGIGet(ctx, method))
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(table.Get("result"))
}

type HTTPType string

const (
	HTTPTypeSinglePart = "singlepart"
	HTTPTypeMultiPart  = "multipart"
)

type AudioStream struct {
	io.ReadCloser
	ContentType string
}

func AudioStreamGet(ctx context.Context, cgi Gen, channel int, httpType HTTPType) (AudioStream, error) {
	method := "audio.cgi"

	query := url.Values{}
	query.Add("action", "getAudio")
	query.Add("channel", strconv.Itoa(channel))
	query.Add("httptype", string(httpType))
	if len(query) > 0 {
		method += "?" + query.Encode()
	}

	res, err := OK(cgi.CGIGet(ctx, method))
	if err != nil {
		return AudioStream{}, err
	}

	contentType := res.Header.Get("Content-Type")

	return AudioStream{
		ReadCloser:  res.Body,
		ContentType: contentType,
	}, nil
}

// WARNING: this has not been tested yet
func AudioStreamPost(ctx context.Context, cgi Gen, channel int, httpType HTTPType, contentType string, body io.Reader) error {
	method := "audio.cgi"

	query := url.Values{}
	query.Add("action", "postAudio")
	query.Add("channel", strconv.Itoa(channel))
	query.Add("httptype", string(httpType))
	if len(query) > 0 {
		method += "?" + query.Encode()
	}

	headers := http.Header{}
	headers.Add("Content-Type", contentType)
	headers.Add("Content-Length", "9999999")

	res, err := OK(cgi.CGIPost(ctx, method, headers, body))
	if err != nil {
		return err
	}
	res.Body.Close()

	return nil
}
