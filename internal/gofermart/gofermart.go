// The package describes the entire business logic of the service
package gofermart

import "github.com/Schalure/gofermart/internal/configs"

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

}

//	Constructor of gofermart service object
func NewGofermart(config *configs.Config, s Storager, l Loggerer) *Gofermart {

	return &Gofermart{
		storager: s,
		loggerer: l,
	}
}

