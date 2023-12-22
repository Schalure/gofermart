package postgrestor

import (
	"context"
	"testing"
	"time"

	"github.com/Schalure/gofermart/internal/storage"
	"github.com/jackc/pgx/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const connectionString = "host=localhost user=aleksandr password=c1f2i3f4 dbname=gofermartdb sslmode=disable"

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
		order storage.Order
	}{
		{
			name: "simple test",
			user: storage.User{
				Login: "Petya",
				Password: "12345678",
			},
			order: storage.Order{
				OrderNumber: "1234567897",
				OrderStatus: storage.OrderStatusNew,
				UploadedOrder:  pgtype.Timestamptz{
					Time: time.Date(2020, 12, 10, 15, 12, 1, 0, time.Local),
					Status: pgtype.Present,
				},
				BonusPoints: 0,
				UploadedBonus: pgtype.Timestamptz{Status: pgtype.Null,},
				UserLogin: "Petya",

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

			err = stor.AddNewOrder(context.Background(), test.order)
			require.NoError(t, err, "can't get order: %s", test.order.String())

			order, err := stor.GetOrderByNumber(context.Background(), test.order.OrderNumber)
			require.NoError(t, err, "can't get order: %s", test.order.OrderNumber)
			assert.EqualValues(t, test.order, order)
		})
	}
}