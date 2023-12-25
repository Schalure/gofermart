package server

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Schalure/gofermart/internal/gofermart/gofermaterrors"
	"github.com/Schalure/gofermart/internal/loggers"
	"github.com/Schalure/gofermart/internal/mocks"
	"github.com/Schalure/gofermart/internal/storage"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LoadOrder(t *testing.T) {

	testMethod := "POST"
	testURL := "/api/user/orders"

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userManager := mocks.NewMockUserManager(mockController)
	orderManager := mocks.NewMockOrderManager(mockController)
	loyaltySystemManager := mocks.NewMockLoyaltySystemManager(mockController)
	logger := loggers.NewLogger("debug")
	tokenCheker := mocks.NewMockTokenCheker(mockController)
	server := NewServer(NewHandler(userManager, orderManager, loyaltySystemManager), NewMidleware(logger, tokenCheker))

	testServer := httptest.NewServer(server.Router)
	defer testServer.Close()

	testCases := []struct {
		name                   string
		requestOrder           string
		tokenChekerGomockCall  *gomock.Call
		orderManagerGomockCall *gomock.Call
		want                   struct {
			statusCode int
		}
	}{
		{
			name:                   "simple test",
			requestOrder:           "1234567897",
			orderManagerGomockCall: orderManager.EXPECT().LoadOrder(gomock.Any(), "Petya", "1234567897").Return(nil),
			tokenChekerGomockCall:  tokenCheker.EXPECT().CheckValidJWT("qqqq").Return("Petya", nil),
			want: struct{ statusCode int }{
				statusCode: http.StatusAccepted,
			},
		},
		{
			name:                   "invalid order test",
			requestOrder:           "123456789",
			orderManagerGomockCall: orderManager.EXPECT().LoadOrder(gomock.Any(), "Petya", "123456789").Return(gofermaterrors.ErrInvalidOrderNumber),
			tokenChekerGomockCall:  tokenCheker.EXPECT().CheckValidJWT("qqqq").Return("Petya", nil),
			want: struct{ statusCode int }{
				statusCode: http.StatusUnprocessableEntity,
			},
		},
		{
			name:                   "dublicate order number by user test",
			requestOrder:           "1234567897",
			orderManagerGomockCall: orderManager.EXPECT().LoadOrder(gomock.Any(), "Petya", "1234567897").Return(gofermaterrors.ErrDublicateOrderNumberByUser),
			tokenChekerGomockCall:  tokenCheker.EXPECT().CheckValidJWT("qqqq").Return("Petya", nil),
			want: struct{ statusCode int }{
				statusCode: http.StatusOK,
			},
		},
		{
			name:                   "dublicate order number test",
			requestOrder:           "1234567897",
			orderManagerGomockCall: orderManager.EXPECT().LoadOrder(gomock.Any(), "Petya", "1234567897").Return(gofermaterrors.ErrDublicateOrderNumber),
			tokenChekerGomockCall:  tokenCheker.EXPECT().CheckValidJWT("qqqq").Return("Petya", nil),
			want: struct{ statusCode int }{
				statusCode: http.StatusConflict,
			},
		},
		{
			name:                  "unauthorized test",
			requestOrder:          "1234567897",
			tokenChekerGomockCall: tokenCheker.EXPECT().CheckValidJWT("qqqq").Return("", errors.New("")),
			want: struct{ statusCode int }{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name:                   "other errors test",
			requestOrder:           "1234567897",
			orderManagerGomockCall: orderManager.EXPECT().LoadOrder(gomock.Any(), "Petya", "1234567897").Return(gofermaterrors.ErrInternal),
			tokenChekerGomockCall:  tokenCheker.EXPECT().CheckValidJWT("qqqq").Return("Petya", nil),
			want: struct{ statusCode int }{
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			req, err := http.NewRequest(testMethod, testServer.URL+testURL, bytes.NewBufferString(test.requestOrder))
			req.AddCookie(&http.Cookie{Name: authorizationCookie, Value: "qqqq"})
			require.NoError(t, err)

			resp, err := testServer.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, test.want.statusCode, resp.StatusCode)
		})
	}
}

func Test_GetOrders(t *testing.T) {

	testMethod := "GET"
	testURL := "/api/user/orders"

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userManager := mocks.NewMockUserManager(mockController)
	orderManager := mocks.NewMockOrderManager(mockController)
	loyaltySystemManager := mocks.NewMockLoyaltySystemManager(mockController)
	logger := loggers.NewLogger("debug")
	tokenCheker := mocks.NewMockTokenCheker(mockController)
	server := NewServer(NewHandler(userManager, orderManager, loyaltySystemManager), NewMidleware(logger, tokenCheker))

	testServer := httptest.NewServer(server.Router)
	defer testServer.Close()

	testCases := []struct {
		name               string
		tokenCheckerLogin  string
		tokenCheckerError  error
		orderManagerOrders []storage.Order
		orderManagerError  error
		want               struct {
			data       string
			statusCode int
		}
	}{
		{
			name:              "simple test",
			tokenCheckerLogin: "Petya",
			tokenCheckerError: nil,
			orderManagerOrders: []storage.Order{
				{
					OrderNumber: "9278923470",
					OrderStatus: storage.OrderStatusProcessed,
					BonusPoints: 500,
					UploadedOrder: pgtype.Timestamptz{
						Time: time.Date(2020, 12, 10, 15, 15, 45, 0, time.FixedZone("", 60*60*3)),
					},
				},
				{
					OrderNumber: "12345678903",
					OrderStatus: storage.OrderStatusProcessing,
					UploadedOrder: pgtype.Timestamptz{
						Time: time.Date(2020, 12, 10, 15, 12, 1, 0, time.FixedZone("", 60*60*3)),
					},
				},
				{
					OrderNumber: "346436439",
					OrderStatus: storage.OrderStatusInvalid,
					UploadedOrder: pgtype.Timestamptz{
						Time: time.Date(2020, 12, 9, 16, 9, 53, 0, time.FixedZone("", 60*60*3)),
					},
				},
			},
			orderManagerError: nil,
			want: struct {
				data       string
				statusCode int
			}{
				data: `[
					{
						"number": "9278923470",
						"status": "PROCESSED",
						"accrual": 500,
						"uploaded_at": "2020-12-10T15:15:45+03:00"
					},
					{
						"number": "12345678903",
						"status": "PROCESSING",
						"uploaded_at": "2020-12-10T15:12:01+03:00"
					},
					{
						"number": "346436439",
						"status": "INVALID",
						"uploaded_at": "2020-12-09T16:09:53+03:00"
					}]`,
				statusCode: http.StatusOK,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			tokenCheker.EXPECT().CheckValidJWT("qqqq").Return(test.tokenCheckerLogin, test.tokenCheckerError)
			orderManager.EXPECT().GetOrders(gomock.Any(), test.tokenCheckerLogin).Return(test.orderManagerOrders, test.orderManagerError)

			req, err := http.NewRequest(testMethod, testServer.URL+testURL, nil)
			req.AddCookie(&http.Cookie{Name: authorizationCookie, Value: "qqqq"})
			require.NoError(t, err)

			resp, err := testServer.Client().Do(req)
			require.NoError(t, err)

			assert.Equal(t, test.want.statusCode, resp.StatusCode)

			var buf bytes.Buffer
			_, err = buf.ReadFrom(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.JSONEq(t, test.want.data, buf.String())

		})
	}
}
