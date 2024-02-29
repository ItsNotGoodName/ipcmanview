package auth

import (
	"context"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/config"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"golang.org/x/crypto/bcrypt"
)

func userFrom(v repo.User) user {
	return user{
		Email:    v.Email,
		Username: v.Username,
	}
}

type user struct {
	Email    string `validate:"required,lte=128,email,excludes= "`
	Username string `validate:"gte=3,lte=64,excludes=@,excludes= "`
	Password string `validate:"gte=8,lte=64"`
}

func (u *user) normalizeEmailAndUsername() {
	u.Email = strings.ToLower(u.Email)
	u.Username = strings.ToLower(u.Username)
}

func (u *user) passwordHash() (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

type CreateUserParams struct {
	Email    string
	Username string
	Password string
	Admin    bool
	Disabled bool
}

func CreateUser(ctx context.Context, cfg config.Config, db sqlite.DB, arg CreateUserParams) (int64, error) {
	actor := core.UseActor(ctx)
	if !actor.Admin && !cfg.EnableSignUp {
		return 0, core.ErrForbidden
	}

	model := user{
		Email:    arg.Email,
		Username: arg.Username,
		Password: arg.Password,
	}
	model.normalizeEmailAndUsername()

	if err := core.ValidateStruct(ctx, model); err != nil {
		return 0, err
	}

	password, err := model.passwordHash()
	if err != nil {
		return 0, err
	}

	var (
		disabled bool
		admin    bool
	)
	if actor.Admin {
		disabled = arg.Disabled
		admin = arg.Admin
	}

	tx, err := db.BeginTx(ctx, true)
	if err != nil {
		return 0, nil
	}
	defer tx.Rollback()

	now := types.NewTime(time.Now())
	id, err := tx.C().AuthCreateUser(ctx, repo.AuthCreateUserParams{
		Email:     arg.Email,
		Username:  arg.Username,
		Password:  password,
		CreatedAt: now,
		UpdatedAt: now,
		DisabledAt: types.NullTime{
			Time:  now.Time,
			Valid: disabled,
		},
	})
	if err != nil {
		return 0, err
	}

	if admin {
		tx.C().AuthUpsertAdmin(ctx, repo.AuthUpsertAdminParams{
			UserID:    id,
			CreatedAt: now,
		})
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

type UpdateUserParams struct {
	ID       int64
	Email    string
	Username string
}

func UpdateUser(ctx context.Context, db sqlite.DB, arg UpdateUserParams) error {
	if _, err := core.AssertAdminOrUser(ctx, arg.ID); err != nil {
		return err
	}

	dbModel, err := db.C().AuthGetUser(ctx, arg.ID)
	if err != nil {
		return err
	}

	model := userFrom(dbModel)

	// Mutate
	model.Email = arg.Email
	model.Username = arg.Username
	model.normalizeEmailAndUsername()

	if err := core.ValidateStructPartial(ctx, model, "Email", "Username"); err != nil {
		return err
	}

	_, err = db.C().AuthPatchUser(ctx, repo.AuthPatchUserParams{
		Username:  core.NewNullString(model.Username),
		Email:     core.NewNullString(model.Email),
		UpdatedAt: types.NewTime(time.Now()),
		ID:        dbModel.ID,
	})
	return err
}

func DeleteUser(ctx context.Context, db sqlite.DB, id int64) error {
	actor, err := core.AssertAdminOrUser(ctx, id)
	if err != nil {
		return err
	}
	if actor.Admin && actor.UserID == id {
		return core.ErrForbidden
	}
	return db.C().DeleteUser(ctx, id)
}

type UpdateUserPasswordParams struct {
	UserID           int64
	OldPasswordSkip  bool
	OldPassword      string
	NewPassword      string
	CurrentSessionID int64
}

func UpdateUserPassword(ctx context.Context, db sqlite.DB, bus *event.Bus, arg UpdateUserPasswordParams) error {
	if _, err := core.AssertAdminOrUser(ctx, arg.UserID); err != nil {
		return err
	}

	dbModel, err := db.C().AuthGetUser(ctx, arg.UserID)
	if err != nil {
		return err
	}
	model := userFrom(dbModel)

	if !arg.OldPasswordSkip {
		if err := CheckUserPassword(dbModel.Password, arg.OldPassword); err != nil {
			return core.NewFieldError("OldPassword", err.Error())
		}
	}

	// Mutate
	model.Password = arg.NewPassword

	if err := core.ValidateStructPartial(ctx, model, "Password"); err != nil {
		return err
	}

	password, err := model.passwordHash()
	if err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx, true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.C().AuthPatchUser(ctx, repo.AuthPatchUserParams{
		Password:  core.NewNullString(password),
		UpdatedAt: types.NewTime(time.Now()),
		ID:        dbModel.ID,
	})
	if err != nil {
		return err
	}

	err = tx.C().AuthDeleteUserSessionForUserAndNotSession(ctx, repo.AuthDeleteUserSessionForUserAndNotSessionParams{
		UserID: dbModel.ID,
		ID:     arg.CurrentSessionID,
	})
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	bus.UserSecurityUpdated(event.UserSecurityUpdated{
		UserID: dbModel.ID,
	})

	return nil
}

func UpdateUserUsername(ctx context.Context, db sqlite.DB, userID int64, newUsername string) error {
	if _, err := core.AssertAdminOrUser(ctx, userID); err != nil {
		return err
	}

	dbModel, err := db.C().AuthGetUser(ctx, userID)
	if err != nil {
		return err
	}
	model := userFrom(dbModel)

	// Mutate
	model.Username = newUsername
	model.normalizeEmailAndUsername()

	if err := core.ValidateStructPartial(ctx, model, "Username"); err != nil {
		return err
	}

	_, err = db.C().AuthPatchUser(ctx, repo.AuthPatchUserParams{
		Username:  core.NewNullString(model.Username),
		UpdatedAt: types.NewTime(time.Now()),
		ID:        dbModel.ID,
	})
	return err
}

func UpdateUserDisabled(ctx context.Context, db sqlite.DB, bus *event.Bus, id int64, disable bool) error {
	actor, err := core.AssertAdmin(ctx)
	if err != nil {
		return err
	}
	if actor.UserID == id {
		return core.ErrForbidden
	}

	if disable {
		_, err := db.C().AuthUpdateUserDisabledAt(ctx, repo.AuthUpdateUserDisabledAtParams{
			DisabledAt: types.NewNullTime(time.Now()),
			ID:         id,
		})
		if err != nil {
			return err
		}
	} else {
		_, err := db.C().AuthUpdateUserDisabledAt(ctx, repo.AuthUpdateUserDisabledAtParams{
			DisabledAt: types.NullTime{},
			ID:         id,
		})
		if err != nil {
			return err
		}
	}

	bus.UserSecurityUpdated(event.UserSecurityUpdated{
		UserID: id,
	})

	return nil
}

func UpdateUserAdmin(ctx context.Context, db sqlite.DB, bus *event.Bus, id int64, admin bool) error {
	actor, err := core.AssertAdmin(ctx)
	if err != nil {
		return err
	}
	if actor.UserID == id {
		return core.ErrForbidden
	}

	if admin {
		_, err := db.C().AuthUpsertAdmin(ctx, repo.AuthUpsertAdminParams{
			UserID:    id,
			CreatedAt: types.NewTime(time.Now()),
		})
		if err != nil {
			return err
		}
	} else {
		err := db.C().AuthDeleteAdmin(ctx, id)
		if err != nil {
			return err
		}
	}

	bus.UserSecurityUpdated(event.UserSecurityUpdated{
		UserID: id,
	})

	return nil
}

func CheckUserPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GetUserByUsernameOrEmail(ctx context.Context, db sqlite.DB, usernameOrEmail string) (repo.User, error) {
	return db.C().AuthGetUserByUsernameOrEmail(ctx, strings.ToLower(strings.TrimSpace(usernameOrEmail)))
}
