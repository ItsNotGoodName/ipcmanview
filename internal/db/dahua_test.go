package db

import (
	"context"
	"testing"
	"time"

	"github.com/ItsNotGoodName/ipcmango/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestDahuaCameraCreate(t *testing.T) {
	context, close := Connect(context.Background())
	defer close()

	coreCam, err := core.DahuaCameraNew(core.DahuaCameraCreate{
		Address:  "localhost",
		Username: "Username",
		Password: "Password",
	})
	assert.NoError(t, err)

	dbCam, err := DahuaCameraCreate(context, coreCam)
	assert.NoError(t, err)

	assert.NotEqual(t, coreCam.ID, dbCam.ID, "should have new id")
	coreCam.ID = dbCam.ID

	assert.NotEqual(t, time.Time{}, dbCam.CreatedAt, "should not have default CreatedAt")
	coreCam.CreatedAt = dbCam.CreatedAt

	assert.Equal(t, coreCam, dbCam)
}

func TestDahuaCameraUpdate(t *testing.T) {
	context, close := Connect(context.Background())
	defer close()

	// Seed
	coreCam, _ := core.DahuaCameraNew(core.DahuaCameraCreate{
		Address:  "localhost",
		Username: "Username",
		Password: "Password",
	})
	dbCam, _ := DahuaCameraCreate(context, coreCam)

	address := "user"

	update := core.DahuaCameraUpdateNew(dbCam.ID)
	update.UpdateAddress(address)

	dbCam, err := DahuaCameraUpdate(context, update)
	assert.NoError(t, err)
	assert.Equal(t, address, dbCam.Address)
}
