package router

import (
	"errors"
	"net/http"

	"github.com/iypetrov/go-indie-hacking-starter/internal/toast"
)

// This is a wrapper fucntion that should wrap all our API handlers
// When they return an error it automatically will return to the client toast with the error 
func Make(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			var t toast.Toast
			if errors.As(err, &t) {
				toast.AddToast(w, t)
			}
		}
	}
}
