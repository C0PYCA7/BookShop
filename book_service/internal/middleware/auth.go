package middleware

import (
	"BookShop/book_service/internal/config"
	"BookShop/book_service/internal/lib/jwt"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler, cfg config.JwtConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		if !jwt.Verify(tokenString, cfg) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
