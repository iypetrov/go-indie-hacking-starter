//go:build prod
// +build prod

package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/iypetrov/go-indie-hacking-starter/logger"
)

//go:embed static
var staticFS embed.FS

func static() http.Handler {
	fs, err := fs.Sub(staticFS, "static")
	if err != nil {
		logger.Get().Error("failed to create sub filesystem: %s", err.Error())
		return http.StripPrefix("/", http.FileServer(http.Dir("static")))
	}
	return http.FileServer(http.FS(fs))
}
