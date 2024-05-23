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
	middleware2 "BookShop/book_service/internal/middleware"
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

	router.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return middleware2.AuthMiddleware(next, cfg.Jwt)
		})
		r.Get("/newauthor", create.NewAuthorPage)
		r.Post("/newauthor", create.New(log, database, cfg.Jwt))
		r.Get("/author/{id}/page", check.ServeAuthorPage)
		r.Get("/author/{id}", check.New(log, database))
		r.Delete("/author/{id}", delete2.New(log, database, cfg.Jwt))

		r.Get("/book/{id}/page", get.ServeBookPage)
		r.Get("/book/{id}", get.New(log, database))
		r.Delete("/book/{id}", delete3.New(log, database, cfg.Jwt))

		r.Get("/newbook", add.NewBookPage)
		r.Post("/newbook", add.New(log, database, cfg.Jwt))
	})

	router.Get("/", list.ServeListPage)
	router.Get("/list", list.New(log, database))

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	log.Info("server started on port: ", slog.String("address", cfg.HttpServer.Address))

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
