package repository

import (
	"context"
	"encoding/base64"

	"github.com/improver2108/jwt-login/db/sqlc"
	"github.com/improver2108/jwt-login/internal/model"
)

type UserRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(q *sqlc.Queries) *UserRepository {
	return &UserRepository{queries: q}
}

func (u *UserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	newUser := sqlc.CreateUserParams{
		Username:    user.Username,
		Email:       user.Email,
		PasswordHsh: base64.StdEncoding.EncodeToString(user.PasswordHash),
		Salt:        base64.StdEncoding.EncodeToString(user.Salt),
		Phone:       user.Phone,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
	}
	result, err := u.queries.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:        result.ID.String(),
		Username:  result.Username,
		AvatarUrl: result.AvatarUrl.String,
	}, nil
}

func (u *UserRepository) CheckUserExist(ctx context.Context, params *model.UserIdentifier) (bool, error) {
	newParams := sqlc.CheckUserExistParams{
		Email:    params.Email,
		Username: params.Username,
		Phone:    params.Phone,
	}
	exists, err := u.queries.CheckUserExist(ctx, newParams)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (u *UserRepository) GetLoginUser(ctx context.Context, email string) (*model.User, error) {

	user, err := u.queries.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, err
	}
	storedHash, _ := base64.StdEncoding.DecodeString(user.PasswordHsh)
	storedSalt, _ := base64.StdEncoding.DecodeString(user.Salt)
	return &model.User{
		ID:           user.ID.String(),
		AvatarUrl:    user.AvatarUrl.String,
		Username:     user.Username,
		PasswordHash: storedHash,
		Salt:         storedSalt,
	}, nil
}
