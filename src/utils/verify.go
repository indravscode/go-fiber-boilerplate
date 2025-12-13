//revive:disable:var-naming
package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(tokenStr, secret, tokenType string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(_ *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	jwtType, ok := claims["type"].(string)
	if !ok || jwtType != tokenType {
		return "", errors.New("invalid token type")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid token sub")
	}

	return userID, nil
}
