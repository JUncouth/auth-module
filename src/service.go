package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	users  UserRepository
	hasher PasswordHasher
	tokens TokenManager
}

func NewAuthService(users UserRepository, hasher PasswordHasher, tokens TokenManager) AuthService {
	return &Service{
		users:  users,
		hasher: hasher,
		tokens: tokens,
	}
}

func (s *Service) Register(
	ctx context.Context,
	email string,
	firstName string,
	lastName string,
	phoneNumber string,
	password string,
) (AuthResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || firstName == "" || lastName == "" || password == "" {
		return AuthResult{}, ErrInvalidInput
	}

	user := User{
		ID:          uuid.New(),
		Email:       email,
		FirstName:   firstName,
		LastName:    lastName,
		PhoneNumber: phoneNumber,
	}

	hash, err := s.hasher.Hash(password)
	if err != nil {
		return AuthResult{}, err
	}

	user.PasswordHash = hash

	if err := s.users.Create(ctx, user); err != nil {
		return AuthResult{}, err
	}

	token, err := s.tokens.Generate(user)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{
		UserID: user.ID,
		Token:  token,
	}, nil
}

func (s *Service) Authenticate(
	ctx context.Context,
	email string,
	password string,
) (AuthResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		if err == ErrUserNotFound {
			return AuthResult{}, ErrInvalidCredentials
		}
		return AuthResult{}, err
	}

	if err := s.hasher.Compare(user.PasswordHash, password); err != nil {
		return AuthResult{}, ErrInvalidCredentials
	}

	token, err := s.tokens.Generate(user)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{
		UserID: user.ID,
		Token:  token,
	}, nil
}

func (s *Service) ChangePassword(
	ctx context.Context,
	email string,
	oldPassword string,
	newPassword string,
) error {
	email = strings.ToLower(strings.TrimSpace(email))
	if newPassword == "" {
		return ErrInvalidInput
	}

	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	if err := s.hasher.Compare(user.PasswordHash, oldPassword); err != nil {
		return ErrInvalidCredentials
	}

	hash, err := s.hasher.Hash(newPassword)
	if err != nil {
		return err
	}

	return s.users.UpdatePassword(ctx, email, hash)
}

func (s *Service) ResetPassword(ctx context.Context, email string) error {
	email = strings.ToLower(strings.TrimSpace(email))

	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	password, err := generateTemporaryPassword()
	if err != nil {
		return err
	}

	hash, err := s.hasher.Hash(password)
	if err != nil {
		return err
	}

	if err := s.users.UpdatePassword(ctx, email, hash); err != nil {
		return err
	}
	println("RESET PASSWORD:", user.Email, password)

	return nil
}

func (s *Service) VerifyToken(ctx context.Context, token string) (TokenInfo, error) {
	return s.tokens.Verify(token)
}

func generateTemporaryPassword() (string, error) {
	data := make([]byte, 12)

	if _, err := rand.Read(data); err != nil {
		return "", err
	}

	return hex.EncodeToString(data), nil
}
