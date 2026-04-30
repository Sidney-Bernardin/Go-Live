package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"users/src"
	"users/src/apis/http"
	"users/src/domain/service"
	"users/src/repos/blobstore"
	"users/src/repos/cache"
	"users/src/repos/database"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	// Create configuration.
	config, err := src.NewConfig()
	if err != nil {
		slog.Error("Failed creating configuration", src.ErrGroup(err))
		return
	}

	// Create logger.
	logger := src.NewLogger(config)

	// Create database repository.
	databaseRepo, err := database.New(ctx, config)
	if err != nil {
		logger.Error("Failed creating database repository", src.ErrGroup(err))
		return
	}

	// Create cache repository.
	cacheRepo, err := cache.New(ctx, config)
	if err != nil {
		logger.Error("Failed creating cache repository", src.ErrGroup(err))
		return
	}

	// Create blob-store repository.
	blobStoreRepo, err := blobstore.New(ctx, config)
	if err != nil {
		logger.Error("Failed creating blob-store repository", src.ErrGroup(err))
		return
	}

	// Create service.
	svc := service.New(config, logger.With("resource", "domain-service"), databaseRepo, cacheRepo, blobStoreRepo)

	// Create and run HTTP API.
	http := http.New(config, logger.With("resource", "http-api"), svc)
	go func() { http.Run(); cancel() }()

	<-ctx.Done()
	cancel()

	// Shutdown APIs.
	http.Shutdown()
}
