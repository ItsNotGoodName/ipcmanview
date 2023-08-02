package db

import (
	"context"
	"testing"
	"time"

	"github.com/ItsNotGoodName/ipcmango/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestDahuaCamera(t *testing.T) {
	dbCtx, close := TestConnect(context.Background())
	defer close()

	// Create
	coreCam, err := core.NewDahuaCamera(core.DahuaCameraCreate{
		Address:  "localhost",
		Username: "Username",
		Password: "Password",
	})
	assert.NoError(t, err)

	createCam, err := DahuaCameraCreate(dbCtx, coreCam)
	assert.NoError(t, err)

	assert.NotEqual(t, coreCam.ID, createCam.ID, "should have new id")
	coreCam.ID = createCam.ID

	assert.NotEqual(t, time.Time{}, createCam.CreatedAt, "should not have default CreatedAt")
	coreCam.CreatedAt = createCam.CreatedAt

	assert.Equal(t, coreCam, createCam)

	// Update
	updateAddress := "user"

	update := core.
		NewDahuaCameraUpdate(createCam.ID).
		AddressUpdate(updateAddress)

	updateCam, err := DahuaCameraUpdate(dbCtx, update)
	assert.NoError(t, err)
	assert.Equal(t, updateAddress, updateCam.Address)

	// Get
	{
		getCam, err := DahuaCameraGet(dbCtx, updateCam.ID)
		assert.NoError(t, err)
		assert.Equal(t, updateCam, getCam)
	}

	// Delete
	{
		value, err := update.Value()
		assert.NoError(t, err)
		err = DahuaCameraDelete(dbCtx, value.ID)
		assert.NoError(t, err)
		err = DahuaCameraDelete(dbCtx, value.ID)
		assert.Error(t, err)
	}
}
