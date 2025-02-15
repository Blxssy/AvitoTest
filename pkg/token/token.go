package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenGenerator interface {
	NewToken(username string) (string, error)
	ParseToken(tokenString string) (string, error)
}

type TokenGen struct {
	tokenKey string
	tokenTTL time.Duration
}

type TokenConfig struct {
	TokenKey string
	TokenTTL time.Duration
}

func NewTokenGen(cfg TokenConfig) TokenGenerator {
	return &TokenGen{
		tokenKey: cfg.TokenKey,
		tokenTTL: cfg.TokenTTL,
	}
}

func (t *TokenGen) NewToken(username string) (string, error) {
	//expirationTime := time.Now().Add(t.tokenTTL)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		//"exp":      expirationTime,
		"username": username,
	})

	token, err := claims.SignedString([]byte(t.tokenKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (t *TokenGen) ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKeyType
		}
		return []byte(t.tokenKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", jwt.ErrTokenInvalidClaims
		}
		return username, nil
	}

	return "", errors.New("invalid token")
}
