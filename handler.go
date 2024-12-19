package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/bwmarrin/snowflake"
	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
	"github.com/iypetrov/go-indie-hacking-starter/database"
	"github.com/iypetrov/go-indie-hacking-starter/templates/components"
	"github.com/iypetrov/go-indie-hacking-starter/templates/views"
	"github.com/mattn/go-sqlite3"
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

func (hnd *Handler) LoginView(w http.ResponseWriter, r *http.Request) {
	Render(w, r, views.Login())
}

func (hnd *Handler) AddEmailToMailingList(ctx context.Context, logger Logger, sf *snowflake.Node, w http.ResponseWriter, r *http.Request) error {
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

	id := sf.Generate()
	output, err := hnd.queries.AddEmailToMailingList(
		ctx,
		database.AddEmailToMailingListParams{
			ID:    int64(id),
			Email: input.Email,
		},
	)
	if err != nil {
		if sqlErr, ok := err.(sqlite3.Error); ok {
			if sqlErr.Code == sqlite3.ErrConstraint {
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
