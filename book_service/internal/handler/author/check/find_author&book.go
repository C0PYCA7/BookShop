package check

import (
	"BookShop/book_service/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type Response struct {
	model.AuthorInfo
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

type FindAuthor interface {
	GetAuthor(id int) (*model.AuthorInfo, error)
}

func New(log *slog.Logger, find FindAuthor) http.HandlerFunc {
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

		resp, err := find.GetAuthor(id)
		if err != nil {
			log.Error("failed to get author&book into")

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
