package gofermart

import (
	"testing"
	"time"

	"github.com/Schalure/gofermart/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_isOrderValid(t *testing.T) {

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	stor := mocks.NewMockStorager(mockController)
	logger := mocks.NewMockLoggerer(mockController)
	orderChecker := mocks.NewMockOrderChecker(mockController)
	service := NewGofermart(stor, logger, orderChecker, `[0-9a-zA-Z@._]`, `[0-9a-zA-Z]`, `[0-9]`, time.Hour*1)

	testCases := []struct {
		name      string
		inpString string
		want      bool
	}{
		{
			name:      "simple even test",
			inpString: "4561261212345467",
			want:      true,
		},
		{
			name:      "bad even test",
			inpString: "4561261212345464",
			want:      false,
		}, {
			name:      "simple odd test",
			inpString: "1234567897",
			want:      true,
		},
		{
			name:      "bad odd test",
			inpString: "1234547897",
			want:      false,
		},
		{
			name:      "empty seq test",
			inpString: "",
			want:      false,
		},
		{
			name:      "bad seq test",
			inpString: "4561261212%45467",
			want:      false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			ok := service.isOrderValid(test.inpString)
			assert.EqualValues(t, test.want, ok)
		})
	}
}
