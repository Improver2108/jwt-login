package service

import (
	"context"

	"github.com/improver2108/jwt-login/internal/auth"
	"github.com/improver2108/jwt-login/internal/cache"
	customerrors "github.com/improver2108/jwt-login/internal/errors"
	"github.com/improver2108/jwt-login/internal/model"
	"github.com/improver2108/jwt-login/internal/pkg/utility"
	"github.com/improver2108/jwt-login/internal/repository"
)

type AuthService struct {
	repository *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repository: repo}
}

func (a *AuthService) RegisterService(ctx context.Context, req *model.RegisterUserRequest, c *cache.JTICache) (*auth.Tokens, error) {
	if req.Email == "" || req.Password == "" || req.Phone == "" || req.FirstName == "" || req.Username == "" {
		return nil, customerrors.ErrCredentialValidation
	}
	if req.Password != req.ConfirmPassword {
		return nil, customerrors.ErrPasswordsMismatch
	}
	params := model.UserIdentifier{
		Email:    req.Email,
		Username: req.Username,
		Phone:    req.Phone,
	}
	if exists, err := a.repository.CheckUserExist(ctx, &params); err != nil {
		return nil, err
	} else if exists {
		return nil, customerrors.ErrUserAlreadyExists
	}

	hashedPassword, salt, err := utility.HashPassword(req.Password)

	if err != nil {
		return nil, customerrors.ErrPasswordHash
	}

	newUser := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		Phone:        req.Phone,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PasswordHash: hashedPassword,
		Salt:         salt,
		AvatarUrl:    req.AvatarUrl,
		Bio:          req.Bio,
	}

	createdUser, err := a.repository.CreateUser(ctx, newUser)

	if err != nil {
		return nil, customerrors.ErrCreatingUser
	}

	tokens, err := auth.IssueTokens(createdUser.ID)

	if err != nil {
		return nil, customerrors.ErrIssuingTokens
	}

	if err := auth.PersistSession(ctx, c, tokens); err != nil {
		return nil, customerrors.ErrTokenPersist
	}

	return tokens, nil
}

func (a *AuthService) LoginService(ctx context.Context, req *model.LoginUserRequest, c *cache.JTICache) (*auth.Tokens, error) {

	if req.Email == "" || req.Password == "" {
		return nil, customerrors.ErrCredentialValidation
	}

	userDetails, err := a.repository.GetLoginUser(ctx, req.Email)
	if err != nil {
		return nil, customerrors.ErrUserNotFound
	}

	isPasswordSame := utility.VerifyPassword(req.Password, userDetails.Salt, userDetails.PasswordHash)

	if !isPasswordSame {
		return nil, customerrors.ErrWrongPassword
	}

	tokens, err := auth.IssueTokens(userDetails.ID)

	if err != nil {
		return nil, customerrors.ErrIssuingTokens
	}

	if err := auth.PersistSession(ctx, c, tokens); err != nil {
		return nil, customerrors.ErrTokenPersist
	}
	return tokens, nil
}

func (a *AuthService) LogoutService(ctx context.Context, refreshStr, accessStr string, c *cache.JTICache) error {
	if claims, err := auth.ParseRefresh(refreshStr); err != nil {
		_ = c.DelJTI(ctx, "refresh:"+claims.ID)
	}
	if claims, err := auth.ParseAccess(accessStr); err != nil {
		_ = c.DelJTI(ctx, "access:"+claims.ID)
	}
	return nil
}

func (a *AuthService) RefreshTokensService(ctx context.Context, refreshStr string, c *cache.JTICache) (*auth.Tokens, error) {
	claims, err := auth.ParseRefresh(refreshStr)
	if err != nil {
		return nil, err
	}
	_, err = c.GetJTI(ctx, "refresh:"+claims.ID)
	if err != nil {
		return nil, err
	}
	_ = c.DelJTI(ctx, "refresh:"+claims.ID)

	token, err := auth.IssueTokens(claims.Subject)

	if err != nil {
		return nil, customerrors.ErrIssuingTokens
	}

	if err := auth.PersistSession(ctx, c, token); err != nil {
		return nil, customerrors.ErrTokenPersist
	}

	return token, nil

}
