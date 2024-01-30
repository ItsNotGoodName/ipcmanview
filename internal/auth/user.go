package auth

import (
	"context"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/internal/validate"
	"golang.org/x/crypto/bcrypt"
)

func userFromDB(v repo.User) user {
	return user{
		Email:    v.Email,
		Username: v.Username,
	}
}

type user struct {
	Email    string `validate:"required,lte=128,email,excludes= "`
	Username string `validate:"gte=3,lte=64,excludes=@,excludes= "`
	Password string `validate:"gte=8"`
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
}

func CreateUser(ctx context.Context, db repo.DB, arg CreateUserParams) (int64, error) {
	model := user{
		Email:    arg.Email,
		Username: arg.Username,
		Password: arg.Password,
	}

	if err := validate.Validate.Struct(model); err != nil {
		return 0, err
	}

	password, err := model.passwordHash()
	if err != nil {
		return 0, err
	}

	now := types.NewTime(time.Now())
	return db.AuthCreateUser(ctx, repo.AuthCreateUserParams{
		Email:     arg.Email,
		Username:  arg.Username,
		Password:  password,
		CreatedAt: now,
		UpdatedAt: now,
	})
}

func UpdateUserPassword(ctx context.Context, db repo.DB, dbModel repo.User, newPassword string) error {
	model := userFromDB(dbModel)

	// Mutate
	model.Password = newPassword

	if err := validate.Validate.StructPartial(model, "Password"); err != nil {
		return err
	}

	password, err := model.passwordHash()
	if err != nil {
		return err
	}

	_, err = db.AuthUpdateUser(ctx, repo.AuthUpdateUserParams{
		Password:  core.NewNullString(password),
		UpdatedAt: types.NewTime(time.Now()),
		ID:        dbModel.ID,
	})
	return err
}

func UpdateUserUsername(ctx context.Context, db repo.DB, dbModel repo.User, newUsername string) error {
	model := userFromDB(dbModel)

	// Mutate
	model.Username = newUsername
	model.normalizeEmailAndUsername()

	if err := validate.Validate.StructPartial(model, "Username"); err != nil {
		return err
	}

	_, err := db.AuthUpdateUser(ctx, repo.AuthUpdateUserParams{
		Username:  core.NewNullString(model.Username),
		UpdatedAt: types.NewTime(time.Now()),
		ID:        dbModel.ID,
	})
	return err
}

func UpdateUserDisable(ctx context.Context, db repo.DB, userID int64, disable bool) error {
	if disable {
		_, err := db.AuthUpdateUserDisabledAt(ctx, repo.AuthUpdateUserDisabledAtParams{
			DisabledAt: types.NewNullTime(time.Now()),
			ID:         userID,
		})
		return err
	}
	_, err := db.AuthUpdateUserDisabledAt(ctx, repo.AuthUpdateUserDisabledAtParams{
		DisabledAt: types.NullTime{},
		ID:         userID,
	})
	return err
}

func UpdateUserAdmin(ctx context.Context, db repo.DB, userID int64, admin bool) error {
	if admin {
		_, err := db.AuthUpsertAdmin(ctx, repo.AuthUpsertAdminParams{
			UserID:    userID,
			CreatedAt: types.NewTime(time.Now()),
		})
		return err
	}
	return db.AuthDeleteAdmin(ctx, userID)
}

func CheckUserPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
