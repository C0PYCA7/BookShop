package delete

import (
	"BookShop/book_service/internal/config"
	"BookShop/book_service/internal/lib/jwt"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type DeleteBook interface {
	DelBook(id int) error
}

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

func New(log *slog.Logger, book DeleteBook, cfg config.JwtConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		file, err := os.OpenFile("data.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Error("open file: ", err)
		}
		defer file.Close()

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

		authHeader := r.Header.Get("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		uid := jwt.GetData(tokenString, cfg)

		date := time.Now().Format("2006-01-02 15:04:05")

		data := fmt.Sprintf("DELETE: [%s] user: with id:%s delete book: %s", date, uid, idStr)
		_, err = fmt.Fprintf(file, data)
		if err != nil {
			log.Error("Failed to write file", slog.Any("data", data), slog.String("err", err.Error()))
		}
		_, _ = fmt.Fprintf(file, "\n")

		log.Info("delete success")
		render.JSON(w, r, Response{
			Status: http.StatusOK,
		})
	}
}
