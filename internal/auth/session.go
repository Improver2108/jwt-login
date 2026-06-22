package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/improver2108/jwt-login/internal/cache"
)

func PersistSession(ctx context.Context, c *cache.JTICache, t *Tokens) error {
	if err := c.SetJTI(ctx, fmt.Sprintf("access:%s", t.JTIAcc), t.UserID, time.Until(t.ExpAcc)); err != nil {
		return err
	}
	if err := c.SetJTI(ctx, fmt.Sprintf("refresh:%s", t.JTIRef), t.UserID, time.Until(t.ExpRef)); err != nil {
		return err
	}
	return nil
}
