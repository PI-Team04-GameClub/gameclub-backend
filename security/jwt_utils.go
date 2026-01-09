package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("your-secret-key-change-this-in-production")

func GenerateToken(id uint, email, firstName, lastName string) (string, error) {
	claims := jwt.MapClaims{
		"id":         id,
		"email":      email,
		"first_name": firstName,
		"last_name":  lastName,
		"exp":        time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

func ExtractClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenMalformed
	}

	return claims, nil
}
