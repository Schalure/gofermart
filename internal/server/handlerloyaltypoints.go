package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Schalure/gofermart/internal/gofermart/gofermaterrors"
)

// Get loyality points balance. GET /api/user/balance
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {

	type Balance struct {
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}

	login := h.getLoginFromContext(r.Context())
	if login == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := h.userManager.GetUserInfo(r.Context(), login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(Balance{
		Current:   user.LoyaltyPoints,
		Withdrawn: user.WithdrawnPoints,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *Handler) WithdrawLoyaltyPoints(w http.ResponseWriter, r *http.Request) {

	type Order struct {
		OrderNumber string  `json:"order"`
		Sum         float64 `json:"sum"`
	}

	login := h.getLoginFromContext(r.Context())
	if login == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var order Order
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.loyaltySystemManager.Withdraw(r.Context(), login, order.OrderNumber, order.Sum); err != nil {
		if errors.Is(err, gofermaterrors.ErrInvalidOrderNumber) {
			http.Error(w, gofermaterrors.ErrInvalidOrderNumber.Error(), http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, gofermaterrors.ErrInsufficientFunds) {
			http.Error(w, gofermaterrors.ErrInsufficientFunds.Error(), http.StatusPaymentRequired)
			return
		}
		http.Error(w, gofermaterrors.ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Receiving information on withdrawal of funds. GET /api/user/withdrawals
func (h *Handler) GetOrdersWithdrawals(w http.ResponseWriter, r *http.Request) {

	type Withdrawal struct {
		OrderNumber string  `json:"order"`
		Sum         float64 `json:"sum"`
		ProcessedAt string  `json:"processed_at"`
	}

	login := h.getLoginFromContext(r.Context())
	if login == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	orders, err := h.loyaltySystemManager.GetWithdraws(r.Context(), login)
	if err != nil {
		if errors.Is(err, gofermaterrors.ErrNoData) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, gofermaterrors.ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	withdrawals := make([]Withdrawal, len(orders))
	for i, order := range orders {
		withdrawals[i] = Withdrawal{
			OrderNumber: order.OrderNumber,
			Sum:         order.BonusWithdraw,
			ProcessedAt: order.UploadedBonus.Time.Format(time.RFC3339),
		}
	}

	data, err := json.Marshal(withdrawals)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}
