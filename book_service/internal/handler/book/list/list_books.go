package list

import (
	"BookShop/book_service/internal/model"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type ListBook interface {
	GelAllBooks() ([]model.Book, error)
}

type Response struct {
	Books  []model.Book `json:"books"`
	Status int          `json:"status"`
	Error  string       `json:"error,omitempty"`
}

func New(log *slog.Logger, book ListBook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := book.GelAllBooks()
		if err != nil {
			log.Error("failed to get all books")

			render.JSON(w, r, Response{Status: http.StatusInternalServerError, Error: "failed to get all books"})
		}
	}
}
