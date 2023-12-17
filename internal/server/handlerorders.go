package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Schalure/gofermart/internal/gofermart/gofermaterrors"
)

//	Order registr handler. POST /api/user/orders
func (h *Handler) LoadOrder(w http.ResponseWriter, r* http.Request) {

	login := h.getLoginFromContext(r.Context())
	if login == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}	

	//	get url
	orderNumber, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Errorf("can`t read request body: %s", err.Error()).Error(), http.StatusBadRequest)
		return
	}

	if err = h.orderManager.LoadOrder(login, string(orderNumber)); err != nil {
		if errors.Is(err, gofermaterrors.InvalidOrderNumber) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, gofermaterrors.DublicateOrderNumberByUser) {
			w.WriteHeader(http.StatusOK)
			return
		}
		if errors.Is(err, gofermaterrors.DublicateOrderNumber) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

//	Get orders for user handler. GET /api/user/orders
func (h *Handler) GetOrders(w http.ResponseWriter, r* http.Request) {

	type orderInfo struct {
		Number string `json:"number"`
		Status string `json:"status"`
		Accrual int `json:"accrual"`
		UploadedAt time.Time `json:"uploaded_at"`
	}

	login := h.getLoginFromContext(r.Context())
	if login == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	orders, err := h.orderManager.GetOrders(login)
	if err != nil {
		if errors.Is(err, gofermaterrors.NoData){
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}

	ordersInfo := make([]orderInfo, len(orders))
	for i, order := range orders {
		ordersInfo[i] = orderInfo{
			Number: order.OrderNumber,
			Status: string(order.OrderStatus),
			Accrual: order.BonusPoints,
			UploadedAt: order.UploadedAt,
		}
	}

	data, err := json.Marshal(ordersInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}