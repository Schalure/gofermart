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
	defaultSecretKey = `poKdq834nFElq71`
)


//	Main service object struct
type Gofermart struct {
	storager Storager
	loggerer Loggerer

	validPassword *regexp.Regexp
	validLogin *regexp.Regexp
	tokenTTL time.Duration
	secretKey string
}

//go:generate mockgen -destination=../mocks/mock_loggerer.go -package=mocks github.com/Schalure/gofermart/internal/gofermart Loggerer
type Loggerer interface {
	Info(args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
}

//	Interface of work with the repository
//go:generate mockgen -destination=../mocks/mock_storager.go -package=mocks github.com/Schalure/gofermart/internal/gofermart Storager
type Storager interface {
	AddNewUser(ctx context.Context, user storage.User) error
	GetUserByLogin(ctx context.Context, login string) (storage.User, error)
}

//	Constructor of gofermart service object
func NewGofermart(s Storager, l Loggerer, loginRules, passRules string, tokenTTL time.Duration) *Gofermart {

	validLogin := regexp.MustCompile(`^` + loginRules + `+$`)
	validPassword := regexp.MustCompile(`^` + passRules + `+$`)

	return &Gofermart{
		storager: s,
		loggerer: l,

		validLogin: validLogin,
		validPassword: validPassword,
		tokenTTL: tokenTTL,
		secretKey: defaultSecretKey,
	}
}



