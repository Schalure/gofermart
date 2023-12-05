// The package describes the entire business logic of the service
package gofermart

//	Main service object struct
type Gofermart struct {
	storage *Storage
	logger *Logger
}

//	Service logging interface
type Logger interface {

}

//	Interface of work with the repository
type Storage interface {

}

//	Constructor of gofermart service object
func NewGofermart(s *Storage, l *Logger) (*Gofermart, error) {

	return &Gofermart{
		storage: s,
		logger: l,
	}, nil
}

