package dahua

import (
	"context"
	"net/url"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

func NewConn(v repo.DahuaDevice, seed int64) Conn {
	return Conn{
		ID:       v.ID,
		URL:      v.Url.URL,
		Username: v.Username,
		Password: v.Password,
		Location: v.Location.Location,
		Feature:  v.Feature,
		Seed:     int(seed),
	}
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

func ListConns(ctx context.Context, db repo.DB) ([]Conn, error) {
	rows, err := db.DahuaListDevices(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]Conn, 0, len(rows))
	for i := range rows {
		res = append(res, NewConn(rows[i].DahuaDevice, rows[i].Seed))
	}
	return res, nil
}

func GetConn(ctx context.Context, db repo.DB, id int64) (Conn, error) {
	row, err := db.DahuaGetDevice(ctx, id)
	if err != nil {
		return Conn{}, err
	}
	return NewConn(row.DahuaDevice, row.Seed), nil
}
