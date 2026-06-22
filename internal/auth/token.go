package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.RegisteredClaims
}

type Tokens struct {
	Access   string
	Refresh  string
	JTIAcc   string
	JTIRef   string
	ExpAcc   time.Time
	ExpRef   time.Time
	UserID   string
	Issuer   string
	Audience string
}

func IssueTokens(userID string) (*Tokens, error) {
	now := time.Now()
	t := &Tokens{
		UserID:   userID,
		JTIAcc:   uuid.NewString(),
		JTIRef:   uuid.NewString(),
		ExpAcc:   now.Add(15 * time.Minute),
		ExpRef:   now.Add(7 * 24 * time.Hour),
		Issuer:   "jwt-login-app",
		Audience: "jwt-login-client",
	}
	acc := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ID:        t.JTIAcc,
			Issuer:    t.Issuer,
			Audience:  jwt.ClaimStrings{t.Audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(t.ExpAcc),
		},
	})

	ref := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ID:        t.JTIRef,
			Issuer:    t.Issuer,
			Audience:  jwt.ClaimStrings{t.Audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(t.ExpRef),
		},
	})

	var err error

	t.Access, err = acc.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	t.Refresh, err = ref.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return t, nil
}

func ParseAccess(tokenStr string) (*Claims, error) {
	var secret = os.Getenv("ACCESS_SECRET")
	return parseWithSecret(tokenStr, secret)
}

func ParseRefresh(tokenStr string) (*Claims, error) {
	var secret = os.Getenv("REFRESH_SECRET")
	return parseWithSecret(tokenStr, secret)
}

func parseWithSecret(tokenStr, secret string) (*Claims, error) {
	if secret == "" {
		return nil, errors.New("jwt secret not configured")
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	token, err := parser.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
