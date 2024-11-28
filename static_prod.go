//go:build prod
// +build prod

package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static
var staticFS embed.FS

func static(logger Logger) http.Handler {
	fs, err := fs.Sub(staticFS, "static")
	if err != nil {
		logger.Info("serving static files from local directory")
		return http.StripPrefix("/", http.FileServer(http.Dir("static")))
	}
	logger.Info("serving static files from embedded FS")
	return http.FileServer(http.FS(fs))
}
