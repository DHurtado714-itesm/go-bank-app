package main

import (
	"context"
	"database/sql"
	"fmt"
	config "go-bank-app/configs"
	"go-bank-app/internal/accounts"
	"go-bank-app/internal/auth"
	"go-bank-app/internal/transactions"
	"go-bank-app/pkg/middleware"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env file")
	}

	dbURL, err := config.GetString("POSTGRES_DB_URI")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}
	defer conn.Close()

	// ─── AUTH ─────────────────────────────────────────────
	authRepo := auth.NewAuthRepository(conn)
	authService := auth.NewAuthService(authRepo)
	authHandler := auth.NewAuthHandler(authService)

	// Rutas públicas
	http.HandleFunc("/auth/register", authHandler.Register)
	http.HandleFunc("/auth/login", authHandler.Login)

	// Ruta protegida con middleware
	http.Handle("/auth/me", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.ContextUserIDKey).(string)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"user_id": "%s"}`, userID)
	})))

	// ─── ACCOUNTS ─────────────────────────────────────────
	accountRepo := accounts.NewAccountRepository(conn)
	accountService := accounts.NewAccountService(accountRepo)
	accountHandler := accounts.NewAccountHandler(accountService)
	
	// Init listener to process commands
	accounts.StartAccountBalanceWorker(accountRepo)

	http.Handle("/accounts", middleware.AuthMiddleware(http.HandlerFunc(accountHandler.Create)))
	http.Handle("/accounts/balance", middleware.AuthMiddleware(http.HandlerFunc(accountHandler.GetBalance)))

	// ─── TRANSACTIONS ─────────────────────────────────────
	txPublisher := &AccountTransferChannelAdapter{}
	accountReader := &AccountReaderAdapter{accountService: accountService}

	txRepo := transactions.NewTransactionRepository(conn)
	txService := transactions.NewTransactionService(txRepo, txPublisher, accountReader)
	txHandler := transactions.NewTransactionHandler(txService)
	
	http.Handle("/transactions/transfer", middleware.AuthMiddleware(http.HandlerFunc(txHandler.Transfer)))
	http.Handle("/transactions/history", middleware.AuthMiddleware(http.HandlerFunc(txHandler.GetHistory)))

	// ─── SERVER ───────────────────────────────────────────
	port := ":8070"
	fmt.Println("🚀 Server running at http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// ─── ADAPTERS ───────────────────────────────────────────────

type AccountTransferChannelAdapter struct{}

func (a *AccountTransferChannelAdapter) PublishTransfer(cmd transactions.UpdateAccountBalanceCommand) error {
	internalCmd := accounts.UpdateAccountBalanceCommand{
		FromAccountID: cmd.FromAccountID,
		ToAccountID:   cmd.ToAccountID,
		Amount:        cmd.Amount,
		ErrChan:       cmd.ErrChan,
	}

	accounts.AccountUpdateChannel <- internalCmd
	return nil
}

type AccountReaderAdapter struct {
	accountService accounts.AccountService
}

func (a *AccountReaderAdapter) GetAccountByUserID(ctx context.Context, userID string) (*transactions.AccountInfo, error) {
	account, err := a.accountService.GetAccountByUserID(ctx, userID)
	if err != nil || account == nil {
		return nil, err
	}
	return &transactions.AccountInfo{ID: account.ID}, nil
}