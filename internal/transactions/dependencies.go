package transactions

import "context"

type AccountTransferPublisher interface {
	PublishTransfer(cmd UpdateAccountBalanceCommand) error
}

type UpdateAccountBalanceCommand struct {
	FromAccountID string
	ToAccountID   string
	Amount        float64
	ErrChan       chan error
}

type AccountReader interface {
	GetAccountByUserID(ctx context.Context, userID string) (*AccountInfo, error)
}

type AccountInfo struct {
	ID string
}
