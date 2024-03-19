package delete

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type DeleteBook interface {
	DelBook(id int) error
}

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

func New(log *slog.Logger, book DeleteBook) http.HandlerFunc {
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

		err = book.DelBook(id)
		if err != nil {
			log.Error("failed to delete book")

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "failed to delete book",
			})

			return
		}

		log.Info("delete success")
		render.JSON(w, r, Response{
			Status: http.StatusOK,
		})
	}
}
