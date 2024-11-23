package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"machineIssuerSystem/internal/model"
	"machineIssuerSystem/pkg/errorlist"
)

func NewToken(user model.User, secretKey string, ttl time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.UUID
	claims["email"] = user.Email
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(ttl).Unix()

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string, secretKey string) (jwt.MapClaims, error) {
	fmt.Printf("token: %s\n", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errorlist.ErrUnexpectedSigningMethod
		}

		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errorlist.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errorlist.ErrInvalidTokenClaims
	}

	return claims, nil
}
