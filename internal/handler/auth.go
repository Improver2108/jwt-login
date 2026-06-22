package handler

import (
	"errors"
	"net/http"

	"github.com/improver2108/jwt-login/internal/auth"
	"github.com/improver2108/jwt-login/internal/cache"
	customerrors "github.com/improver2108/jwt-login/internal/errors"
	"github.com/improver2108/jwt-login/internal/model"
	base "github.com/improver2108/jwt-login/internal/pkg/handler"
	"github.com/improver2108/jwt-login/internal/service"
)

type AuthHandler struct {
	base.Base[*service.AuthService]
	cache *cache.JTICache
}

func NewAuthHandler(svc *service.AuthService, c *cache.JTICache) *AuthHandler {
	return &AuthHandler{
		Base: base.Base[*service.AuthService]{Service: svc}, cache: c,
	}
}

func (a *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterUserRequest
	if ok := a.Decode(w, r, &req); !ok {
		return
	}
	tokens, err := a.Service.RegisterService(r.Context(), &req, a.cache)
	if err != nil {
		switch {
		case errors.Is(err, customerrors.ErrCredentialValidation):
			a.BadRequest(w, "Enter valid details")
		case errors.Is(err, customerrors.ErrPasswordsMismatch):
			a.BadRequest(w, "passwords don't match")
		case errors.Is(err, customerrors.ErrUserAlreadyExists):
			a.Conflict(w, "user already exists")
		default:
			a.InternalError(w, "something went wrong")
		}
		return
	}

	auth.SetAuthCookie(w, tokens)

	a.Created(w, "User Created")
}

func (a *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req model.LoginUserRequest
	if ok := a.Decode(w, r, &req); !ok {
		return
	}
	tokens, err := a.Service.LoginService(r.Context(), &req, a.cache)
	if err != nil {
		switch {
		case errors.Is(err, customerrors.ErrCredentialValidation):
			a.BadRequest(w, "Please fill proper details")
		case errors.Is(err, customerrors.ErrUserNotFound), errors.Is(err, customerrors.ErrWrongPassword):
			a.BadRequest(w, "Invlaid Credentials")
		default:
			a.InternalError(w, "Something went wrong")
		}
		return
	}
	auth.SetAuthCookie(w, tokens)
	a.OK(w, "Successfully login")

}

func (a *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	accessStr, err := auth.GetCookie(r, "access_token")
	if err != nil {
		a.Unauthorized(w, "Access token missing")
		return
	}
	refreshStr, err := auth.GetCookie(r, "refresh_token")
	if err != nil {
		a.Unauthorized(w, "Refresh token missing")
		return
	}

	if err := a.Service.LogoutService(r.Context(), refreshStr, accessStr, a.cache); err != nil {
		a.InternalError(w, "error removing tokens")
		return
	}

	auth.ClearAuthCookie(w)
	a.OK(w, "Logout successfully!")
}

func (a *AuthHandler) TokenRefreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshStr, err := auth.GetCookie(r, "refresh_token")
	if err != nil {
		a.Unauthorized(w, "Refresh token missing")
		return
	}
	token, err := a.Service.RefreshTokensService(r.Context(), refreshStr, a.cache)
	if err != nil {
		a.InternalError(w, "Something went wrong")
		return
	}

	auth.SetAuthCookie(w, token)
	a.OK(w, "token refreshed")
}
