package dahua

import (
	"context"
	"testing"

	"github.com/ItsNotGoodName/ipcmango/internal/core"
	"github.com/ItsNotGoodName/ipcmango/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	context, close := db.TestConnect(context.Background())
	defer close()

	// Seed
	cam1, err := db.DahuaCameraCreate(context, core.DahuaCamera{Address: "127.0.0.1"})
	assert.NoError(t, err)

	cam2, err := db.DahuaCameraCreate(context, core.DahuaCamera{Address: "localhost"})
	assert.NoError(t, err)

	store := NewStore()
	assert.Equal(t, []StoreActor{}, store.actors, "should create store")

	// Create camera 1
	actor1, err := store.GetOrCreate(context, cam1.ID)
	assert.NoError(t, err)
	assert.Equal(t, []StoreActor{actor1}, store.actors)

	// Create camera 2
	actor2, err := store.GetOrCreate(context, cam2.ID)
	assert.NoError(t, err)
	assert.Equal(t, []StoreActor{actor1, actor2}, store.actors)

	// Get camera 2
	{
		_, err = store.GetOrCreate(context, cam2.ID)
		assert.NoError(t, err)
		assert.Equal(t, []StoreActor{actor1, actor2}, store.actors)
	}

	// Delete camera 1
	{
		store.Delete(context, cam1.ID)
		assert.Equal(t, []StoreActor{actor2}, store.actors)
		_, ok := <-actor1.doneC
		assert.False(t, ok)
		actor1, err = store.GetOrCreate(context, cam1.ID)
		assert.NoError(t, err)
		assert.Equal(t, []StoreActor{actor2, actor1}, store.actors)
	}

	// Update camera 1
	{
		// Update database
		update := core.NewDahuaCameraUpdate(cam1.ID).AddressUpdate("hi")
		assert.NoError(t, err)

		cam1Updated, err := db.DahuaCameraUpdate(context, update)
		assert.NoError(t, err)

		// Update store
		actorCam1, err := store.GetOrCreate(context, cam1.ID)
		assert.NoError(t, err)
		assert.True(t, cam1Updated.Equal(actorCam1.Camera))
	}
}
