package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`

	jwt.RegisteredClaims
}

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenManager(cfg Config) TokenManager {
	return &JWTManager{
		secret: []byte(cfg.JWTSecret),
		ttl:    cfg.JWTTTL,
	}
}

func (m *JWTManager) Generate(user User) (string, error) {
	now := time.Now()

	claims := JWTClaims{
		UserID: user.ID.String(),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(m.secret)
}

func (m *JWTManager) Verify(tokenString string) (TokenInfo, error) {
	var claims JWTClaims

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
			}

			return m.secret, nil
		},
		jwt.WithExpirationRequired(),
	)

	if err != nil || !token.Valid {
		return TokenInfo{}, ErrInvalidToken
	}

	return TokenInfo{
		UserID: uuid.MustParse(claims.UserID),
		Email:  claims.Email,
	}, nil
}
