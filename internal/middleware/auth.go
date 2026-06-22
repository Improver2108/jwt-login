package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/improver2108/jwt-login/internal/auth"
	"github.com/improver2108/jwt-login/internal/cache"
)

func bearerFromHeader(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if after, ok := strings.CutPrefix(h, "Bearer "); ok {
		return after
	}
	return ""
}

func AuthMiddleware(c *cache.JTICache) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenStr string
			if cookie, err := r.Cookie("access_token"); err != nil || cookie.Value == "" {
				tokenStr = bearerFromHeader(r)
			} else {
				tokenStr = cookie.Value
			}

			if tokenStr == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			claims, err := auth.ParseAccess(tokenStr)

			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			if _, err := c.GetJTI(ctx, "access:"+claims.ID); err != nil {
				http.Error(w, "token revoked", http.StatusUnauthorized)
				return
			}

			newCtx := context.WithValue(ctx, userKey, &authUser{
				ID: claims.Subject,
			})

			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}
