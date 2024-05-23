package delete

import (
	"BookShop/user_service/internal/database/postgres"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

type DeleteUser interface {
	DeleteUser(id int) error
}

func New(log *slog.Logger, user DeleteUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userIdStr := chi.URLParam(r, "id")

		if userIdStr == "" {
			log.Error("Empty url parameter")

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "Empty URL parameter",
			})

			return
		}

		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			log.Error("Invalid url parameter")

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "Invalid URL parameter",
			})

			return
		}

		err = user.DeleteUser(userId)
		if err != nil {
			if errors.Is(err, postgres.ErrUserNotFound) {
				log.Error("User not found")

				render.JSON(w, r, Response{
					Status: http.StatusNotFound,
					Error:  "User not found",
				})

				return
			}
			log.Error("Internal server error")

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "Internal server error",
			})

			return
		}
		log.Info("user delete: ", userId)

		render.JSON(w, r, Response{
			Status: http.StatusOK,
		})
	}
}
