package get

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

		logFile, err := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Error("open file: ", slog.Any("err", err))
		}
		defer logFile.Close()

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("failed to parse id")
			_, err = fmt.Fprintln(logFile, "GETBOOK failed to parse id")
			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "internal server error",
			})

			return
		}

		rBook, err := book.GetBookInfo(id)
		if err != nil {
			log.Error("failed to get book info")
			_, err = fmt.Fprintln(logFile, fmt.Sprintf("GETBOOK failed to get book info: %d", id))
			render.JSON(w, r, Response{Status: http.StatusInternalServerError, Error: "internal server error"})

			return
		}

		log.Info("got book info")

		render.JSON(w, r, Response{Book: *rBook, Status: http.StatusOK})
	}
}
