package auth_lib

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUnexpectedSigningMethod = fmt.Errorf("unexpected signing method")
	ErrTokenInvalid            = fmt.Errorf("invalid token")
	ErrTokenExpired            = fmt.Errorf("token expired")
)

type Token struct {
	UID  string
	Type string
	jwt.RegisteredClaims
}

const TokenExpireLife = 200

func GenerateToken(uid string, duration time.Duration, meta string) (string, error) {
	claims := Token{
		uid,
		meta,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokens string) (*Token, error) {
	tkn, err := jwt.ParseWithClaims(tokens, &Token{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		if errors.Is(err, ErrUnexpectedSigningMethod) {
			return nil, ErrTokenInvalid
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		} else if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, ErrTokenInvalid
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, ErrTokenInvalid
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenInvalid
		}
		return nil, err
	}

	if claims, ok := tkn.Claims.(*Token); ok && tkn.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}
