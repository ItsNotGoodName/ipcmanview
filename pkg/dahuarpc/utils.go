package dahuarpc

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type AuthParam struct {
	Encryption string `json:"encryption"`
	Random     string `json:"random"`
	Realm      string `json:"realm"`
}

// HashPassword runs the hashing algorithm for the password.
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
				strings.ToUpper(fmt.Sprintf(
					"%x",
					md5.Sum([]byte(fmt.Sprintf("%s:%s:%s", username, a.Realm, password))))))))))
	default:
		return password
	}
}

type Timestamp string

// NewTimestamp converts the given UTC time to the given location and returns the timestamp.
func NewTimestamp(date time.Time, cameraLocation *time.Location) Timestamp {
	return Timestamp(date.In(cameraLocation).Format("2006-01-02 15:04:05"))
}

// Parse returns the UTC time for the given timestamp and camera location.
func (t Timestamp) Parse(cameraLocation *time.Location) (time.Time, error) {
	if strings.HasSuffix(string(t), "PM") || strings.HasSuffix(string(t), "AM") {
		date, err := time.ParseInLocation("2006-01-02 03:04:05 PM", string(t), cameraLocation)
		if err != nil {
			return date, err
		}

		return date.UTC(), nil
	} else {
		date, err := time.ParseInLocation("2006-01-02 15:04:05", string(t), cameraLocation)
		if err != nil {
			return date, err
		}

		return date.UTC(), nil
	}
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

// Integer is for types that are supposed to integer but for some reason the camera returns a float.
type Integer int64

func (s *Integer) UnmarshalJSON(data []byte) error {
	var number float64
	if err := json.Unmarshal(data, &number); err != nil {
		return err
	}

	*s = Integer(number)

	return nil
}

func (s Integer) Integer() int64 {
	return int64(s)
}

func LoadFileURL(httpAddress, path string) string {
	return fmt.Sprintf("%s/RPC_Loadfile%s", httpAddress, path)
}

func Cookie(session string) string {
	return fmt.Sprintf("WebClientSessionID=%s; DWebClientSessionID=%s; DhWebClientSessionID=%s", session, session, session)
}
