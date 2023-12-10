package sqlc

import "github.com/ItsNotGoodName/ipcmanview/internal/models"

func (c ListDahuaCameraByIDsRow) Convert() models.DahuaCamera {
	return models.DahuaCamera{
		ID:        c.ID,
		Address:   c.Address,
		Username:  c.Username,
		Password:  c.Password,
		Location:  c.Location,
		Seed:      int(c.Seed),
		CreatedAt: c.CreatedAt.Time,
	}
}

func ConvertListDahuaCameraByIDsRow(list []ListDahuaCameraByIDsRow) []models.DahuaCamera {
	res := make([]models.DahuaCamera, 0, len(list))
	for _, row := range list {
		res = append(res, row.Convert())
	}
	return res
}

func (c GetDahuaCameraRow) Convert() models.DahuaCamera {
	return models.DahuaCamera{
		ID:        c.ID,
		Address:   c.Address,
		Username:  c.Username,
		Password:  c.Password,
		Location:  c.Location,
		Seed:      int(c.Seed),
		CreatedAt: c.CreatedAt.Time,
	}
}

func (c ListDahuaCameraRow) Convert() models.DahuaCamera {
	return models.DahuaCamera{
		ID:        c.ID,
		Address:   c.Address,
		Username:  c.Username,
		Password:  c.Password,
		Location:  c.Location,
		Seed:      int(c.Seed),
		CreatedAt: c.CreatedAt.Time,
	}
}

func ConvertListDahuaCameraRow(list []ListDahuaCameraRow) []models.DahuaCamera {
	res := make([]models.DahuaCamera, 0, len(list))
	for _, row := range list {
		res = append(res, row.Convert())
	}
	return res
}
