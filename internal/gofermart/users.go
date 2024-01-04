package gofermart

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Schalure/gofermart/internal/storage"
	"github.com/golang-jwt/jwt/v4"
)

// Create new user. Return JSON Web Token (JWT)
func (g *Gofermart) CreateUser(ctx context.Context, login, password string) (string, error) {

	pc := "func (g *Gofermart) CreateUser(ctx context.Context, login, password string) error"

	if err := g.isLoginValid(login); err != nil {
		g.loggerer.Debugw(
			pc,
			"error", err,
			"login", login,
		)
		return "", ErrInvalidLogin
	}

	if _, err := g.storager.GetUserByLogin(ctx, login); err == nil {
		g.loggerer.Debugw(
			pc,
			"error", err,
		)
		return "", ErrLoginAlreadyTaken
	}

	if err := g.isPasswordValid(password); err != nil {
		g.loggerer.Debugw(
			pc,
			"error", err,
		)
		return "", err
	}

	user := storage.User{
		Login:    login,
		Password: g.generatePasswordHash(password),
	}

	if err := g.storager.AddNewUser(ctx, user); err != nil {
		g.loggerer.Debugw(
			pc,
			"error", err,
		)
		return "", ErrInternal
	}

	token, err := g.generateJWT(login)
	if err != nil {
		g.loggerer.Debugw(
			pc,
			"error", err,
		)
		return "", ErrInternal
	}

	g.loggerer.Debugw(
		pc,
		"create new user", login,
		"token", token,
	)

	return token, nil
}

// Authentication process for user. Return JSON Web Token (JWT)
func (g *Gofermart) AuthenticationUser(ctx context.Context, login, password string) (string, error) {

	pc := "func (g *Gofermart) AuthenticationUser(ctx context.Context, login, password string) (string, error)"

	user, err := g.storager.GetUserByLogin(ctx, login)
	if err != nil {
		g.loggerer.Debugw(
			pc,
			"error", err,
		)
		return "", ErrInvalidLoginPassword
	}

	if g.generatePasswordHash(password) != user.Password {
		g.loggerer.Debugw(
			pc,
			"error", "the password hash didn't match",
		)
		return "", ErrInvalidLoginPassword
	}

	return g.generateJWT(login)
}

// Return user info
func (g *Gofermart) GetUserInfo(ctx context.Context, login string) (storage.User, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	user, err := g.storager.GetUserByLogin(ctx, login)
	return user, err
}

// Claims of JSON Web Token (JWT)
type Claims struct {
	jwt.RegisteredClaims
	Login string
}

// Check valid JSON Web Token (JWT)
func (g *Gofermart) CheckValidJWT(tokenString string) (string, error) {

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(g.secretKey), nil
	})
	if err != nil {
		return "", errors.New("can't parse token string")
	}
	if !token.Valid {
		return "", errors.New("token not valid")
	}

	if claims.RegisteredClaims.ExpiresAt.Time.Before(time.Now()) {
		return "", fmt.Errorf("token obsolete: token time = %s, now time = %s",
			claims.RegisteredClaims.ExpiresAt.Time.Format(time.RFC3339),
			time.Now().Format(time.RFC3339),
		)
	}

	return claims.Login, nil
}

// Generate JSON Web Token (JWT)
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

	return tokenString, nil
}

// Check to valid login
func (g *Gofermart) isLoginValid(login string) error {

	if len(login) == 0 {
		return ErrInvalidLogin
	}

	if !g.validLogin.MatchString(login) {
		return ErrInvalidLogin
	}

	return nil
}

// Check to valid password
func (g *Gofermart) isPasswordValid(password string) error {

	var errs []error

	if len(password) < PasswordMinLenght {
		errs = append(errs, ErrPasswordShort)
	}

	if !g.validPassword.MatchString(password) {
		errs = append(errs, ErrPasswordBad)
	}

	return errors.Join(errs...)
}

// Generate password hash
func (g *Gofermart) generatePasswordHash(password string) string {

	salt := "m1xFdMsf"
	out := []byte(strings.Join([]string{password, salt}, ""))

	hash := sha256.Sum256(out)
	return fmt.Sprintf("%x", hash[:])
}
