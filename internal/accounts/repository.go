package accounts

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, acc *Account) error
	GetBalance(ctx context.Context, accountID string) (float64, error)
	GetAccountByUserID(ctx context.Context, userID string) (*Account, error)
	Transfer(ctx context.Context, fromID, toID string, amount float64) error
}

type accountRepository struct {
	db *sql.DB
}

// Transfer implements AccountRepository.
func (r *accountRepository) Transfer(ctx context.Context, fromID string, toID string, amount float64) error {
	log.Printf("💸 Starting transfer of %.2f from [%s] to [%s]", amount, fromID, toID)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("❌ Failed to begin transaction: %v", err)
		return err
	}

	// 1. Leer saldo del emisor
	var fromBalance float64
	err = tx.QueryRowContext(ctx, `SELECT balance FROM accounts WHERE id = $1`, fromID).Scan(&fromBalance)
	if err != nil {
		log.Printf("❌ Failed to get balance for fromAccount [%s]: %v", fromID, err)
		tx.Rollback()
		return err
	}
	log.Printf("💼 FromAccount balance: %.2f", fromBalance)

	// 2. Verificar fondos
	if fromBalance < amount {
		log.Printf("❌ Insufficient funds in [%s]: has %.2f, needs %.2f", fromID, fromBalance, amount)
		tx.Rollback()
		return errors.New("insufficient funds")
	}

	// 3. Descontar del emisor
	_, err = tx.ExecContext(ctx, `UPDATE accounts SET balance = balance - $1 WHERE id = $2`, amount, fromID)
	if err != nil {
		log.Printf("❌ Failed to debit fromAccount [%s]: %v", fromID, err)
		tx.Rollback()
		return err
	}
	log.Printf("✅ Debited %.2f from [%s]", amount, fromID)

	// 4. Agregar al receptor
	_, err = tx.ExecContext(ctx, `UPDATE accounts SET balance = balance + $1 WHERE id = $2`, amount, toID)
	if err != nil {
		log.Printf("❌ Failed to credit toAccount [%s]: %v", toID, err)
		tx.Rollback()
		return err
	}
	log.Printf("✅ Credited %.2f to [%s]", amount, toID)

	// 5. Commit
	err = tx.Commit()
	if err != nil {
		log.Printf("❌ Failed to commit transfer: %v", err)
		return err
	}

	log.Printf("✅ Transfer of %.2f from [%s] to [%s] completed successfully", amount, fromID, toID)
	return nil
}


// CreateAccount implements AccountRepository.
func (r *accountRepository) CreateAccount(ctx context.Context, acc *Account) error {
	query := `
	INSERT INTO accounts (id, user_id, balance, currency, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
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
