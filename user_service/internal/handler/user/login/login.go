package login

import (
	"BookShop/user_service/internal/config"
	"BookShop/user_service/internal/database/postgres"
	"BookShop/user_service/internal/lib/jwt"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
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

type RedisSave interface {
	SaveUserPermissions(id int, permissions string) error
}

//func LoginPage(w http.ResponseWriter, r *http.Request) {
//	tmpl, err := template.ParseFiles("user_service/web/template/login.html")
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	// Отправляем HTML-страницу в ответ на запрос
//	if err := tmpl.Execute(w, nil); err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}

func New(log *slog.Logger, logIn LogIn, cfg config.JwtConfig, red RedisSave) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		err = red.SaveUserPermissions(id, permissions)
		if err != nil {
			log.Error("Internal server error в редисе")

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "Internal server error",
			})

			return
		}

		token, err := jwt.NewToken(id, cfg)
		if err != nil {

			log.Error("Internal server error")

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "Internal server error",
			})

			return
		}

		render.JSON(w, r, Response{Token: token,
			Status: http.StatusOK,
		})
	}
}
