package get

import (
	"BookShop/book_service/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type Response struct {
	Book   model.BookInfo `json:"book"`
	Status int            `json:"status"`
	Error  string         `json:"error,omitempty"`
}

type GetBook interface {
	GetBookInfo(id int) (*model.BookInfo, error)
}

func ServeBookPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "book_service/web/template/book.html")
}

func New(log *slog.Logger, book GetBook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("failed to parse id")

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "internal server error",
			})

			return
		}

		rBook, err := book.GetBookInfo(id)
		if err != nil {
			log.Error("failed to get book info")

			render.JSON(w, r, Response{Status: http.StatusInternalServerError, Error: "internal server error"})

			return
		}

		log.Info("got book info")

		render.JSON(w, r, Response{Book: *rBook, Status: http.StatusOK})
	}
}
