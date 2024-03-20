package main

import (
	"BookShop/book_service/internal/config"
	"BookShop/book_service/internal/database/postgres"
	"BookShop/book_service/internal/handler/author/check"
	"BookShop/book_service/internal/handler/author/create"
	delete2 "BookShop/book_service/internal/handler/author/delete"
	"BookShop/book_service/internal/handler/book/add"
	delete3 "BookShop/book_service/internal/handler/book/delete"
	"BookShop/book_service/internal/handler/book/get"
	"BookShop/book_service/internal/handler/book/list"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := newLogger()

	database, err := postgres.New(cfg.Database)
	if err != nil {
		log.Error("failed to init database")
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	fs := http.FileServer(http.Dir("book_service/web/static/js"))
	router.Handle("/js/*", http.StripPrefix("/js", fs))

	router.Get("/newauthor", create.NewAuthorPage)
	router.Post("/newauthor", create.New(log, database))

	router.Get("/author/{id}/page", check.ServeAuthorPage)
	router.Get("/author/{id}", check.New(log, database))
	router.Delete("/author/{id}", delete2.New(log, database))

	router.Get("/book/{id}/page", get.ServeBookPage)
	router.Get("/book/{id}", get.New(log, database))
	router.Delete("/book/{id}", delete3.New(log, database))

	router.Get("/newbook", add.NewBookPage)
	router.Post("/newbook", add.New(log, database))

	router.Get("/", list.ServeListPage)
	router.Get("/list", list.New(log, database))
	//todo: общая страничка со списком книг, создать, найти, удалить

	_ = database

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server is stopped")
}

func newLogger() *slog.Logger {
	var log *slog.Logger

	log = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	return log
}
