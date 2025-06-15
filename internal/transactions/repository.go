package transactions

import (
	"context"
	"database/sql"
	"log"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *Transaction) error
	GetByAccount(ctx context.Context, accountID string, filter TransactionFilter) ([]Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, t *Transaction) error {
	log.Printf("💾 Saving transaction: FROM %s TO %s AMOUNT %.2f CURRENCY %s",
		t.FromAccountID, t.ToAccountID, t.Amount, t.Currency)

	query := `
                INSERT INTO transactions (id, from_account_id, to_account_id, amount, currency, description, category, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        `
	_, err := r.db.ExecContext(ctx, query, t.ID, t.FromAccountID, t.ToAccountID, t.Amount, t.Currency, t.Description, t.Category, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		log.Printf("❌ Failed to save transaction: %v", err)
	}
	return err
}

func (r *transactionRepository) GetByAccount(ctx context.Context, accountID string, filter TransactionFilter) ([]Transaction, error) {
	baseQuery := `
                SELECT id, from_account_id, to_account_id, amount, currency, description, category, created_at, updated_at
                FROM transactions
                WHERE (from_account_id = $1 OR to_account_id = $1)`

	var args []interface{}
	args = append(args, accountID)
	if filter.Category != "" {
		baseQuery += " AND category = $2"
		args = append(args, filter.Category)
	}

	baseQuery += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction

	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.FromAccountID, &t.ToAccountID, &t.Amount, &t.Currency, &t.Description, &t.Category, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}
