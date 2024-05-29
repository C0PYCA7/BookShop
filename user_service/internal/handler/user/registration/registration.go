package registration

import (
	"BookShop/user_service/internal/database/postgres"
	"BookShop/user_service/internal/model"
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
	model.UserRegistration
}

// todo - наверное убрать Id из Response
type Response struct {
	Id     int    `json:"id"`
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

type CreateUser interface {
	CreateUser(user *model.UserRegistration) (int, error)
}

func RegistrationPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "user_service/web/template/registration.html")
}

func New(log *slog.Logger, create CreateUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		var req model.UserRegistration

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("Failed to decode request body")
			_, err = fmt.Fprintln(logFile, "REG failed to decode request body")
			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "Failed to decode request body",
			})

			return
		}

		log.Info("Request decoded ", req)

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request")
			_, err = fmt.Fprintln(logFile, "REG failed to decode request body")
			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "Invalid request",
			})
			return
		}

		password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to hash password")
			_, err = fmt.Fprintln(logFile, "REG failed to hash password")
			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "Internal server error",
			})

			return
		}

		req.Password = string(password)

		id, err := create.CreateUser(&req)
		if err != nil {
			if errors.Is(err, postgres.ErrLoginExists) {
				log.Info("login exists")
				_, err = fmt.Fprintln(logFile, "REG failed to register user")
				render.JSON(w, r, Response{
					Status: http.StatusConflict,
					Error:  "Login exists",
				})

				return
			}
			log.Info("internal error")

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "Internal error",
			})

			return
		}

		date := time.Now().Format("2006-01-02 15:04:05")

		data := fmt.Sprintf("REG: [%s] user:%s with id:%d registered successfully", date, req.Login, id)

		_, err = fmt.Fprintf(file, data)
		if err != nil {
			log.Error("Failed to write file", slog.Any("data", data), slog.String("err", err.Error()))
			_, err = fmt.Fprintln(file, "REG failed to write in file data")
		}
		_, _ = fmt.Fprintf(file, "\n")
		render.JSON(w, r, Response{
			Id:     id,
			Status: http.StatusOK,
		})
	}
}
