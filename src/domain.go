package auth

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID
	Email        string
	FirstName    string
	LastName     string
	PhoneNumber  string
	PasswordHash string
}

type AuthResult struct {
	UserID uuid.UUID
	Token  string
}

type TokenInfo struct {
	UserID uuid.UUID
	Email  string
}
