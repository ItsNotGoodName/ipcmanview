package auth

import (
	"context"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/internal/validate"
	"golang.org/x/crypto/bcrypt"
)

func normalizeUser(arg *models.User) {
	arg.Email = strings.ToLower(arg.Email)
	arg.Username = strings.ToLower(arg.Username)
}

func hashUserPassword(arg *models.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(arg.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	arg.Password = string(hash)

	return nil
}

func CreateUser(ctx context.Context, db repo.DB, arg models.User) (int64, error) {
	normalizeUser(&arg)

	if err := validate.Validate.Struct(arg); err != nil {
		return 0, err
	}

	if err := hashUserPassword(&arg); err != nil {
		return 0, err
	}

	now := types.NewTime(time.Now())
	return db.CreateUser(ctx, repo.CreateUserParams{
		Email:     arg.Email,
		Username:  arg.Username,
		Password:  arg.Password,
		CreatedAt: now,
		UpdatedAt: now,
	})
}

func UpdateUser(ctx context.Context, db repo.DB, arg models.User, newPassword string) (int64, error) {
	if newPassword != "" {
		arg.Password = newPassword
	}

	normalizeUser(&arg)

	if err := validate.Validate.Struct(arg); err != nil {
		return 0, err
	}

	if newPassword != "" {
		if err := hashUserPassword(&arg); err != nil {
			return 0, err
		}
	}

	return db.UpdateUser(ctx, repo.UpdateUserParams{
		Email:     arg.Email,
		Username:  arg.Username,
		Password:  arg.Password,
		UpdatedAt: types.NewTime(time.Now()),
		ID:        arg.ID,
	})
}

func CheckUserPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func UpdateUserDisable(ctx context.Context, db repo.DB, userID int64, disable bool) error {
	if disable {
		_, err := db.UpdateUserDisabledAt(ctx, repo.UpdateUserDisabledAtParams{
			DisabledAt: types.NewNullTime(time.Now()),
			ID:         userID,
		})
		return err
	}
	_, err := db.UpdateUserDisabledAt(ctx, repo.UpdateUserDisabledAtParams{
		DisabledAt: types.NullTime{},
		ID:         userID,
	})
	return err
}

func UpdateUserAdmin(ctx context.Context, db repo.DB, userId int64, admin bool) error {
	return nil
}
