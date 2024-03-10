package mediamtx

import (
	"bytes"
	"fmt"
	"net/url"
	"text/template"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/rs/zerolog/log"
)

func NewConfig(host, streamProtocol string, webrtcPort, hlsPort int) (Config, error) {
	pathTemplate := "ipcmanview_dahua_{{.DeviceID}}_{{.Channel}}_{{.Subtype}}"

	var tmpl *template.Template
	if pathTemplate != "" {
		var err error
		tmpl, err = template.New("").Parse(pathTemplate)
		if err != nil {
			return Config{}, err
		}
	}

	var embedAddress string
	switch streamProtocol {
	case "webrtc":
		embedAddress = fmt.Sprintf("http://%s:%d", host, webrtcPort)
	case "hls":
		embedAddress = fmt.Sprintf("http://%s:%d", host, hlsPort)
	default:
		return Config{}, fmt.Errorf("invalid stream protocol: %s", streamProtocol)
	}
	urL, err := url.Parse(embedAddress)
	if err != nil {
		return Config{}, err
	}

	return Config{
		url:          urL,
		pathTemplate: tmpl,
	}, nil
}

type Config struct {
	url          *url.URL
	pathTemplate *template.Template
}

func (c Config) URL() *url.URL {
	return c.url
}

type DahuaStream struct {
	ID           int64
	DeviceID     int64
	Name         string
	Channel      int64
	Subtype      int64
	MediamtxPath string
}

func (c Config) DahuaEmbedPath(stream repo.DahuaStream) string {
	if stream.MediamtxPath != "" {
		return stream.MediamtxPath
	}

	if c.pathTemplate != nil {
		var buffer bytes.Buffer
		err := c.pathTemplate.Execute(&buffer, DahuaStream{
			ID:       stream.ID,
			DeviceID: stream.DeviceID,
			Name:     stream.Name,
			Channel:  stream.Channel,
			Subtype:  stream.Subtype,
		})
		if err != nil {
			log.Err(err).Str("package", "mediamtx").Send()
			return ""
		}

		return buffer.String()
	}

	return ""
}
