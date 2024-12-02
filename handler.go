package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/iypetrov/go-indie-hacking-starter/database"
	"github.com/iypetrov/go-indie-hacking-starter/templates/components"
	"github.com/iypetrov/go-indie-hacking-starter/templates/views"
	"github.com/lib/pq"
)

type Handler struct {
	formDecoder   *form.Decoder
	formValidator *validator.Validate
	conn          *sql.DB
	queries       *database.Queries
}

func (hnd *Handler) HomeView(w http.ResponseWriter, r *http.Request) {
	Render(w, r, views.PublicHome())
}

func (hnd *Handler) AddEmailToMailingList(ctx context.Context, logger Logger, w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		AddToast(w, ErrorInternalServerError(ErrParsingFrom))
		return Render(w, r, components.PublicMailingListForm(components.PublicMailingListFormInput{}))
	}
	var input components.PublicMailingListFormInput
	err = hnd.formDecoder.Decode(&input, r.Form)
	if err != nil {
		AddToast(w, ErrorInternalServerError(ErrDecodingForm))
		return Render(w, r, components.PublicMailingListForm(components.PublicMailingListFormInput{}))
	}

	err = hnd.formValidator.Struct(input)
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

	output, err := hnd.queries.AddEmailToMailingList(
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
}
