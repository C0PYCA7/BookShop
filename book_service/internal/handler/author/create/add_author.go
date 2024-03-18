package create

import (
	"BookShop/book_service/internal/database/postgres"
	"BookShop/book_service/internal/model"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Response struct {
	Id     int    `json:"id"`
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

type CreateAuthor interface {
	AddAuthor(author *model.AddAuthor) (int, error)
}

func New(log *slog.Logger, create CreateAuthor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.AddAuthor

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body")

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "failed to decode request body",
			})

			return
		}

		log.Info("decoded: ", req)

		if err := validator.New().Struct(req); err != nil {
			log.Error("failed to decode request body")

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "failed to decode request body",
			})

			return
		}

		log.Info("validate success")

		id, err := create.AddAuthor(&req)
		if err != nil {
			if errors.Is(err, postgres.ErrAuthorExists) {
				log.Error("author exists")

				render.JSON(w, r, Response{
					Status: http.StatusConflict,
					Error:  "author exists",
				})

				return
			}

			log.Error("failed to add author")

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "internal server error",
			})

			return

		}

		log.Info("add author ", id)

		render.JSON(w, r, Response{
			Id:     id,
			Status: http.StatusOK,
		})

	}
}
