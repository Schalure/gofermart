// The package describes the entire business logic of the service
package gofermart

import (
	"context"
	"regexp"
	"time"

	"github.com/Schalure/gofermart/internal/configs"
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
	AddNewUser(ctx context.Context, user  storage.User) error
	GetUserByLogin(ctx context.Context, login string) (storage.User, error)
}

//	Constructor of gofermart service object
func NewGofermart(config *configs.Config, s Storager, l Loggerer) *Gofermart {

	validLogin := regexp.MustCompile(`^` + config.AppConfig.LoginRules + `+$`)
	validPassword := regexp.MustCompile(`^` + config.AppConfig.PassRules + `+$`)

	return &Gofermart{
		storager: s,
		loggerer: l,

		validLogin: validLogin,
		validPassword: validPassword,
		tokenTTL: config.AppConfig.TokenTTL,
		secretKey: defaultSecretKey,
	}
}



