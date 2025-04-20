package transactions

import (
	"context"
	"errors"
	"go-bank-app/pkg/csvwriter"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type TransactionService interface {
	Transfer(ctx context.Context, fromID, toID string, amount float64, currency string) (*Transaction, error)
	GetByAccount(ctx context.Context, accountID string) ([]Transaction, error)
	GenerateStatementCSV(transactions []Transaction, filePath string) error
}

type transactionService struct {
	repo      TransactionRepository
	publisher AccountTransferPublisher
	reader    AccountReader
}

// GetByAccount implements TransactionService.
func (s *transactionService) GetByAccount(ctx context.Context, accountID string) ([]Transaction, error) {
	return s.repo.GetByAccount(ctx, accountID)

}

// Transfer implements TransactionService.
func (s *transactionService) Transfer(ctx context.Context, fromID string, toID string, amount float64, currency string) (*Transaction, error) {
	if fromID == toID {
		return nil, errors.New("cannot transfer to the same account")
	}

	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	account, err := s.reader.GetAccountByUserID(ctx, fromID)
	if err != nil || account == nil {
		return nil, errors.New("origin account not found")
	}

	errChan := make(chan error)
	cmd := UpdateAccountBalanceCommand{
		FromAccountID: account.ID,
		ToAccountID:   toID,
		Amount:        amount,
		ErrChan:       errChan,
	}

	if err := s.publisher.PublishTransfer(cmd); err != nil {
		return nil, err
	}

	if err := <-errChan; err != nil {
		return nil, err
	}

	tx := &Transaction{
		ID:            uuid.New().String(),
		FromAccountID: account.ID,
		ToAccountID:   toID,
		Amount:        amount,
		Currency:      currency,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.Create(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *transactionService) GenerateStatementCSV(transactions []Transaction, filePath string) error {
	w := csvwriter.NewCSVWriter([]string{
		"Transaction ID",
		"From Account",
		"To Account",
		"Amount",
		"Currency",
		"Date",
	})

	for _, t := range transactions {
		row := []string{
			t.ID,
			t.FromAccountID,
			t.ToAccountID,
			strconv.FormatFloat(t.Amount, 'f', 2, 64),
			t.Currency,
			t.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if err := w.AddRow(row); err != nil {
			return err
		}
	}

	return w.WriteToFile(filePath)
}

func NewTransactionService(repo TransactionRepository, publisher AccountTransferPublisher, reader AccountReader) TransactionService {
	return &transactionService{repo: repo, publisher: publisher, reader: reader}
}
