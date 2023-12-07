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
	tokenTTL time.Duration
	secretKey string
}

//	Service logging interface
type Loggerer interface {
	Info(args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
}

//	Interface of work with the repository
type Storager interface {
	AddNewUser(ctx context.Context, user  storage.User) error
	GetUserByLogin(ctx context.Context, login string) (storage.User, error)
}

//	Constructor of gofermart service object
func NewGofermart(config *configs.Config, s Storager, l Loggerer) *Gofermart {

	validPassword := regexp.MustCompile(`^` + config.AppConfig.PassRules + `+$`)

	return &Gofermart{
		storager: s,
		loggerer: l,

		validPassword: validPassword,
		tokenTTL: config.AppConfig.TokenTTL,
		secretKey: defaultSecretKey,
	}
}



