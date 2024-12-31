package main

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
)

//go:embed static
var staticFS embed.FS

func StaticFiles(logger Logger) http.Handler {
	if Profile == "local" {
		logger.Info("serving static files from local directory")
		return http.StripPrefix("/static", http.FileServerFS(os.DirFS("static")))
	}

	fs, err := fs.Sub(staticFS, "static")
	if err != nil {
		logger.Info("serving static files from local directory")
		return http.StripPrefix("/static", http.FileServerFS(os.DirFS("static")))
	}

	logger.Info("serving static files from embedded FS")
	// return http.FileServer(http.FS(fs))
	return http.StripPrefix("/static", http.FileServer(http.FS(fs)))
}
