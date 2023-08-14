package dahuacgi

import (
	"context"
	"io"
	"strconv"
)

func AudioInputChannelCount(ctx context.Context, c Client) (int, error) {
	req := NewRequest("devAudioInput.cgi").
		QueryString("action", "getCollect")

	table, err := OKTable(c.CGIGet(ctx, req))
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(table.Get("result"))
}

func AudioOutputChannelCount(ctx context.Context, c Client) (int, error) {
	req := NewRequest("devAudioOutput.cgi").
		QueryString("action", "getCollect")

	table, err := OKTable(c.CGIGet(ctx, req))
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

func AudioStreamGet(ctx context.Context, c Client, channel int, httpType HTTPType) (AudioStream, error) {
	if channel == 0 {
		channel = 1
	}

	req := NewRequest("audio.cgi").
		QueryString("action", "getAudio").
		QueryInt("channel", channel).
		QueryString("httptype", string(httpType))

	res, err := OK(c.CGIGet(ctx, req))
	if err != nil {
		return AudioStream{}, err
	}

	contentType := res.Header.Get("Content-Type")

	return AudioStream{
		ReadCloser:  res.Body,
		ContentType: contentType,
	}, nil
}

// INFO: The following reasons are why streaming audio to the camera cannot be added.
// - The HTTP digest library (github.com/icholy/digest) copies the body before sending the real request,
//   this does not work if the body is infinite like what we are doing.
// - I swear that my camera SD2A500-GN-A-PV (Build Date: 2022-08-26) has a broken AudioStreamPost CGI API.
//   I tested it with cURL and it would keep doing a connection reset after sending a bit of audio data.
// - The current Client interface does not support POST, but that is an easy fix the HTTP digest library is fixed.

// func AudioStreamPost(ctx context.Context, c Client, channel int, httpType HTTPType, contentType string, body io.Reader) error {
// 	if channel == 0 {
// 		channel = 1
// 	}
//
// 	req := NewRequest("audio.cgi").
// 		QueryString("action", "postAudio").
// 		QueryInt("channel", channel).
// 		QueryString("httptype", string(httpType)).
// 		Header("Content-Type", contentType).
// 		Body(body)
//
// 	res, err := OK(c.CGIPost(ctx, req))
// 	if err != nil {
// 		return err
// 	}
// 	res.Body.Close()
//
// 	return nil
// }
