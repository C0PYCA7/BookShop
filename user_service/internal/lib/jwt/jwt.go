package jwt

import (
	"BookShop/user_service/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

func NewToken(id int, permission string, cfg config.JwtConfig) (string, error) {
	claims := jwt.MapClaims{
		"uid":        id,
		"permission": permission,
		"exp":        time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
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
