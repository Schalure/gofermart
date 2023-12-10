package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Schalure/gofermart/internal/configs"
	"github.com/Schalure/gofermart/internal/gofermart"
	"github.com/Schalure/gofermart/internal/loggers"
	"github.com/Schalure/gofermart/internal/storage/mockstor"
	"github.com/stretchr/testify/assert"
)



func Test_UserRegistration(t *testing.T) {

	config, _ := configs.NewConfig()
	logger := loggers.NewLogger(config)
	stor := mockstor.NewStorage()
	service := gofermart.NewGofermart(config, stor, logger)
	server := NewServer(config, NewHandler(service), NewMidleware())

	testServer := httptest.NewServer(server.router)

	testCases := []struct{
		name string
		requestBody string
		want struct{
			statusCode int
			token string
		}
	}{
		{
			name: "simple test",
			requestBody: `{"login":"Mihail","password": "q1w2e3r4"}`,
			want: struct{statusCode int; token string}{
				statusCode: http.StatusOK,
			},
		},
	}

	for _, test := range testCases{
		t.Run(test.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/register", bytes.NewBufferString(test.requestBody))

			server.router.ServeHTTP(w, req)

			assert.Equal(t, test.want.statusCode, req.Response.StatusCode)
		})
	}
	
}