package dahua

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/validate"
)

func validateCamera(camera models.DahuaCamera) error {
	return validate.Validate.Struct(camera)
}

