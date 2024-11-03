package common

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/iypetrov/go-indie-hacking-starter/internal/toast"
)

func Render(w http.ResponseWriter, r *http.Request, c templ.Component) error {
	w.Header().Set("Content-Type", "text/html")

	err := c.Render(r.Context(), w)
	if err != nil {
		return toast.ErrorInternalServerError(fmt.Errorf("server failed to render this component"))
	}

	return nil
}
