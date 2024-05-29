package check

import (
	"BookShop/book_service/internal/model"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

type Response struct {
	model.AuthorInfo `json:"author"`
	Status           int    `json:"status"`
	Error            string `json:"error,omitempty"`
}

type FindAuthor interface {
	GetAuthor(id int) (*model.AuthorInfo, error)
}

func ServeAuthorPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "book_service/web/template/author.html")
}

func New(log *slog.Logger, find FindAuthor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logFile, err := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Error("open file: ", slog.Any("err", err))
		}
		defer logFile.Close()

		idStr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("failed to parse id")
			_, err = fmt.Fprintln(logFile, "")
			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "internal server error",
			})

			return
		}

		resp, err := find.GetAuthor(id)
		if err != nil {
			log.Error("failed to get author&book info")
			_, err = fmt.Fprintln(logFile, fmt.Sprintf("CHECK AUTHOR failed to get author: %d info", id))
			render.JSON(w, r, Response{Status: http.StatusInternalServerError, Error: "internal server error"})

			return
		}

		log.Info("resp: ", resp)

		render.JSON(w, r, Response{
			AuthorInfo: *resp,
			Status:     http.StatusOK,
		})
	}
}
