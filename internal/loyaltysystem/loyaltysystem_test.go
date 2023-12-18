package loyaltysystem

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			status OrderStatus
			accrual float64
			statusCode int
		}
	}{
		{
			name: "simple test",
			requestOrder: "1234567897",
			responseBody: `{"order":"1234567897","status":"PROCESSED","accrual":500}`,
			responseCode: http.StatusOK,
			responseDelay: time.Second * 0,
			want: struct{order string; status OrderStatus; accrual float64; statusCode int}{
				order: "1234567897",
				status: Processed,
				accrual: 500,
				statusCode: http.StatusOK,
			},
		},
		{
			name: "timeout test",
			requestOrder: "1234567897",
			responseBody: `{"order":"1234567897","status":"PROCESSED","accrual":500}`,
			responseCode: http.StatusOK,
			responseDelay: time.Second * 2,
			want: struct{order string; status OrderStatus; accrual float64; statusCode int}{
				statusCode: 0,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(test.responseDelay)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(test.responseCode)
				w.Write([]byte(test.responseBody))
			}))
			defer testServer.Close()
			loyaltySystem := NewLoyaltySystem(testServer.URL)


			ctx, cancel := context.WithTimeout(context.Background(), time.Second * 1)
			order, statusCode := loyaltySystem.OrderCheck(ctx, test.requestOrder)
			cancel()

			require.Equal(t, test.want.statusCode, statusCode)

			assert.Equal(t, test.want.order, order.OrderNumber)
			assert.Equal(t, string(test.want.status), string(order.OrderStatus))
			assert.Equal(t, test.want.accrual, order.BonusPoints)

		})
	}
}