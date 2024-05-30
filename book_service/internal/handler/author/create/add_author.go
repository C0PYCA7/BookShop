package create

import (
	"BookShop/book_service/internal/config"
	"BookShop/book_service/internal/database/postgres"
	"BookShop/book_service/internal/lib/jwt"
	"BookShop/book_service/internal/model"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type Response struct {
	Id     int    `json:"id"`
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

type CreateAuthor interface {
	AddAuthor(author *model.AddAuthor) (int, error)
}

func NewAuthorPage(w http.ResponseWriter, r *http.Request) {
	log.Println("header from request: ", r.Header.Get("Authorization"))
	http.ServeFile(w, r, "book_service/web/template/newauthor.html")
}

func New(log *slog.Logger, create CreateAuthor, cfg config.JwtConfig) http.HandlerFunc {
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

		var req model.AddAuthor

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body: ", err)
			_, err = fmt.Fprintln(logFile, "NEWAUTHOR invalid request body")
			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "failed to decode request body",
			})

			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request")
			_, err = fmt.Fprintln(logFile, "NEWAUTHOR invalid request")
			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "Invalid request",
			})

			return
		}

		log.Info("validate success")

		id, err := create.AddAuthor(&req)
		if err != nil {
			if errors.Is(err, postgres.ErrAuthorExists) {
				log.Error("author exists")
				_, err = fmt.Fprintln(logFile, fmt.Sprintf("NEWAUTHOR author exists: %s %s", req.Name, req.Surname))
				render.JSON(w, r, Response{
					Status: http.StatusConflict,
					Error:  "author exists",
				})

				return
			}

			log.Error("failed to add author")
			_, err = fmt.Fprintln(logFile, fmt.Sprintf("NEWAUTHOR failed to add author: %s %s", req.Name, req.Surname))
			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "internal server error",
			})

			return

		}

		authHeader := r.Header.Get("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		uid := jwt.GetData(tokenString, cfg)

		date := time.Now().Format("2006-01-02 15:04:05")

		data := fmt.Sprintf("NEWAUTHOR: [%s] user: with id:%s create author:%d", date, uid, id)
		_, err = fmt.Fprintf(file, data)
		if err != nil {
			log.Error("Failed to write file", slog.Any("data", data), slog.String("err", err.Error()))
		}
		_, _ = fmt.Fprintf(file, "\n")

		log.Info("add author ", id)
		_, err = fmt.Fprintln(logFile, fmt.Sprintf("NEWAUTHOR with id %d created", id))

		render.JSON(w, r, Response{
			Id:     id,
			Status: http.StatusOK,
		})

	}
}
