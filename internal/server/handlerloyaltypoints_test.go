package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Schalure/gofermart/internal/loggers"
	"github.com/Schalure/gofermart/internal/mocks"
	"github.com/Schalure/gofermart/internal/storage"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetBalance(t *testing.T) {

	testMethod := "GET"
	testURL := "/api/user/balance"

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
		name                       string
		tokenCheckerLogin          string
		tokenCheckerError          error
		userManagerPoints          float64
		userManagerWithdrawnPoints float64

		want struct {
			data       string
			statusCode int
		}
	}{
		{
			name:                       "simple test",
			tokenCheckerLogin:          "Petya",
			tokenCheckerError:          nil,
			userManagerPoints:          500.5,
			userManagerWithdrawnPoints: 42,
			want: struct {
				data       string
				statusCode int
			}{
				data: `
				{
					"current": 500.5, 
					"withdrawn": 42
				}`,
				statusCode: http.StatusOK,
			},
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			tokenCheker.EXPECT().CheckValidJWT("qqqq").Return(test.tokenCheckerLogin, test.tokenCheckerError)
			userManager.EXPECT().GetUserInfo(gomock.Any(), test.tokenCheckerLogin).Return(storage.User{
				LoyaltyPoints:   test.userManagerPoints,
				WithdrawnPoints: 42,
			}, nil)

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

			assert.JSONEq(t, buf.String(), test.want.data)
		})
	}
}

func Test_WithdrawLoyaltyPoints(t *testing.T) {

	testMethod := "POST"
	testURL := "/api/user/balance/withdraw"

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
		name                      string
		requestBody               string
		tokenCheckerLogin         string
		tokenCheckerError         error
		loyaltySystemManagerLogin string
		loyaltySystemManagerOrder string
		loyaltySystemManagerSum   float64

		want struct {
			statusCode int
		}
	}{
		{
			name: "simple test",
			requestBody: `
			{
				"order": "2377225624",
				"sum": 751
			}`,
			tokenCheckerLogin:         "Petya",
			tokenCheckerError:         nil,
			loyaltySystemManagerLogin: "Petya",
			loyaltySystemManagerOrder: "2377225624",
			loyaltySystemManagerSum:   751,
			want: struct{ statusCode int }{
				statusCode: http.StatusOK,
			},
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			tokenCheker.EXPECT().CheckValidJWT("qqqq").Return(test.tokenCheckerLogin, test.tokenCheckerError)
			loyaltySystemManager.EXPECT().Withdraw(gomock.Any(), test.loyaltySystemManagerLogin, test.loyaltySystemManagerOrder, test.loyaltySystemManagerSum).Return(nil)

			req, err := http.NewRequest(testMethod, testServer.URL+testURL, bytes.NewReader([]byte(test.requestBody)))
			req.AddCookie(&http.Cookie{Name: authorizationCookie, Value: "qqqq"})
			require.NoError(t, err)

			resp, err := testServer.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, test.want.statusCode, resp.StatusCode)

		})
	}
}

func Test_GetOrdersWithdrawals(t *testing.T) {

	testMethod := "GET"
	testURL := "/api/user/withdrawals"

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
		name                       string
		tokenCheckerLogin          string
		tokenCheckerError          error
		loyaltySystemManagerLogin  string
		loyaltySystemManagerOrders []storage.Order
		want                       struct {
			statusCode int
			data       string
		}
	}{
		{
			name:                      "simple test",
			tokenCheckerLogin:         "Petya",
			tokenCheckerError:         nil,
			loyaltySystemManagerLogin: "Petya",
			loyaltySystemManagerOrders: []storage.Order{
				{
					OrderNumber: "2377225624",
					BonusPoints: 500,
					UploadedBonus: pgtype.Timestamptz{
						Time: time.Date(2020, 12, 9, 16, 9, 57, 0, time.FixedZone("", 60*60*3)),
					}, //time.Date(2020, 12, 9, 16, 9, 57, 0, time.FixedZone("", 60*60*3)),
				},
			},
			want: struct {
				statusCode int
				data       string
			}{
				statusCode: http.StatusOK,
				data: `  
				[
					{
						"order": "2377225624",
						"sum": 500,
						"processed_at": "2020-12-09T16:09:57+03:00"
					}
				]`,
			},
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			tokenCheker.EXPECT().CheckValidJWT("qqqq").Return(test.tokenCheckerLogin, test.tokenCheckerError)
			loyaltySystemManager.EXPECT().GetWithdraws(gomock.Any(), test.loyaltySystemManagerLogin).Return(test.loyaltySystemManagerOrders, nil)

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
