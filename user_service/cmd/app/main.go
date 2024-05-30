package main

import (
	"BookShop/user_service/internal/config"
	"BookShop/user_service/internal/database/postgres"
	delete2 "BookShop/user_service/internal/handler/user/delete"
	"BookShop/user_service/internal/handler/user/login"
	"BookShop/user_service/internal/handler/user/registration"
	"BookShop/user_service/internal/handler/user/update"
	middleware3 "BookShop/user_service/internal/middleware"
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

	router.Group(func(router chi.Router) {
		router.Use(func(next http.Handler) http.Handler {
			return middleware3.AuthMiddleware(next, cfg.Jwt)
		})
		router.Get("/update", update.ShowPage)
		router.Post("/update", update.New(log, database, cfg.Jwt))
		router.Get("/delete", delete2.ShowPage)
		router.Delete("/delete", delete2.New(log, database, cfg.Jwt))
	})

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	log.Info("server started on port: ", slog.String("address", cfg.HttpServer.Address))

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
