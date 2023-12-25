package server

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Schalure/gofermart/internal/gofermart"
)

//go:generate mockgen -destination=../mocks/mock_tokencheker.go -package=mocks github.com/Schalure/gofermart/internal/server TokenCheker
type TokenCheker interface {
	CheckValidJWT(tokenString string) (string, error)
}

type Middleware struct {
	logger      gofermart.Loggerer
	tokenCheker TokenCheker
}

func NewMidleware(logger gofermart.Loggerer, tokenCheker TokenCheker) *Middleware {

	return &Middleware{
		logger:      logger,
		tokenCheker: tokenCheker,
	}
}

type contextKey string

const contextLoginKey contextKey = "LoginKey"

func (m *Middleware) WithAuthentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		pc := "func (m *Middleware) WithAuthentication(h http.Handler) http.Handler"

		tokenCookie, err := r.Cookie(authorizationCookie)
		if err != nil {
			m.logger.Infow(
				pc,
				"message", "no authentication token found",
				"error", err,
			)
			http.Error(w, errors.New("no authentication token found").Error(), http.StatusUnauthorized)
			return
		}

		login, err := m.tokenCheker.CheckValidJWT(tokenCookie.Value)
		if err != nil {
			m.logger.Infow(
				pc,
				"message", "token failed",
				"error", err,
			)
			http.Error(w, errors.New("token failed").Error(), http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextLoginKey, login)))
	})
}


type (

	//	Date from response
	responseData struct {
		status int
		size   int
		data   string
	}

	//	Response writer with login
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

// ------------------------------------------------------------
//
//	Override Write() method by http.ResponseWriter
//	Receiver:
//		r *loggingResponseWriter
//	Input:
//		b []byte
//	Output:
//		int - count of write bytes
//		err
func (r *loggingResponseWriter) Write(b []byte) (int, error) {

	size, err := r.ResponseWriter.Write(b)

	r.responseData.data = string(b)
	r.responseData.size += size
	return size, err
}

// ------------------------------------------------------------
//
//	Override WriteHeader() method by http.ResponseWriter
//	Receiver:
//		r *loggingResponseWriter
//	Input:
//		statusCode int
func (r *loggingResponseWriter) WriteHeader(statusCode int) {

	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func (m *Middleware) WithLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		responseData := new(responseData)
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		buf, _ := io.ReadAll(r.Body)
		rdr1 := io.NopCloser(bytes.NewBuffer(buf))	
		r.Body = rdr1

		m.logger.Infow("Information about request",
			"Request URI", r.RequestURI,
			"Request method", r.Method,
			"Request headers", r.Header,
			"Request cookie", r.Cookies(),
			"Request body", string(buf),
		)

		start := time.Now()
		h.ServeHTTP(&lw, r)
		duration := time.Since(start)

		m.logger.Infow(
			"Information about response",
			"Response status", responseData.status,
			"Response headers", lw.ResponseWriter.Header(),
			"Response cookie", r.Cookies(),
			"Response data", responseData.data,

			"duration", duration,
		)

	})
}