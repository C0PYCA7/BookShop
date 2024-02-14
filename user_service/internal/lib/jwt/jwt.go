package jwt

import (
	"BookShop/user_service/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(id int, cfg config.JwtConfig) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = id
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(cfg.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
