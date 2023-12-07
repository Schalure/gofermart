package gofermart

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/Schalure/gofermart/internal/storage"
	"github.com/golang-jwt/jwt/v4"
)

//	Create new user
func (g *Gofermart) CreateUser(ctx context.Context, login, password string) error {

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

	user := storage.User {
		Login: login,
		Password: g.generatePasswordHash(password),
	}

	return g.storager.AddNewUser(ctx, user)
}

//	Authentication process for user. Return JSON Web Token (JWT)
func (g *Gofermart) AuthenticationUser(ctx context.Context, login, password string) (string, error) {

	pc := "func (g *Gofermart) AuthenticationUser(ctx context.Context, login, password string) (string, error)"

	user, err := g.storager.GetUserByLogin(ctx, login)
	if err != nil {
		g.loggerer.Debugw(
			pc,
			"error", err,
		)
		return "", fmt.Errorf("invalid login or password")
	}

	if g.generatePasswordHash(password) != user.Password {
		g.loggerer.Debugw(
			pc,
			"error", "the password hash didn't match",
		)
		return "", fmt.Errorf("invalid login or password")
	}

	return g.generateJWT(login)
}

//	Claims of JSON Web Token (JWT)
type Claims struct {
	jwt.RegisteredClaims
	Login string
}

//	Generate JSON Web Token (JWT)
func (g *Gofermart) generateJWT(login string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.tokenTTL)),
		},
		Login: login,
	})

	tokenString, err := token.SignedString([]byte(g.secretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
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
