package auth

import (
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	UpdatePassword(ctx context.Context, email string, passwordHash string) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash string, password string) error
}

type TokenManager interface {
	Generate(user User) (string, error)
	Verify(token string) (TokenInfo, error)
}

type AuthService interface {
	Register(
		ctx context.Context,
		email string,
		firstName string,
		lastName string,
		phoneNumber string,
		password string,
	) (AuthResult, error)

	Authenticate(
		ctx context.Context,
		email string,
		password string,
	) (AuthResult, error)

	ChangePassword(
		ctx context.Context,
		email string,
		oldPassword string,
		newPassword string,
	) error

	ResetPassword(
		ctx context.Context,
		email string,
	) error

	VerifyToken(
		ctx context.Context,
		token string,
	) (TokenInfo, error)
}
