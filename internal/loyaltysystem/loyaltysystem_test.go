package loyaltysystem

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Schalure/gofermart/internal/storage"
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
			accrual int
			err error
		}
	}{
		{
			name: "simple test",
			requestOrder: "1234567897",
			responseBody: `{"order":"1234567897","status":"PROCESSED","accrual":500}`,
			responseCode: http.StatusOK,
			responseDelay: time.Second * 0,
			want: struct{order string; status OrderStatus; accrual int; err error}{
				order: "1234567897",
				status: Processed,
				accrual: 500,
				err: nil,
			},
		},
		{
			name: "timeout test",
			requestOrder: "1234567897",
			responseBody: `{"order":"1234567897","status":"PROCESSED","accrual":500}`,
			responseCode: http.StatusOK,
			responseDelay: time.Second * 2,
			want: struct{order string; status OrderStatus; accrual int; err error}{
				order: "1234567897",
				status: Processed,
				accrual: 500,
				err: fmt.Errorf("the response timeout time has been exceeded"),
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

			resultCh := make(chan struct{order storage.Order; statusCode int})
			var err error
			ctx, cancel := context.WithTimeout(context.Background(), time.Second * 1)

			go loyaltySystem.OrderCheck(ctx, storage.Order{OrderNumber: test.requestOrder}, resultCh)

			select {
			case <-ctx.Done():
				cancel()
				err = fmt.Errorf("the response timeout time has been exceeded")
			case order := <-resultCh:
				assert.Equal(t, test.want.order, order.order.OrderNumber)
				assert.Equal(t, string(test.want.status), string(order.order.OrderStatus))
				assert.Equal(t, test.want.accrual, order.order.BonusPoints)
				assert.Equal(t, test.responseCode, order.statusCode)
			}
			close(resultCh)
			cancel()

			require.Equal(t, test.want.err, err)
		})
	}
}