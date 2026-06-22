package customerrors

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exits")
	ErrUserNotFound      = errors.New("user not found")
	ErrCreatingUser      = errors.New("error creating user")
	ErrCheckinUserExists = errors.New("error checking users exist")
)
