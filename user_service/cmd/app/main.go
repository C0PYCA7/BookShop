package main

import (
	"BookShop/user_service/internal/config"
	"BookShop/user_service/internal/database/postgres"
	"BookShop/user_service/internal/handler/user/login"
	"BookShop/user_service/internal/handler/user/registration"
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
		log.Error("init database: ", err)
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	fs := http.FileServer(http.Dir("user_service/web/static/jss"))
	router.Handle("/jss/*", http.StripPrefix("/jss", fs))

	router.Get("/login", login.LoginPage)
	router.Post("/login", login.New(log, database, cfg.Jwt))
	router.Get("/registration", registration.RegistrationPage)
	router.Post("/registration", registration.New(log, database))

	router.Group(func(router chi.Router) {})

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("start server: ", err)
		os.Exit(1)
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
