package transactions

import "time"

// TransactionFilter defines optional fields for querying transactions.
type TransactionFilter struct {
	Category  string
	StartDate *time.Time
	EndDate   *time.Time
}
