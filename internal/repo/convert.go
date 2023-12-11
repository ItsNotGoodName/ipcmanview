package repo

import "github.com/ItsNotGoodName/ipcmanview/internal/models"

func (c ListDahuaCameraByIDsRow) Convert() models.DahuaConn {
	return models.DahuaConn{
		ID:       c.ID,
		Address:  c.Address,
		Username: c.Username,
		Password: c.Password,
		Location: c.Location.Location,
		Seed:     int(c.Seed),
	}
}

func (c GetDahuaCameraRow) Convert() models.DahuaConn {
	return models.DahuaConn{
		ID:       c.ID,
		Address:  c.Address,
		Username: c.Username,
		Password: c.Password,
		Location: c.Location.Location,
		Seed:     int(c.Seed),
	}
}

func (c ListDahuaCameraRow) Convert() models.DahuaConn {
	return models.DahuaConn{
		ID:       c.ID,
		Address:  c.Address,
		Username: c.Username,
		Password: c.Password,
		Location: c.Location.Location,
		Seed:     int(c.Seed),
	}
}
