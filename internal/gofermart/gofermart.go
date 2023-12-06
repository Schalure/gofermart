// The package describes the entire business logic of the service
package gofermart

import (
	"context"

	"github.com/Schalure/gofermart/internal/configs"
	"github.com/Schalure/gofermart/internal/storage"
)

//	Main service object struct
type Gofermart struct {
	storager Storager
	loggerer Loggerer
}

//	Service logging interface
type Loggerer interface {

}

//	Interface of work with the repository
type Storager interface {
	GetUserByLogin(ctx context.Context, ctxlogin string) (storage.User, error)
}

//	Constructor of gofermart service object
func NewGofermart(config *configs.Config, s Storager, l Loggerer) *Gofermart {

	return &Gofermart{
		storager: s,
		loggerer: l,
	}
}

func (g *Gofermart) CreateUser(ctx context.Context, login, password string) error {

	user, err := g.storager.GetUserByLogin(ctx, login)




}

