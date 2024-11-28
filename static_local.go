//go:build local
// +build local

package main

import (
	"net/http"
)

func static(logger Logger) http.Handler {
	logger.Info("serving static files from local directory")
	return http.StripPrefix("/", http.FileServer(http.Dir("static")))
}
