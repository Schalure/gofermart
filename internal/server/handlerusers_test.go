package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Schalure/gofermart/internal/configs"
	"github.com/Schalure/gofermart/internal/gofermart/gofermaterrors"
	"github.com/Schalure/gofermart/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)



func Test_UserRegistration(t *testing.T) {

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	config, _ := configs.NewConfig()
	mockUserManager := mocks.NewMockUserManager(mockController)
	server := NewServer(config, NewHandler(mockUserManager), NewMidleware())

	testServer := httptest.NewServer(server.router)
	defer testServer.Close()

	testMethod := "POST"
	testURL := "/api/user/register"


	testCases := []struct{
		name string
		requestBody string
		gomockCall *gomock.Call
		want struct{
			statusCode int
			token string
		}
	}{
		{
			name: "simple test",
			requestBody: `{"login":"Mihail","password":"q1w2e3r4"}`,
			gomockCall: mockUserManager.EXPECT().CreateUser(gomock.Any(), "Mihail", "q1w2e3r4").Return("token", nil),
			want: struct{statusCode int; token string}{
				statusCode: http.StatusOK,
				token: "token",
			},
		},
		{
			name: "dublicate test",
			requestBody: `{"login":"Mihail","password":"q1w2e3r4"}`,
			gomockCall: mockUserManager.EXPECT().CreateUser(gomock.Any(), "Mihail", "q1w2e3r4").Return("", gofermaterrors.LoginAlreadyTaken),
			want: struct{statusCode int; token string}{
				statusCode: http.StatusConflict,
			},
		},
		{
			name: "bad body test",
			requestBody: `bad request body`,
			want: struct{statusCode int; token string}{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "empty body",
			requestBody: `{}`,
			gomockCall: mockUserManager.EXPECT().CreateUser(gomock.Any(), "", "").Return("", gofermaterrors.InvalidLogin),
			want: struct{statusCode int; token string}{
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	for _, test := range testCases{
		t.Run(test.name, func(t *testing.T) {

			req, err := http.NewRequest(testMethod, testServer.URL + testURL, bytes.NewBufferString(test.requestBody))
			require.NoError(t, err)

			resp, err := testServer.Client().Do(req)
			require.NoError(t, err)

			assert.Equal(t, test.want.statusCode, resp.StatusCode)

			
			assert.Equal(t, test.want.token, func () string {
				cookies := resp.Cookies()
				for _, cookie := range cookies {
					if cookie.Name == authorizationCookie {
						return cookie.Value
					}
				}
				return ""
			}())
		})
	}
}


func Test_UserAuthentication(t *testing.T) {

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	config, _ := configs.NewConfig()
	mockUserManager := mocks.NewMockUserManager(mockController)
	server := NewServer(config, NewHandler(mockUserManager), NewMidleware())

	testServer := httptest.NewServer(server.router)
	defer testServer.Close()

	testMethod := "POST"
	testURL := "/api/user/login"

	testCases := []struct {
		name string
		requestBody string
		gomockCall *gomock.Call
		want struct{
			statusCode int
			token string
		}
	}{
		{
			name: "simple test",
			requestBody: `{"login":"Mihail","password":"q1w2e3r4"}`,
			gomockCall: mockUserManager.EXPECT().AuthenticationUser(gomock.Any(), "Mihail", "q1w2e3r4").Return("token", nil),
			want: struct{statusCode int; token string}{
				statusCode: http.StatusOK,
				token: "token",
			},
		},
		{
			name: "invalid login or password",
			requestBody: `{"login":"Mihail","password":"q1w2e3r4"}`,
			gomockCall: mockUserManager.EXPECT().AuthenticationUser(gomock.Any(), "Mihail", "q1w2e3r4").Return("", gofermaterrors.InvalidLoginPassword),
			want: struct{statusCode int; token string}{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "bad body",
			requestBody: `bad body`,
			want: struct{statusCode int; token string}{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "empty body",
			requestBody: `{}`,
			gomockCall: mockUserManager.EXPECT().AuthenticationUser(gomock.Any(), "", "").Return("", gofermaterrors.InvalidLoginPassword),
			want: struct{statusCode int; token string}{
				statusCode: http.StatusUnauthorized,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			req, err := http.NewRequest(testMethod, testServer.URL + testURL, bytes.NewBufferString(test.requestBody))
			require.NoError(t, err)

			resp, err := testServer.Client().Do(req)
			require.NoError(t, err)

			assert.Equal(t, test.want.statusCode, resp.StatusCode)

			
			assert.Equal(t, test.want.token, func () string {
				cookies := resp.Cookies()
				for _, cookie := range cookies {
					if cookie.Name == authorizationCookie {
						return cookie.Value
					}
				}
				return ""
			}())
		})
	}
}
