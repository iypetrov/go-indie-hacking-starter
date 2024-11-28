package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/iypetrov/go-indie-hacking-starter/database"
	"github.com/iypetrov/go-indie-hacking-starter/templates/components"
	"github.com/iypetrov/go-indie-hacking-starter/templates/views"
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

	mux := chi.NewRouter()
	mux.Handle("/*", static(logger))
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
			mux.Route("/mailing-list", func(mux chi.Router) {
				mux.Post("/", MakeTemplHandler(func(w http.ResponseWriter, r *http.Request) error {
					err := r.ParseForm()
					if err != nil {
						AddToast(w, ErrorInternalServerError(ErrParsingFrom))
						return Render(w, r, components.PublicMailingListForm(components.PublicMailingListFormInput{}))
					}
					var input components.PublicMailingListFormInput
					err = formDecoder.Decode(&input, r.Form)
					if err != nil {
						AddToast(w, ErrorInternalServerError(ErrDecodingForm))
						return Render(w, r, components.PublicMailingListForm(components.PublicMailingListFormInput{}))
					}

					err = formValidator.Struct(input)
					if err != nil {
						if _, ok := err.(*validator.InvalidValidationError); ok {

							AddToast(w, ErrorInternalServerError(ErrFailedtoValidateRequest))
							return Render(w, r, components.PublicMailingListForm(components.PublicMailingListFormInput{}))
						}

						for _, err := range err.(validator.ValidationErrors) {
							switch err.Field() {
							case "Email":
								if err.Tag() == "required" {
									AddToast(w, WarningStatusBadRequest(WarnEmailIsRequred))
									return Render(w, r, components.PublicMailingListForm(components.PublicMailingListFormInput{Email: input.Email}))
								} else if err.Tag() == "email" {
									AddToast(w, WarningStatusBadRequest(WarnInvalidEmailFormat))
									return Render(w, r, components.PublicMailingListForm(components.PublicMailingListFormInput{Email: input.Email}))
								}
							}
						}
					}

					output, err := queries.AddEmailToMailingList(
						ctx,
						database.AddEmailToMailingListParams{
							ID:        uuid.New(),
							Email:     input.Email,
							CreatedAt: time.Now().UTC(),
						},
					)
					if err != nil {
						var pgErr *pq.Error

						ok := errors.As(err, &pgErr)
						if ok {
							if pgErr.Code == "23505" {
								AddToast(w, WarningStatusBadRequest(WarnEmailAlreadyExists))
								return Render(w, r, components.PublicMailingListForm(components.PublicMailingListFormInput{}))
							}
						}
						AddToast(w, ErrorInternalServerError(ErrFailedToAddEmailToMailingList))
						return Render(w, r, components.PublicMailingListForm(components.PublicMailingListFormInput{}))
					}
					logger.Info("email %s was added to the mailing list", output.Email)

					AddToast(w, SuccessStatusCreated(SuccEmailAddedToMailingList))
					return Render(w, r, components.PublicMailingListForm(components.PublicMailingListFormInput{}))
				}))
			})
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
