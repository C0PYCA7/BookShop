package middleware

import (
	"BookShop/user_service/internal/config"
	"BookShop/user_service/internal/lib/jwt"
	"log"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler, cfg config.JwtConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		log.Println("token: ", authHeader)

		if authHeader == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		ok, perm := jwt.Verify(tokenString, cfg)
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if !perm {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
