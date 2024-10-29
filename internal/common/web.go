package common

import (
	"net/http"

	"github.com/a-h/templ"
)

func Render(w http.ResponseWriter, r *http.Request, c templ.Component) error {
	w.Header().Set("Content-Type", "text/html")

	err := c.Render(r.Context(), w)
	if err != nil {
		return err
	}

	return nil
}
