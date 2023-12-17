// The package describes the entire business logic of the service
package gofermart

import (
	"context"
	"regexp"
	"time"

	"github.com/Schalure/gofermart/internal/storage"
)

const (
	PasswordMinLenght = 8
	defaultSecretKey  = `poKdq834nFElq71`
)

// Main service object struct
type Gofermart struct {
	storager Storager
	loggerer Loggerer

	orderChecker OrderChecker
	workChannel chan string
	doneCh chan struct{}
	
	validPassword *regexp.Regexp
	validLogin    *regexp.Regexp
	validOrderNumber *regexp.Regexp
	tokenTTL      time.Duration
	secretKey     string
}


//go:generate mockgen -destination=../mocks/mock_orderchecker.go -package=mocks github.com/Schalure/gofermart/internal/gofermart OrderChecker
type OrderChecker interface {
	OrderCheck(ctx context.Context, ordernumber string) (storage.Order, int)
}


//go:generate mockgen -destination=../mocks/mock_loggerer.go -package=mocks github.com/Schalure/gofermart/internal/gofermart Loggerer
type Loggerer interface {
	Info(args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
}


//go:generate mockgen -destination=../mocks/mock_storager.go -package=mocks github.com/Schalure/gofermart/internal/gofermart Storager
type Storager interface {
	AddNewUser(ctx context.Context, user storage.User) error
	GetUserByLogin(ctx context.Context, login string) (storage.User, error)
	AddNewOrder(ctx context.Context, order storage.Order) error
	GetOrderByNumber(ctx context.Context, orderNumber string) (storage.Order, error)
	GetOrdersByLogin(ctx context.Context, login string) ([]storage.Order, error)
	GetOrdersToUpdateStatus(ctx context.Context) ([]storage.Order, error)
	WithdrawPointsForOrder(ctx context.Context, orderNumber string, sum int) error
	GetPointWithdraws(ctx context.Context, login string) ([]storage.Order, error)
}

// Constructor of gofermart service object
func NewGofermart(s Storager, l Loggerer, orderChecker OrderChecker, loginRules, passRules, OrderNumberRules string, tokenTTL time.Duration) *Gofermart {

	validLogin := regexp.MustCompile(`^` + loginRules + `+$`)
	validPassword := regexp.MustCompile(`^` + passRules + `+$`)
	validOrderNumber := regexp.MustCompile(`^` + OrderNumberRules + `+$`)

	return &Gofermart{
		storager: s,
		loggerer: l,

		orderChecker: orderChecker,

		validLogin:    validLogin,
		validPassword: validPassword,
		validOrderNumber: validOrderNumber,
		tokenTTL:      tokenTTL,
		secretKey:     defaultSecretKey,
	}
}

//	Start service workers and other tasks
func (g *Gofermart) Run(ctx context.Context) {

	//g.orderProvider.Run(ctx)
}
