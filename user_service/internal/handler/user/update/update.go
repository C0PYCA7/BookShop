package update

import (
	"BookShop/user_service/internal/config"
	"BookShop/user_service/internal/database/postgres"
	"BookShop/user_service/internal/lib/jwt"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type Updater interface {
	UpdatePermission(login, permission string) error
}

type Request struct {
	Login      string `json:"login" validate:"required"`
	Permission string `json:"permission" validate:"required"`
}

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"err"`
}

func ShowPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "user_service/web/template/login.html")
}

func New(log *slog.Logger, updater Updater, cfg config.JwtConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		file, err := os.OpenFile("data.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Error("open file: ", slog.Any("err", err))
		}
		defer file.Close()

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("Failed to decode request body")

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "Failed to decode request body",
			})

			return
		}

		log.Info("Request decoded ", req)

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request")

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "Invalid request",
			})
			return
		}

		err = updater.UpdatePermission(req.Login, req.Permission)
		if err != nil {
			if errors.Is(err, postgres.ErrUserNotFound) {
				log.Error("User not found")

				render.JSON(w, r, Response{
					Status: http.StatusNotFound,
					Error:  "User not found",
				})
				return
			}
			log.Error("Failed to update user")

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "Failed to update user",
			})
			return
		}

		authHeader := r.Header.Get("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		uid := jwt.GetData(tokenString, cfg)

		date := time.Now().Format("2006-01-02 15:04:05")

		data := fmt.Sprintf("UPDATE: [%s] %s's permission updated to %s, by user with id:%s", date, req.Login, req.Permission, uid)

		_, err = fmt.Fprintf(file, data)
		if err != nil {
			log.Error("Failed to write file", slog.Any("data", data), slog.String("err", err.Error()))
		}
		_, _ = fmt.Fprintf(file, "\n")

	}
}
