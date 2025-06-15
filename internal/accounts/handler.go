package accounts

import (
	"encoding/json"
	"go-bank-app/pkg/middleware"
	"net/http"
)

type AccountHandler struct {
	service AccountService
}

func NewAccountHandler(service AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

type createAccountRequest struct {
	Currency Currency `json:"currency"`
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserIDKey).(string)

	var req createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	account, err := h.service.Create(r.Context(), userID, req.Currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func (h *AccountHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserIDKey).(string)

	account, err := h.service.GetAccountByUserID(r.Context(), userID)
	if err != nil || account == nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"account_id": account.ID,
		"balance":    account.Balance,
		"currency":   account.Currency,
	})
}

// GetPublic exposes account balance without authentication, emulating an open banking endpoint.
func (h *AccountHandler) GetPublic(w http.ResponseWriter, r *http.Request) {
	accountID := r.URL.Query().Get("id")
	if accountID == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	account, err := h.service.GetAccountByID(r.Context(), accountID)
	if err != nil || account == nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"account_id": account.ID,
		"balance":    account.Balance,
		"currency":   account.Currency,
	})
}
