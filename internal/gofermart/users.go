package gofermart

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

	"github.com/Schalure/gofermart/internal/storage"
)

//	Create new user
func (g *Gofermart) CreateUser(ctx context.Context, login, password string) error {

	var errs []error

	if _, err := g.storager.GetUserByLogin(ctx, login); err == nil {
		g.loggerer.Debugw(
			"func (g *Gofermart) CreateUser(ctx context.Context, login, password string) error",
			"error", err,
		)
		return fmt.Errorf("a user with this login already exists")
	}

	if err := g.isPasswordValid(password); err != nil {
		return err
	}

	if err := errors.Join(errs...); err != nil {
		g.loggerer.Infow(
			"func (g *Gofermart) CreateUser(ctx context.Context, login, password string) error",
			"error", err,
		)
		return err
	}

	user := storage.User {
		Login: login,
		Password: g.generatePasswordHash(password),
	}

	return g.storager.AddNewUser(ctx, user)
}


//	Check to valid password
func (g *Gofermart) isPasswordValid(password string) error {

	if len(password) < PasswordMinLenght {
		return fmt.Errorf("password is too short")
	}

	if !g.validPassword.MatchString(password) {
		return fmt.Errorf("the password can only be made of characters: 0-9, a-z, A-Z")
	}
	return nil
}

//	Generate password hash
func (g *Gofermart) generatePasswordHash(password string) string {

	salt := "m1xFdMsf"
	out := []byte(strings.Join([]string{password, salt}, ""))

	hash := sha256.Sum256(out)
	return fmt.Sprintf("%x", hash[:])
}
