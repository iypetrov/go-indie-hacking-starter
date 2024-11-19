package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/iypetrov/go-indie-hacking-starter/config"
	"github.com/iypetrov/go-indie-hacking-starter/internal/router"
	"github.com/iypetrov/go-indie-hacking-starter/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Init()
	config.Init()
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Get().App.Port),
		Handler:      router.New(ctx, static()),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Get().Info("server started on %s", config.Get().App.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.Get().Error("cannot start server: %s", err.Error())
	}
}
