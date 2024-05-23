package login

import (
	"BookShop/user_service/internal/config"
	"BookShop/user_service/internal/database/postgres"
	"BookShop/user_service/internal/lib/jwt"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Request struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	Token  string `json:"token,omitempty"`
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

type LogIn interface {
	CheckUser(login string) (int, string, string, error)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "user_service/web/template/login.html")
}

func New(log *slog.Logger, logIn LogIn, cfg config.JwtConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		file, err := os.OpenFile("data.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Error("open file: ", slog.Any("err", err))
		}
		defer file.Close()

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body")

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "Failed to decode request body",
			})

			return
		}

		log.Info("request body decoded: ", req)

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request")

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "Invalid request",
			})

			return
		}

		id, passBd, permissions, err := logIn.CheckUser(req.Login)
		if err != nil {
			if errors.Is(err, postgres.ErrUserNotFound) {
				log.Error("Error", postgres.ErrUserNotFound)

				render.JSON(w, r, Response{
					Status: http.StatusNotFound,
					Error:  "User not found",
				})

				return
			}

			log.Error("Error", postgres.ErrInternalServer)

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "Internal server error",
			})

			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(passBd), []byte(req.Password)); err != nil {
			log.Error("Invalid data")

			render.JSON(w, r, Response{
				Status: http.StatusNotFound,
				Error:  "Invalid data",
			})

			return
		}

		log.Info("user id ", id)

		token, err := jwt.NewToken(id, permissions, cfg)
		if err != nil {

			log.Error("Internal server error")

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "Internal server error",
			})

			return
		}

		r.Header.Set("Authorization", "Bearer "+token)

		date := time.Now().Format("2006-01-02 15:04:05")

		data := fmt.Sprintf("AUTH: [%s] user:%s with id:%d logged in successfully", date, req.Login, id)

		_, err = fmt.Fprintf(file, data)
		if err != nil {
			log.Error("Failed to write file", slog.Any("data", data), slog.String("err", err.Error()))
		}
		_, _ = fmt.Fprintf(file, "\n")

		log.Info("token: ", r.Header.Get("Authorization"))
		render.JSON(w, r, Response{Token: token,
			Status: http.StatusOK,
		})

	}
}
