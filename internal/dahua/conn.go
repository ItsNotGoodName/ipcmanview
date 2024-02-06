package dahua

import (
	"net/url"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

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
