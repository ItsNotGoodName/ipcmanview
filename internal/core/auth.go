package core

import (
	"context"
	"fmt"
)

func userOrAdmin(actor Actor, userID int64) error {
	if actor.Admin || actor.UserID == userID {
		return nil
	}
	return fmt.Errorf("not user or admin")
}

func UserOrAdmin(ctx context.Context, userID int64) error {
	actor := UseActor(ctx)
	return userOrAdmin(actor, userID)
}

func UserOrAdminActor(ctx context.Context, userID int64) (Actor, error) {
	actor := UseActor(ctx)
	return actor, userOrAdmin(actor, userID)
}

func Admin(ctx context.Context) error {
	if UseActor(ctx).Admin {
		return nil
	}
	return fmt.Errorf("not admin")
}
