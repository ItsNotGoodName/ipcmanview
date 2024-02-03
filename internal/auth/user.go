package auth

import (
	"context"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
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

func CreateUser(ctx context.Context, db repo.DB, arg CreateUserParams) (int64, error) {
	model := user{
		Email:    arg.Email,
		Username: arg.Username,
		Password: arg.Password,
	}
	model.normalizeEmailAndUsername()

	if err := core.Validate.Struct(model); err != nil {
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
	if core.UseActor(ctx).Admin {
		disabled = arg.Disabled
		admin = arg.Admin
	}

	tx, err := db.BeginTx(ctx, true)
	if err != nil {
		return 0, nil
	}
	defer tx.Rollback()

	now := types.NewTime(time.Now())
	id, err := tx.AuthCreateUser(ctx, repo.AuthCreateUserParams{
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
		tx.AuthUpsertAdmin(ctx, repo.AuthUpsertAdminParams{
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
	Email    string
	Username string
}

func UpdateUser(ctx context.Context, db repo.DB, dbModel repo.User, arg UpdateUserParams) error {
	if err := core.UserOrAdmin(ctx, dbModel.ID); err != nil {
		return err
	}

	model := userFrom(dbModel)

	// Mutate
	model.Email = arg.Email
	model.Username = arg.Username
	model.normalizeEmailAndUsername()

	if err := core.Validate.StructPartial(model, "Email", "Username"); err != nil {
		return err
	}

	_, err := db.AuthPatchUser(ctx, repo.AuthPatchUserParams{
		Username:  core.NewNullString(model.Username),
		Email:     core.NewNullString(model.Email),
		UpdatedAt: types.NewTime(time.Now()),
		ID:        dbModel.ID,
	})
	return err
}

func DeleteUser(ctx context.Context, db repo.DB, id int64) error {
	if err := core.UserOrAdmin(ctx, id); err != nil {
		return err
	}
	return db.DeleteUser(ctx, id)
}

type UpdateUserPasswordParams struct {
	NewPassword      string
	CurrentSessionID int64
}

func UpdateUserPassword(ctx context.Context, db repo.DB, dbModel repo.User, arg UpdateUserPasswordParams) error {
	if err := core.UserOrAdmin(ctx, dbModel.ID); err != nil {
		return err
	}

	model := userFrom(dbModel)

	// Mutate
	model.Password = arg.NewPassword

	if err := core.Validate.StructPartial(model, "Password"); err != nil {
		return err
	}

	password, err := model.passwordHash()
	if err != nil {
		return err
	}

	_, err = db.AuthPatchUser(ctx, repo.AuthPatchUserParams{
		Password:  core.NewNullString(password),
		UpdatedAt: types.NewTime(time.Now()),
		ID:        dbModel.ID,
	})
	if err != nil {
		return err
	}

	err = db.AuthDeleteUserSessionForUserAndNotSession(ctx, repo.AuthDeleteUserSessionForUserAndNotSessionParams{
		UserID: dbModel.ID,
		ID:     arg.CurrentSessionID,
	})
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserUsername(ctx context.Context, db repo.DB, dbModel repo.User, newUsername string) error {
	if err := core.UserOrAdmin(ctx, dbModel.ID); err != nil {
		return err
	}

	model := userFrom(dbModel)

	// Mutate
	model.Username = newUsername
	model.normalizeEmailAndUsername()

	if err := core.Validate.StructPartial(model, "Username"); err != nil {
		return err
	}

	_, err := db.AuthPatchUser(ctx, repo.AuthPatchUserParams{
		Username:  core.NewNullString(model.Username),
		UpdatedAt: types.NewTime(time.Now()),
		ID:        dbModel.ID,
	})
	return err
}

func UpdateUserDisabled(ctx context.Context, db repo.DB, id int64, disable bool) error {
	if err := core.Admin(ctx); err != nil {
		return err
	}

	if disable {
		_, err := db.AuthUpdateUserDisabledAt(ctx, repo.AuthUpdateUserDisabledAtParams{
			DisabledAt: types.NewNullTime(time.Now()),
			ID:         id,
		})
		return err
	}
	_, err := db.AuthUpdateUserDisabledAt(ctx, repo.AuthUpdateUserDisabledAtParams{
		DisabledAt: types.NullTime{},
		ID:         id,
	})
	return err
}

func UpdateUserAdmin(ctx context.Context, db repo.DB, id int64, admin bool) error {
	if err := core.Admin(ctx); err != nil {
		return err
	}

	if admin {
		_, err := db.AuthUpsertAdmin(ctx, repo.AuthUpsertAdminParams{
			UserID:    id,
			CreatedAt: types.NewTime(time.Now()),
		})
		return err
	}
	return db.AuthDeleteAdmin(ctx, id)
}

func CheckUserPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GetUserByUsernameOrEmail(ctx context.Context, db repo.DB, usernameOrEmail string) (repo.User, error) {
	return db.AuthGetUserByUsernameOrEmail(ctx, strings.ToLower(strings.TrimSpace(usernameOrEmail)))
}
