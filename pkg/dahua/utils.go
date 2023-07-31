package dahua

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

type JSON = map[string]any

type AuthParam struct {
	Encryption string `json:"encryption,omitempty"`
	Random     string `json:"random,omitempty"`
	Realm      string `json:"realm,omitempty"`
}

// HashPassword runs the bespoke hashing algorithm for the password.
func (a AuthParam) HashPassword(username, password string) string {
	switch a.Encryption {
	case "Basic":
		return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	case "Default":
		return strings.ToUpper(fmt.Sprintf("%x",
			md5.Sum([]byte(fmt.Sprintf(
				"%s:%s:%s",
				username,
				a.Random,
				strings.ToUpper(fmt.Sprintf("%x",
					md5.Sum([]byte(fmt.Sprintf("%s:%s:%s", username, a.Realm, password))),
				)))))))
	default:
		return password
	}
}

type Timestamp string

// Parse returns the UTC time for the given timestamp and location.
func (t Timestamp) Parse(cameraLocation *time.Location) (time.Time, error) {
	date, err := time.ParseInLocation("2006-01-02 15:04:05", string(t), cameraLocation)
	if err != nil {
		return date, err
	}

	return date.UTC(), nil
}

// NewTimestamp converts the given UTC time to the given location and returns the timestamp.
func NewTimestamp(date time.Time, cameraLocation *time.Location) Timestamp {
	return Timestamp(date.In(cameraLocation).Format("2006-01-02 15:04:05"))
}

// ExtractFilePathTags extracts tags that are surrounded by brackets from the given file path.
func ExtractFilePathTags(filePath string) []string {
	search := filePath
	idx := strings.LastIndex(filePath, "/")
	if idx != -1 {
		search = filePath[idx:]
	}

	var tags []string
	tokens := strings.Split(search, "[")
	for i := 1; i < len(tokens); i++ {
		if end := strings.Index(tokens[i], "]"); end != -1 {
			tags = append(tags, tokens[i][:end])
		}
	}

	return tags
}
