package loyaltysystem

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Schalure/gofermart/internal/storage"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func Test_OrderCheck(t *testing.T) {

	testCases := []struct{
		name string
		requestOrder string
		responseBody string
		responseCode int
		responseDelay time.Duration
		want struct{
			order string
			status string
			accrual int
		}
	}{
		{
			name: "simple test",
			requestOrder: "1234567897",
			responseBody: `{"order":"1234567897","status":"PROCESSED","accrual":500}`,
			responseCode: http.StatusOK,
			responseDelay: time.Second * 0,
			want: struct{order string; status string; accrual int}{
				order: "1234567897",
				status: "PROCESSED",
				accrual: 500,
			},
		},
	}

	r := chi.NewMux()
	testServer := httptest.NewServer(r)
	defer testServer.Close()

	loyaltySystem := NewLoyaltySystem(testServer.URL)

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			r.Get("/api/orders/{number}", func (w http.ResponseWriter, r *http.Request) {
				time.Sleep(test.responseDelay)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(test.responseCode)
				w.Write([]byte(test.responseBody))
			})

			resultCh := make(chan storage.Order)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
			defer cancel()
			loyaltySystem.OrderCheck(ctx, storage.Order{OrderNumber: test.requestOrder}, resultCh)
			close(resultCh)

			order := <-resultCh

			assert.Equal(t, test.want.order, order.OrderNumber)
			assert.Equal(t, test.want.status, order.OrderStatus)
			assert.Equal(t, test.want.accrual, order.BonusPoints)
		})
	}
}