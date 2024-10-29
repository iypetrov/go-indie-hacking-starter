//go:build local
// +build local

package main

import (
	"net/http"
)

func static() http.Handler {
	return http.StripPrefix("/", http.FileServer(http.Dir("static")))
}
