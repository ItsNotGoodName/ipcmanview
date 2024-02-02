package core

import "context"

var actorCtxKey contextKey = contextKey("actor")

type ActorType string

const (
	ActorTypeSystem = "system"
	ActorTypeUser   = "user"
)

type Actor struct {
	Type   ActorType
	UserID int64
	Admin  bool
}

func WithUserActor(ctx context.Context, userID int64, admin bool) context.Context {
	return context.WithValue(ctx, actorCtxKey, Actor{
		Type:   ActorTypeUser,
		UserID: userID,
		Admin:  admin,
	})
}

func UseActor(ctx context.Context) Actor {
	actor, ok := ctx.Value(actorCtxKey).(Actor)
	if !ok {
		return Actor{
			Type:   ActorTypeSystem,
			UserID: 0,
			Admin:  true,
		}
	}
	return actor
}
