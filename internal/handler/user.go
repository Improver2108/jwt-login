package handler

import (
	base "github.com/improver2108/jwt-login/internal/pkg/handler"
	"github.com/improver2108/jwt-login/internal/service"
)

type UserHandler struct {
	base.Base[*service.UserService]
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{Base: base.Base[*service.UserService]{Service: svc}}
}
