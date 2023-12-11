package gofermart

import (
	"testing"
	"time"

	"github.com/Schalure/gofermart/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_isOrderValid(t *testing.T) {
	
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	stor := mocks.NewMockStorager(mockController)
	logger := mocks.NewMockLoggerer(mockController)

	service := NewGofermart(stor, logger, `[0-9a-zA-Z@._]`, `[0-9a-zA-Z]`, `[0-9]`, time.Hour*1)

	ok := service.isOrderValid("123456789")

	require.NotEqualValues(t, true, ok)
}