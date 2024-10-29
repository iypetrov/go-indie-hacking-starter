package router

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iypetrov/go-indie-hacking-starter/config"
	"github.com/iypetrov/go-indie-hacking-starter/internal/common"
	"github.com/iypetrov/go-indie-hacking-starter/template/view/client"
)

func New(ctx context.Context, static http.Handler) *chi.Mux {
	mux := chi.NewRouter()

	// files & images
	mux.Handle("/*", static)

	// views
	mux.Route(config.Get().ViewPrefix(), func(mux chi.Router) {
		// public views
		mux.Route(config.Get().PublicViewPrefix(), func(mux chi.Router) {
		})

		// client views
		mux.Route(config.Get().ClientViewPrefix(), func(mux chi.Router) {
			mux.Get("/home", func(w http.ResponseWriter, r *http.Request) {
				common.Render(w, r, client.Home())
			})
		})

		// public views
		mux.Route(config.Get().AdminViewPrefix(), func(mux chi.Router) {
		})
	})

	// health check
	mux.Get("/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{})
	})

	return mux
}
