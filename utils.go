package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/bwmarrin/snowflake"
)

func Render(w http.ResponseWriter, r *http.Request, c templ.Component) error {
	w.Header().Set("Content-Type", "text/html")

	err := c.Render(r.Context(), w)
	if err != nil {
		return ErrorInternalServerError(fmt.Errorf("server failed to render this component"))
	}

	return nil
}

func MakeTemplHandler(ctx context.Context, logger Logger, sf *snowflake.Node, f func(ctx context.Context, logger Logger, sf *snowflake.Node, w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(ctx, logger, sf, w, r); err != nil {
			var t Toast
			if errors.As(err, &t) {
				AddToast(w, t)
			}
		}
	}
}
