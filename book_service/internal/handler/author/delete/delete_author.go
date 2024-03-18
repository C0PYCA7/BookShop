package delete

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

type DeleteAuthor interface {
	DelAuthor(id int) error
}

func New(log *slog.Logger, delete DeleteAuthor) http.HandlerFunc {
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

		err = delete.DelAuthor(id)
		if err != nil {
			log.Error("failed to delete author")

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "failed to delete author",
			})

			return
		}

		log.Info("delete success")

		render.JSON(w, r, Response{Status: http.StatusOK})
	}
}
