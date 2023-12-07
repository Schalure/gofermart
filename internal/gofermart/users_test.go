package gofermart

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Schalure/gofermart/internal/configs"
	"github.com/Schalure/gofermart/internal/loggers"
	"github.com/Schalure/gofermart/internal/storage/mockstor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateUser(t *testing.T) {

	testCases := []struct{
		name string
		login string
		password string
		want struct {
			err error
			passwordHash string
		}
	}{
		{
			name: "simple test",
			login: "Misha",
			password: "q1w2e3r4",
			want: struct{err error; passwordHash string}{
				err:  nil,
				passwordHash: "e360f368fcca8779da96cba0267a5b2c523afd2909036f643e38fb6cc451163c",
			},
		},
		{
			name: "dublicate test",
			login: "Misha",
			password: "q1w2e3r4",
			want: struct{err error; passwordHash string}{
				err:  fmt.Errorf("a user with this login already exists"),
				passwordHash: "e360f368fcca8779da96cba0267a5b2c523afd2909036f643e38fb6cc451163c",
			},
		},		
		{
			name: "smol password test",
			login: "Nikita",
			password: "q1",
			want: struct{err error; passwordHash string}{
				err:  fmt.Errorf("a user with this login already exists"),
				passwordHash: "e360f368fcca8779da96cba0267a5b2c523afd2909036f643e38fb6cc451163c",
			},
		},
		{
			name: "bad password test",
			login: "Vova",
			password: "q1w2e3r4%",
			want: struct{err error; passwordHash string}{
				err:  fmt.Errorf("a user with this login already exists"),
				passwordHash: "e360f368fcca8779da96cba0267a5b2c523afd2909036f643e38fb6cc451163c",
			},
		},
	}

	config, _ := configs.NewConfig()
	logger := loggers.NewLogger(config)
	stor := mockstor.NewStorage()
	service := NewGofermart(config, stor, logger)

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
			err := service.CreateUser(ctx, test.login, test.password)
			assert.Equal(t, test.want.err, err)

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