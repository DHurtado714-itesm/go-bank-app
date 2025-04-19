package accounts

import (
	"context"
	"database/sql"
	"errors"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, acc *Account) error
	GetBalance(ctx context.Context, accountID string) (float64, error)
	GetAccountByUserID(ctx context.Context, userID string) (*Account, error)
}

type accountRepository struct {
	db *sql.DB
}

// CreateAccount implements AccountRepository.
func (r *accountRepository) CreateAccount(ctx context.Context, acc *Account) error {
	query := `
	INSERT INTO accounts (id, user_id, balance, currency, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
`
	_, err := r.db.ExecContext(ctx, query, acc.ID, acc.UserID, acc.Balance, acc.Currency, acc.CreatedAt, acc.UpdatedAt)
	return err
}

// GetAccountByUserID implements AccountRepository.
func (r *accountRepository) GetAccountByUserID(ctx context.Context, userID string) (*Account, error) {
	query := `
	SELECT id, user_id, balance, currency, created_at
	FROM accounts
	WHERE user_id = $1
	LIMIT 1
`
	row := r.db.QueryRowContext(ctx, query, userID)

	var acc Account
	err := row.Scan(&acc.ID, &acc.UserID, &acc.Balance, &acc.Currency, &acc.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &acc, nil
}

// GetBalance implements AccountRepository.
func (r *accountRepository) GetBalance(ctx context.Context, accountID string) (float64, error) {
	query := `SELECT balance FROM accounts WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, accountID)

	var balance float64
	err := row.Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, sql.ErrConnDone
	}

	return balance, nil
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{db: db}
}
