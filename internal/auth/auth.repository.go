package auth

import (
	"context"
	"database/sql"
	"errors"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users(id, email, hashed_password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.HashedPassword, user.CreatedAt, user.UpdatedAt)

	return err
}

func (r *authRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, hashed_password, created_at, updated_at
		FROM users
		where email = $1
	`

	row := r.db.QueryRowContext(ctx, query, email)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.HashedPassword, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}
