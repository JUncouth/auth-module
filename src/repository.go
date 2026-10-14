package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user User) error {
	_, err := r.pool.Exec(
		ctx,
		`INSERT INTO users (id, email, first_name, last_name, phone_number, password_hash) VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID, user.Email, user.FirstName, user.LastName, user.PhoneNumber, user.PasswordHash,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrEmailAlreadyExists
		}
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User

	err := r.pool.QueryRow(
		ctx,
		`SELECT id, email, first_name, last_name, phone_number, password_hash FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.PhoneNumber, &user.PasswordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}

	if err != nil {
		return User{}, fmt.Errorf("find user by email: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) UpdatePassword(ctx context.Context, email string, passwordHash string) error {
	commandTag, err := r.pool.Exec(
		ctx,
		`UPDATE users SET password_hash = $1 WHERE email = $2`,
		passwordHash, email,
	)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

