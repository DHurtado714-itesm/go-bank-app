package accounts

import (
	"context"
	"log"
)

type UpdateAccountBalanceCommand struct {
	FromAccountID string
	ToAccountID   string
	Amount        float64
	ErrChan       chan error
}

var AccountUpdateChannel = make(chan UpdateAccountBalanceCommand)

func StartAccountBalanceWorker(repo AccountRepository) {
	go func() {
		for cmd := range AccountUpdateChannel {
			ctx := context.Background()

			err := repo.Transfer(ctx, cmd.FromAccountID, cmd.ToAccountID, cmd.Amount)
			if err != nil {
				log.Printf("❌ Transfer error: %v", err)
			}

			cmd.ErrChan <- err
		}
	}()
}
