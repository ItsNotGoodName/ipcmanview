package core

import (
	"context"
	"fmt"
)

var actorCtxKey contextKey = contextKey("actor")

type ActorType string

const (
	ActorTypeSystem = "system"
	ActorTypeUser   = "user"
	ActorTypePublic = "public"
)

func newSystemActor() Actor {
	return Actor{
		Type:   ActorTypeSystem,
		UserID: 0,
		Admin:  true,
	}
}

type Actor struct {
	Type   ActorType
	UserID int64
	Admin  bool
}

// WithSystemActor sets actor to system.
func WithSystemActor(ctx context.Context) context.Context {
	return context.WithValue(ctx, actorCtxKey, newSystemActor())
}

// WithPublicActor sets actor to public.
func WithPublicActor(ctx context.Context) context.Context {
	return context.WithValue(ctx, actorCtxKey, Actor{
		Type:   ActorTypePublic,
		UserID: 0,
		Admin:  false,
	})
}

// WithUserActor sets actor to user.
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
		return newSystemActor()
	}
	return actor
}

func AssertAdminOrUser(ctx context.Context, userID int64) (Actor, error) {
	actor := UseActor(ctx)
	if actor.Admin || actor.UserID == userID {
		return actor, nil
	}
	return actor, fmt.Errorf("%w: not admin or user", ErrForbidden)
}

func AssertAdmin(ctx context.Context) (Actor, error) {
	actor := UseActor(ctx)
	if actor.Admin {
		return actor, nil
	}
	return actor, fmt.Errorf("%w: not admin", ErrForbidden)
}
