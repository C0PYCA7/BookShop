package add

import (
	"BookShop/book_service/internal/database/postgres"
	"BookShop/book_service/internal/model"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type AddBook interface {
	AddBook(books *model.AddBook) (int, error)
}

type Response struct {
	Id     int    `json:"id"`
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

func NewBookPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "book_service/web/template/newbook.html")
}

func New(log *slog.Logger, book AddBook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req model.AddBook

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body: ", err)

			render.JSON(w, r, Response{Status: http.StatusBadRequest, Error: "failed to decode request body"})

			return
		}

		log.Info("decode success")

		if err := validator.New().Struct(&req); err != nil {
			log.Error("invalid request")

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "Invalid request",
			})
			return
		}

		log.Info("validate success")

		id, err := book.AddBook(&req)
		if err != nil {
			if errors.Is(err, postgres.ErrAuthorNotFound) {
				log.Error("author not found")

				render.JSON(w, r, Response{Status: http.StatusNotFound, Error: "author not found"})

				return
			}
			log.Error("failed to insert book")

			render.JSON(w, r, Response{Status: http.StatusInternalServerError, Error: "failed to insert data"})

			return
		}

		log.Info("insert success")

		render.JSON(w, r, Response{
			Id:     id,
			Status: http.StatusOK,
		})
	}
}
