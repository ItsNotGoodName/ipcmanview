package mediamtx

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/rs/zerolog/log"
)

type Config struct {
	embedAddress string
	pathTemplate *template.Template
}

func NewConfig(host, pathTemplate, streamProtocol string, webrtcPort, hlsPort int) (Config, error) {
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

	return Config{
		embedAddress: embedAddress,
		pathTemplate: tmpl,
	}, nil
}

type DahuaStream struct {
	ID           int64
	DeviceID     int64
	Name         string
	Channel      int64
	Subtype      int64
	MediamtxPath string
}

func (c Config) DahuaEmbedURL(stream repo.DahuaStream) string {
	if c.embedAddress == "" {
		return ""
	}

	path, err := c.dahuaPath(stream)
	if err != nil {
		log.Err(err).Msg("Failed to get mediamtx path")
		return ""
	}

	return fmt.Sprintf("%s/%s", c.embedAddress, path)
}

func (c Config) dahuaPath(stream repo.DahuaStream) (string, error) {
	if stream.MediamtxPath != "" {
		return stream.MediamtxPath, nil
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
			return "", err
		}

		return buffer.String(), nil
	}

	return "", nil
}
