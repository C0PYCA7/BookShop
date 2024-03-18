package main

import (
	"BookShop/book_service/internal/config"
	"BookShop/book_service/internal/database/postgres"
	"BookShop/book_service/internal/handler/author/check"
	"BookShop/book_service/internal/handler/author/create"
	delete2 "BookShop/book_service/internal/handler/author/delete"
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

	router.Delete("/author/{id}", delete2.New(log, database))
	router.Post("/newauthor", create.New(log, database))
	router.Get("/author/{id}", check.New(log, database))

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
