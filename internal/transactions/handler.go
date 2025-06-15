package transactions

import (
	"encoding/json"
	"go-bank-app/pkg/middleware"
	"net/http"
)

type TransactionHandler struct {
	service TransactionService
}

func NewTransactionHandler(service TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

type transferRequest struct {
	ToAccountID string  `json:"to_account_id"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"` // Optional, you can validate if needed
}

func (h *TransactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserIDKey).(string)

	var req transferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 || req.ToAccountID == "" {
		http.Error(w, "Invalid transfer data", http.StatusBadRequest)
		return
	}

	tx, err := h.service.Transfer(r.Context(), userID, req.ToAccountID, req.Amount, req.Currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tx)
}

func (h *TransactionHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserIDKey).(string)

	transactions, err := h.service.GetByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "Error retrieving history", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}
