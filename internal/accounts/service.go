package accounts

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type AccountService interface {
	Create(ctx context.Context, userID string, currency Currency) (*Account, error)
	GetBalance(ctx context.Context, accountID string) (float64, error)
	GetAccountByUserID(ctx context.Context, userID string) (*Account, error)
}

type accountService struct {
	repo AccountRepository
}

// Create implements AccountService.
func (s *accountService) Create(ctx context.Context, userID string, currency Currency) (*Account, error) {
	if currency != CurrencyMXN {
		return nil, errors.New("invalid currency")
	}

	// Verificar si ya tiene cuenta
	existing, _ := s.repo.GetAccountByUserID(ctx, userID)
	if existing != nil {
		return nil, errors.New("user already has an account")
	}

	account := &Account{
		ID:        uuid.New().String(),
		UserID:    userID,
		Balance:   0,
		Currency:  currency,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.repo.CreateAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// GetAccountByUserID implements AccountService.
func (s *accountService) GetAccountByUserID(ctx context.Context, userID string) (*Account, error) {
	return s.repo.GetAccountByUserID(ctx, userID)
}

// GetBalance implements AccountService.
func (s *accountService) GetBalance(ctx context.Context, accountID string) (float64, error) {
	return s.repo.GetBalance(ctx, accountID)
}

func NewAccountService(repo AccountRepository) AccountService {
	return &accountService{repo: repo}
}
