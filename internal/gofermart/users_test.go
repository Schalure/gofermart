package gofermart

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Schalure/gofermart/internal/configs"
	"github.com/Schalure/gofermart/internal/gofermart/gofermaterrors"
	"github.com/Schalure/gofermart/internal/loggers"
	"github.com/Schalure/gofermart/internal/storage"
	"github.com/Schalure/gofermart/internal/storage/mockstor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateUser(t *testing.T) {

	testCases := []struct {
		name     string
		login    string
		password string
		want     struct {
			err          error
			passwordHash string
		}
	}{
		{
			name:     "simple test",
			login:    "Misha",
			password: "q1w2e3r4",
			want: struct {
				err          error
				passwordHash string
			}{
				err:          nil,
				passwordHash: "e360f368fcca8779da96cba0267a5b2c523afd2909036f643e38fb6cc451163c",
			},
		},
		{
			name:     "dublicate test",
			login:    "Misha",
			password: "q1w2e3r4",
			want: struct {
				err          error
				passwordHash string
			}{
				err:          gofermaterrors.LoginAlreadyTaken,
				passwordHash: "",
			},
		},
		{
			name:     "empty login test",
			login:    "",
			password: "q1w2e3r4",
			want: struct {
				err          error
				passwordHash string
			}{
				err:          gofermaterrors.InvalidLogin,
				passwordHash: "e360f368fcca8779da96cba0267a5b2c523afd2909036f643e38fb6cc451163c",
			},
		},
		{
			name:     "bad login test",
			login:    "Tema#",
			password: "q1w2e3r4",
			want: struct {
				err          error
				passwordHash string
			}{
				err:          gofermaterrors.InvalidLogin,
				passwordHash: "e360f368fcca8779da96cba0267a5b2c523afd2909036f643e38fb6cc451163c",
			},
		},
		{
			name:     "smol password test",
			login:    "Nikita",
			password: "q1",
			want: struct {
				err          error
				passwordHash string
			}{
				err:          gofermaterrors.PasswordShort,
				passwordHash: "",
			},
		},
		{
			name:     "bad password test",
			login:    "Vova",
			password: "q1w2e3r4%",
			want: struct {
				err          error
				passwordHash string
			}{
				err:          gofermaterrors.PasswordBad,
				passwordHash: "",
			},
		},
	}

	logger := loggers.NewLogger(configs.Debug)
	stor := mockstor.NewStorage()
	service := NewGofermart(stor, logger, `[0-9a-zA-Z@._]`, `[0-9a-zA-Z]`, time.Hour*1)

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := service.CreateUser(ctx, test.login, test.password)

			assert.ErrorIs(t, err, test.want.err)

			if err == nil {

				user, ok := stor.Users[test.login]
				require.True(t, ok, fmt.Sprintf("user not found: %s", test.login))
				assert.Equal(t, test.want.passwordHash, user.Password)
				assert.Equal(t, test.login, user.Login)
			}

			cancel()
		})
	}
}

func Test_UserAuthentication(t *testing.T) {

	testCases := []struct {
		name            string
		loginToSave     string
		passwordToSave  string
		loginToCheck    string
		passwordToCheck string
		want            struct {
			err error
		}
	}{
		{
			name:            "simple test",
			loginToSave:     "Mihail",
			passwordToSave:  "q1w2e3r4",
			loginToCheck:    "Mihail",
			passwordToCheck: "q1w2e3r4",
			want: struct{ err error }{
				err: nil,
			},
		},
		{
			name:            "bad login",
			loginToSave:     "Mihail",
			passwordToSave:  "q1w2e3r4",
			loginToCheck:    "Sasha",
			passwordToCheck: "q1w2e3r4",
			want: struct{ err error }{
				err: gofermaterrors.InvalidLoginPassword,
			},
		}, {
			name:            "bad password",
			loginToSave:     "Mihail",
			passwordToSave:  "q1w2e3r4",
			loginToCheck:    "Sasha",
			passwordToCheck: "",
			want: struct{ err error }{
				err: gofermaterrors.InvalidLoginPassword,
			},
		},
	}

	logger := loggers.NewLogger(configs.Debug)
	stor := mockstor.NewStorage()
	service := NewGofermart(stor, logger, `[0-9a-zA-Z@._]`, `[0-9a-zA-Z]`, time.Hour*1)

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			user := storage.User{
				Login:    test.loginToSave,
				Password: service.generatePasswordHash(test.passwordToSave),
			}
			stor.AddNewUser(context.Background(), user)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := service.AuthenticationUser(ctx, test.loginToCheck, test.passwordToCheck)

			assert.Equal(t, test.want.err, err)

			cancel()
		})
	}
}
