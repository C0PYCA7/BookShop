package delete

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

type Request struct {
	Login string `json:"login"`
}

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"err_msg,omitempty"`
}

type UserDeleter interface {
	DeleteUser(login string) error
}

func ShowPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "user_service/web/template/delete.html")
}

func New(log *slog.Logger, deleter UserDeleter, cfg config.JwtConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		date := time.Now().Format("2006-01-02 15:04:05")

		file, err := os.OpenFile("data.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Error("open file: ", slog.Any("err", err))
		}
		defer file.Close()

		logFile, err := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Error("open file: ", slog.Any("err", err))
		}
		defer logFile.Close()

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body")
			_, err = fmt.Fprintln(logFile, "DELETE USER failed to decode request body")
			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "Failed to decode request body",
			})

			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request")
			_, err = fmt.Fprintln(logFile, "DELETE USER invalid request")
			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "Invalid request",
			})

			return
		}

		err = deleter.DeleteUser(req.Login)
		if err != nil {
			if errors.Is(err, postgres.ErrUserNotFound) {
				log.Error("user not found")
				_, err = fmt.Fprintln(logFile, fmt.Sprintf("DELETE USER failed to delete user: %s", req.Login))
				render.JSON(w, r, Response{
					Status: http.StatusNotFound,
					Error:  "User not found",
				})
				return
			}
			log.Error("failed to delete user")
			_, err = fmt.Fprintln(logFile, fmt.Sprintf("DELETE USER failed to delete user %s", req.Login))
			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "Failed to delete user",
			})
			return
		}

		authHeader := r.Header.Get("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		uid := jwt.GetData(tokenString, cfg)

		data := fmt.Sprintf("DELETE USER: [%s] user with id: %s deleted user with login %s", date, uid, req.Login)

		_, err = fmt.Fprintf(file, data)
		if err != nil {
			log.Error("Failed to write file", slog.Any("data", data), slog.String("err", err.Error()))
		}
		_, _ = fmt.Fprintf(file, "\n")
		_, _ = fmt.Fprintf(logFile, "DELETE USER deleted user with login %s", req.Login)

		render.JSON(w, r, Response{
			Status: http.StatusOK,
		})
	}
}
