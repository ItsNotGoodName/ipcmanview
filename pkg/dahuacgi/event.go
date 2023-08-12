package dahuacgi

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
)

type EventBoundary string

const defaultEventBoundary EventBoundary = "--myboundary"

// EventManager attaches to all events.
func EventManager(ctx context.Context, cgi Gen, heartbeat int) (EventSession, error) {
	method := "eventManager.cgi"

	query := url.Values{}
	query.Add("action", "attach")
	query.Add("codes", "[All]")
	if heartbeat != 0 {
		query.Add("heartbeat", strconv.Itoa(heartbeat))
	}
	if len(query) > 0 {
		method += "?" + query.Encode()
	}

	rd, err := OKBody(cgi.CGIGet(ctx, method))
	if err != nil {
		return EventSession{}, err
	}

	return EventSession{
		br:       bufio.NewReader(rd),
		boundary: string(defaultEventBoundary), // TODO: parse boundary from Content-Type
	}, nil
}

type EventSession struct {
	br       *bufio.Reader
	boundary string
}

func NewEventSession(br *bufio.Reader) EventSession {
	return EventSession{
		br:       br,
		boundary: string(defaultEventBoundary),
	}
}

type Event struct {
	ContentType   string
	ContentLength int
	Code          string
	Action        string
	Index         int
	Data          []byte
}

// Poll waits for the next event.
func (es EventSession) Poll() error {
	for {
		s, _, err := es.br.ReadLine()
		if err != nil {
			return err
		}
		if strings.HasPrefix(string(s), es.boundary) {
			return nil
		}
	}
}

// Read parses the event. It should only be called after Poll.
func (es EventSession) Read() (Event, error) {
	// Parse headers
	headers, err := es.seekAfterEmptyLine()
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
	body, err := es.seekAfterEmptyLine()
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
func (es EventSession) seekAfterEmptyLine() (string, error) {
	var tokens string
	for {
		var s string
		for {
			b, isPrefix, err := es.br.ReadLine()
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
