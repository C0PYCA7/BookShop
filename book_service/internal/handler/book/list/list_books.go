package list

import (
	"BookShop/book_service/internal/model"
	"fmt"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"os"
)

type ListBook interface {
	GelAllBooks() ([]model.Book, error)
}

type Response struct {
	Books  []model.Book `json:"books"`
	Status int          `json:"status"`
	Error  string       `json:"error,omitempty"`
}

func ServeListPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "book_service/web/template/list.html")
}

func New(log *slog.Logger, book ListBook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logFile, err := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Error("open file: ", slog.Any("err", err))
		}
		defer logFile.Close()

		books, err := book.GelAllBooks()
		if err != nil {
			log.Error("failed to get all books")
			_, err = fmt.Fprintf(logFile, "LIST failed to get all books")
			render.JSON(w, r, Response{Status: http.StatusInternalServerError, Error: "failed to get all books"})
		}

		log.Info("got books: ", books)

		render.JSON(w, r, Response{Books: books, Status: http.StatusOK})
	}
}
