package auth

import (
	"context"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/config"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"golang.org/x/crypto/bcrypt"
)

func userFrom(v repo.User) _User {
	return _User{
		Email:    v.Email,
		Username: v.Username,
	}
}

type _User struct {
	Email    string `validate:"required,lte=128,email,excludes= "`
	Username string `validate:"gte=3,lte=64,excludes=@,excludes= "`
	Password string `validate:"gte=8,lte=64"`
}

func (u *_User) normalizeEmailAndUsername() {
	u.Email = strings.ToLower(u.Email)
	u.Username = strings.ToLower(u.Username)
}

func (u *_User) passwordHash() (string, error) {
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

func CreateUser(ctx context.Context, cfg config.Config, arg CreateUserParams) (int64, error) {
	actor := core.UseActor(ctx)
	if !actor.Admin && !cfg.EnableSignUp {
		return 0, core.ErrForbidden
	}

	model := _User{
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

	tx, err := app.DB.BeginTx(ctx, true)
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
			Time:  types.NewTime(now.Time),
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

func UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	if _, err := core.AssertAdminOrUser(ctx, arg.ID); err != nil {
		return err
	}

	dbModel, err := app.DB.C().AuthGetUser(ctx, arg.ID)
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

	_, err = app.DB.C().AuthPatchUser(ctx, repo.AuthPatchUserParams{
		Username:  core.NewNullString(model.Username),
		Email:     core.NewNullString(model.Email),
		UpdatedAt: types.NewTime(time.Now()),
		ID:        dbModel.ID,
	})
	return err
}

func DeleteUser(ctx context.Context, id int64) error {
	actor, err := core.AssertAdminOrUser(ctx, id)
	if err != nil {
		return err
	}
	if actor.Admin && actor.UserID == id {
		return core.ErrForbidden
	}
	return app.DB.C().DeleteUser(ctx, id)
}

type UpdateUserPasswordParams struct {
	UserID           int64
	OldPasswordSkip  bool
	OldPassword      string
	NewPassword      string
	CurrentSessionID int64
}

func UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	if _, err := core.AssertAdminOrUser(ctx, arg.UserID); err != nil {
		return err
	}

	dbModel, err := app.DB.C().AuthGetUser(ctx, arg.UserID)
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

	tx, err := app.DB.BeginTx(ctx, true)
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

	app.Hub.UserSecurityUpdated(bus.UserSecurityUpdated{
		UserID: dbModel.ID,
	})

	return nil
}

func UpdateUserUsername(ctx context.Context, userID int64, newUsername string) error {
	if _, err := core.AssertAdminOrUser(ctx, userID); err != nil {
		return err
	}

	dbModel, err := app.DB.C().AuthGetUser(ctx, userID)
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

	_, err = app.DB.C().AuthPatchUser(ctx, repo.AuthPatchUserParams{
		Username:  core.NewNullString(model.Username),
		UpdatedAt: types.NewTime(time.Now()),
		ID:        dbModel.ID,
	})
	return err
}

func UpdateUserDisabled(ctx context.Context, id int64, disable bool) error {
	actor, err := core.AssertAdmin(ctx)
	if err != nil {
		return err
	}
	if actor.UserID == id {
		return core.ErrForbidden
	}

	_, err = app.DB.C().AuthUpdateUserDisabledAt(ctx, repo.AuthUpdateUserDisabledAtParams{
		DisabledAt: types.NullTime{
			Time:  types.NewTime(time.Now()),
			Valid: disable,
		},
		ID: id,
	})
	if err != nil {
		return err
	}

	app.Hub.UserSecurityUpdated(bus.UserSecurityUpdated{
		UserID: id,
	})

	return nil
}

func UpdateUserAdmin(ctx context.Context, id int64, admin bool) error {
	actor, err := core.AssertAdmin(ctx)
	if err != nil {
		return err
	}
	if actor.UserID == id {
		return core.ErrForbidden
	}

	if admin {
		_, err := app.DB.C().AuthUpsertAdmin(ctx, repo.AuthUpsertAdminParams{
			UserID:    id,
			CreatedAt: types.NewTime(time.Now()),
		})
		if err != nil {
			return err
		}
	} else {
		err := app.DB.C().AuthDeleteAdmin(ctx, id)
		if err != nil {
			return err
		}
	}

	app.Hub.UserSecurityUpdated(bus.UserSecurityUpdated{
		UserID: id,
	})

	return nil
}

func CheckUserPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
