package service

import (
	"github.com/improver2108/jwt-login/internal/repository"
)

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repository: *repo}
}
