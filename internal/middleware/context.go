package middleware

import (
	"context"
	"errors"
)

type contextKey string

const userKey contextKey = "user"

type authUser struct {
	ID string
}

func UserFromContext(ctx context.Context) (*authUser, error) {
	u, ok := ctx.Value(userKey).(*authUser)
	if !ok || u == nil {
		return nil, errors.New("no authenticated user in context")
	}
	return u, nil
}
