package dahua

import (
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

// Camera is used internally.
type Camera struct {
	ID        int64
	Address   string
	Username  string
	Password  string
	CreatedAt time.Time
}

func NewCamera(cam models.DahuaCamera) Camera {
	return Camera{
		ID:       cam.ID,
		Address:  cam.Address,
		Username: cam.Username,
		Password: cam.Password,
	}
}

func (c Camera) Different(cam Camera) bool {
	return c.Address != cam.Address ||
		c.Username != cam.Username ||
		c.Password != cam.Password
}
