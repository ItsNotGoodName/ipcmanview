package dahuacgi

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/textproto"
	"strconv"
	"strings"
)

type EventBoundary string

const DefaultEventBoundary EventBoundary = "myboundary"

type EventManager struct {
	io.ReadCloser
	Boundary EventBoundary
}

func EventManagerGet(ctx context.Context, c Client, heartbeat int) (EventManager, error) {
	req := NewRequest("eventManager.cgi").
		QueryString("action", "attach").
		QueryString("codes", "[All]")

	if heartbeat != 0 {
		req.QueryInt("heartbeat", heartbeat)
	}

	res, err := OK(c.CGIGet(ctx, req))
	if err != nil {
		return EventManager{}, err
	}

	// Parse boundary
	contentType := res.Header.Get("Content-Type")
	boundary := DefaultEventBoundary
	for _, token := range strings.Split(contentType, ";") {
		maybeKV := strings.SplitN(token, "=", 2)
		if len(maybeKV) != 2 {
			continue
		}
		if maybeKV[0] != "boundary" {
			continue
		}

		boundary = EventBoundary(maybeKV[1])
		break
	}

	return EventManager{
		ReadCloser: res.Body,
		Boundary:   boundary,
	}, nil
}

func (em EventManager) Reader() EventReader {
	return NewEventReader(em.ReadCloser, em.Boundary)
}

type EventReader struct {
	br                 *bufio.Reader
	boundaryWithPrefix string
}

func NewEventReader(rd io.Reader, boundary EventBoundary) EventReader {
	return EventReader{
		br:                 bufio.NewReader(rd),
		boundaryWithPrefix: "--" + string(boundary),
	}
}

type Event struct {
	ContentType   string
	ContentLength int
	Code          string
	Action        string
	Index         int
	Data          json.RawMessage
}

// Poll waits for the next event boundary.
func (er EventReader) Poll() error {
	for {
		s, _, err := er.br.ReadLine()
		if err != nil {
			return err
		}
		if strings.HasPrefix(string(s), er.boundaryWithPrefix) {
			return nil
		}
	}
}

// ReadEvent parses the next event. Should be called after Poll.
func (er EventReader) ReadEvent() (Event, error) {
	// Parse headers
	headers, err := er.seekAfterEmptyLine()
	if err != nil {
		return Event{}, err
	}
	mimeHeaders, err := textproto.NewReader(bufio.NewReader(strings.NewReader(headers + "\r\n"))).ReadMIMEHeader()
	if err != nil {
		return Event{}, err
	}

	contentType := mimeHeaders.Get("Content-Type")
	contentLength, _ := strconv.Atoi(mimeHeaders.Get("Content-Length"))

	// Parse body
	body, err := er.seekAfterEmptyLine()
	if err != nil {
		return Event{}, err
	}
	kv, data := eventParseBody(bufio.NewReader(strings.NewReader(body)))

	code := kv["code"]
	action := kv["action"]
	index, _ := strconv.Atoi(kv["index"])

	return Event{
		ContentType:   contentType,
		ContentLength: contentLength,
		Code:          code,
		Action:        action,
		Index:         index,
		Data:          data,
	}, nil
}

// seekAfterEmptyLine moves reader until it is after the next empty line.
func (er EventReader) seekAfterEmptyLine() (string, error) {
	var tokens string
	for {
		var s string
		for {
			b, isPrefix, err := er.br.ReadLine()
			if err != nil {
				return "", err
			}
			if isPrefix {
				s += string(b)
				continue
			}
			s = string(b)
			break
		}

		if len(s) == 0 {
			return tokens, nil
		}

		tokens += s + "\n"
	}
}

func eventParseBody(br *bufio.Reader) (map[string]string, []byte) {
	m := make(map[string]string)

	for {
		// Key
		key, err := eventTextUntil(br, '=')
		if err != nil {
			return m, []byte{}
		}
		key = strings.ToLower(key)

		// Break when the next value is JSON
		isJSON, err := eventPeekJSON(br)
		if err != nil {
			return m, []byte{}
		}
		if isJSON {
			break
		}

		// Value
		value, err := eventTextUntil(br, ';')
		if err != nil {
			if errors.Is(err, io.EOF) {
				m[key] = value
			}
			return m, []byte{}
		}
		m[key] = value
	}

	// Naive way of grabbing JSON
	data, err := io.ReadAll(br)
	if err != nil {
		return m, []byte{}
	}

	return m, data
}

// eventTextUntil reads until the delimiter.
// The delimiter is not returned but is still consumed.
func eventTextUntil(br *bufio.Reader, delim byte) (string, error) {
	value, err := br.ReadString(delim)
	if err != nil {
		return "", err
	}
	if len(value) > 0 {
		value = value[:len(value)-1]
	}

	return value, nil
}

// eventPeekJSON checks if the next value is JSON.
func eventPeekJSON(br *bufio.Reader) (bool, error) {
	maybeCurly, _, err := br.ReadRune()
	if err != nil {
		return false, err
	}
	br.UnreadRune()

	return maybeCurly == '{', nil
}
