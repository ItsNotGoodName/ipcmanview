package dahua

import (
	"net/url"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

func ConnFrom(v repo.DahuaFatDevice) Conn {
	return Conn{
		ID:       v.DahuaDevice.ID,
		URL:      v.DahuaDevice.Url.URL,
		Username: v.DahuaDevice.Username,
		Password: v.DahuaDevice.Password,
		Location: v.DahuaDevice.Location.Location,
		Feature:  v.DahuaDevice.Feature,
		Seed:     int(v.Seed),
	}
}

func ConnsFrom(devices []repo.DahuaFatDevice) []Conn {
	conns := make([]Conn, 0, len(devices))
	for _, v := range devices {
		conns = append(conns, ConnFrom(v))
	}
	return conns
}

// Conn is the bare minumum information required to create a connection to a Dahua device.
type Conn struct {
	ID       int64
	URL      *url.URL
	Username string
	Password string
	Location *time.Location
	Feature  models.DahuaFeature
	Seed     int
}

func (lhs Conn) EQ(rhs Conn) bool {
	return lhs.URL.String() == rhs.URL.String() &&
		lhs.Username == rhs.Username &&
		lhs.Password == rhs.Password &&
		lhs.Location.String() == rhs.Location.String() &&
		lhs.Feature == rhs.Feature &&
		lhs.Seed == rhs.Seed
}
