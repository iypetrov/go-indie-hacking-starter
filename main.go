package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"github.com/iypetrov/go-indie-hacking-starter/templates/views"
)

func Render(w http.ResponseWriter, r *http.Request, c templ.Component) error {
	w.Header().Set("Content-Type", "text/html")

	err := c.Render(r.Context(), w)
	if err != nil {
		return ErrorInternalServerError(fmt.Errorf("server failed to render this component"))
	}

	return nil
}

func Make(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			var t Toast
			if errors.As(err, &t) {
				AddToast(w, t)
			}
		}
	}
}

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := NewConfig()
	logger := NewLogger()
	mux := chi.NewRouter()

	mux.Handle("/*", static())

	mux.Route("/p", func(mux chi.Router) {
		mux.Route("/public", func(mux chi.Router) {
			mux.Get("/home", func(w http.ResponseWriter, r *http.Request) {
				Render(w, r, views.PublicHome())
			})
		})

		mux.Route("/client", func(mux chi.Router) {
		})

		mux.Route("/admin", func(mux chi.Router) {
		})
	})

	mux.Route("/api", func(mux chi.Router) {
		mux.Route("/public/v0", func(mux chi.Router) {
		})

		mux.Route("/client/v0", func(mux chi.Router) {
		})

		mux.Route("/admin/v0", func(mux chi.Router) {
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
