package dahua

//
// import (
// 	"context"
// 	"testing"
//
// 	"github.com/ItsNotGoodName/ipcmango/internal/core"
// 	"github.com/ItsNotGoodName/ipcmango/internal/dbtest"
// 	"github.com/stretchr/testify/assert"
// )
//
// func TestStore(t *testing.T) {
// 	ctx := context.Background()
// 	db, close := dbtest.Connect(ctx)
// 	defer close()
//
// 	// Seed
// 	cam1, err := DB.CameraCreate(ctx, db, core.DahuaCamera{Address: "127.0.0.1"})
// 	assert.NoError(t, err)
//
// 	cam2, err := DB.CameraCreate(ctx, db, core.DahuaCamera{Address: "localhost"})
// 	assert.NoError(t, err)
//
// 	store := NewStore()
// 	assert.Equal(t, []StoreActorHandle{}, store.actors, "should create store")
//
// 	// Create camera 1
// 	actor1, err := store.GetOrCreate(ctx, db, cam1.ID)
// 	assert.NoError(t, err)
// 	assert.Equal(t, []StoreActorHandle{actor1}, store.actors)
//
// 	// Create camera 2
// 	actor2, err := store.GetOrCreate(ctx, db, cam2.ID)
// 	assert.NoError(t, err)
// 	assert.Equal(t, []StoreActorHandle{actor1, actor2}, store.actors)
//
// 	// Get camera 2
// 	{
// 		_, err = store.GetOrCreate(ctx, db, cam2.ID)
// 		assert.NoError(t, err)
// 		assert.Equal(t, []StoreActorHandle{actor1, actor2}, store.actors)
// 	}
//
// 	// Delete camera 1
// 	{
// 		store.Delete(ctx, cam1.ID)
// 		assert.Equal(t, []StoreActorHandle{actor2}, store.actors)
// 		_, ok := <-actor1.doneC
// 		assert.False(t, ok)
// 		actor1, err = store.GetOrCreate(ctx, db, cam1.ID)
// 		assert.NoError(t, err)
// 		assert.Equal(t, []StoreActorHandle{actor2, actor1}, store.actors)
// 	}
//
// 	// Update camera 1
// 	{
// 		// Update database
// 		update := core.NewDahuaCameraUpdate(cam1.ID).AddressUpdate("hi")
// 		assert.NoError(t, err)
//
// 		cam1Updated, err := DB.CameraUpdate(ctx, db, update)
// 		assert.NoError(t, err)
//
// 		// Update store
// 		actorCam1, err := store.GetOrCreate(ctx, db, cam1.ID)
// 		assert.NoError(t, err)
// 		assert.True(t, cam1Updated.Equal(actorCam1.cam))
// 	}
// }
