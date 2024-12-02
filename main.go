package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/iypetrov/go-indie-hacking-starter/database"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := NewConfig()
	logger := NewLogger()

	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSL,
	)
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		panic(err)
	}
	queries := database.New(conn)

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error("failed to set dialect: %s", err.Error())
	}
	if err := goose.Up(conn, "sql/migrations"); err != nil {
		logger.Error("failed to run migrations: %s", err.Error())
	}

	formDecoder := form.NewDecoder()
	formValidator := validator.New(validator.WithRequiredStructEnabled())

	h := Handler{
		formDecoder:   formDecoder,
		formValidator: formValidator,
		conn:          conn,
		queries:       queries,
	}

	mux := chi.NewRouter()
	mux.Handle("/*", static(logger))
	mux.Route("/p", func(mux chi.Router) {
		mux.Route("/public", func(mux chi.Router) {
			mux.Get("/home", h.HomeView)
		})
		mux.Route("/client", func(mux chi.Router) {
			// No handlers yet
		})
		mux.Route("/admin", func(mux chi.Router) {
			// No handlers yet
		})
	})
	mux.Route("/api", func(mux chi.Router) {
		mux.Route("/public/v0", func(mux chi.Router) {
			mux.Route("/mailing-list", func(mux chi.Router) {
				mux.Post("/", MakeTemplHandler(ctx, logger, h.AddMailToMailingList))
			})
		})

		mux.Route("/client/v0", func(mux chi.Router) {
			// No handlers yet
		})
		mux.Route("/admin/v0", func(mux chi.Router) {
			// No handlers yet
		})
	})
	mux.Get("/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{})
	})
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      mux,
	}
	logger.Info("server started on %s", cfg.App.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("cannot start server: %s", err.Error())
	}
}
