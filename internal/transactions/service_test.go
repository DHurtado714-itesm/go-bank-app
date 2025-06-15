package transactions

import (
	"context"
	"reflect"
	"testing"
)

type mockRepo struct {
	gotAccountID string
	transactions []Transaction
}

func (m *mockRepo) Create(ctx context.Context, tx *Transaction) error { return nil }
func (m *mockRepo) GetByAccount(ctx context.Context, accountID string) ([]Transaction, error) {
	m.gotAccountID = accountID
	return m.transactions, nil
}

type mockReader struct {
	acc *AccountInfo
	err error
}

func (m *mockReader) GetAccountByUserID(ctx context.Context, userID string) (*AccountInfo, error) {
	return m.acc, m.err
}

func TestTransactionService_GetByUser(t *testing.T) {
	repo := &mockRepo{transactions: []Transaction{{ID: "tx1"}}}
	reader := &mockReader{acc: &AccountInfo{ID: "acc123"}}
	svc := NewTransactionService(repo, nil, reader)

	txs, err := svc.GetByUser(context.Background(), "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.gotAccountID != "acc123" {
		t.Errorf("expected repo to be called with accountID acc123, got %s", repo.gotAccountID)
	}
	if !reflect.DeepEqual(txs, repo.transactions) {
		t.Errorf("expected %v, got %v", repo.transactions, txs)
	}
}

func TestTransactionService_GetByUser_NoAccount(t *testing.T) {
	repo := &mockRepo{}
	reader := &mockReader{acc: nil}
	svc := NewTransactionService(repo, nil, reader)

	_, err := svc.GetByUser(context.Background(), "user1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
