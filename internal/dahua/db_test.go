package dahua

import (
	"context"
	"testing"
	"time"

	"github.com/ItsNotGoodName/ipcmango/internal/core"
	"github.com/ItsNotGoodName/ipcmango/internal/dbtest"
	"github.com/stretchr/testify/assert"
)

func TestCamera(t *testing.T) {
	ctx := context.Background()
	db, close := dbtest.Connect(ctx)
	defer close()

	// Create
	coreCam, err := core.NewDahuaCamera(core.DahuaCameraCreate{
		Address:  "localhost",
		Username: "Username",
		Password: "Password",
	})
	assert.NoError(t, err)

	createCam, err := DB.CameraCreate(ctx, db, coreCam)
	assert.NoError(t, err)

	assert.NotEqual(t, coreCam.ID, createCam.ID, "should have new id")
	coreCam.ID = createCam.ID

	assert.NotEqual(t, time.Time{}, createCam.CreatedAt, "should not have default CreatedAt")
	coreCam.CreatedAt = createCam.CreatedAt

	assert.Equal(t, coreCam.Address, createCam.Address)
	assert.Equal(t, coreCam.Username, createCam.Username)
	assert.Equal(t, coreCam.Password, createCam.Password)

	// Update
	updateAddress := "user"

	update := core.
		NewDahuaCameraUpdate(createCam.ID).
		AddressUpdate(updateAddress)

	updateCam, err := DB.CameraUpdate(ctx, db, update)
	assert.NoError(t, err)
	assert.Equal(t, updateAddress, updateCam.Address)

	// Get
	{
		getCam, err := DB.CameraGet(ctx, db, updateCam.ID)
		assert.NoError(t, err)
		assert.Equal(t, updateCam, getCam)
	}

	// Delete
	{
		value, err := update.Value()
		assert.NoError(t, err)
		err = DB.CameraDelete(ctx, db, value.ID)
		assert.NoError(t, err)
		err = DB.CameraDelete(ctx, db, value.ID)
		assert.Error(t, err)
	}
}
