package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret = []byte("your_secret_key")

func GenerateToken(userID uint, isRefreshToken bool) (string, error) {
	days := 24 * time.Hour
	if isRefreshToken {
		days = 30 * 24 * time.Hour
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(days).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (uint, error) {
	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return uint(claims["user_id"].(float64)), nil
	}
	return 0, jwt.ErrSignatureInvalid
}
