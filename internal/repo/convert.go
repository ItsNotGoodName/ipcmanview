package repo

import "github.com/ItsNotGoodName/ipcmanview/internal/models"

// this is stupid

func (c User) Convert() models.User {
	return models.User{
		ID:       c.ID,
		Email:    c.Email,
		Username: c.Username,
		Password: c.Password,
	}
}

func (c Group) Convert() models.Group {
	return models.Group{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
	}
}
