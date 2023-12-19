package postgrestor

import (
	"context"
	"testing"

	"github.com/Schalure/gofermart/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const connectionString = "user=aleksandr password=qwerty dbname=gofermartdb sslmode=disable"

func Test_User(t *testing.T) {

	stor, err := NewStorage(connectionString)
	require.NoErrorf(t, err, "can't connect to database: %s", connectionString)
	defer func() {
		_, err := stor.db.Exec(context.Background(), `DROP TABLE orders, users;`)
		require.NoError(t, err, "can't drop tables")
	}()

	testCases := []struct {
		name string
		user storage.User
	}{
		{
			name: "simple test",
			user: storage.User{
				Login: "Petya",
				Password: "12345678",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			err := stor.AddNewUser(context.Background(), test.user)
			require.NoError(t, err, "can't save new user: %s", test.user.String())


			user, err := stor.GetUserByLogin(context.Background(), test.user.Login)
			require.NoError(t, err, "can't get user: %s", test.user.Login)
			assert.EqualValues(t, test.user, user)


		})
	}
}