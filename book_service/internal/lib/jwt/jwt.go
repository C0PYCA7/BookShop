package jwt

import (
	"BookShop/book_service/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
)

func Verify(tokenString string, config config.JwtConfig) bool {
	secretKey := []byte(config.SecretKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return false
	}

	if !token.Valid {
		return false
	}
	return true
}

func GetData(tokenString string, config config.JwtConfig) string {
	secretKey := []byte(config.SecretKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return ""
	}

	if !token.Valid {
		return ""
	}
	claims := token.Claims.(jwt.MapClaims)
	uid := claims["uid"].(float64)
	return strconv.Itoa(int(uid))
}
