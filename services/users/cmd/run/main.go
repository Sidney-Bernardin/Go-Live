package main

import (
	"context"
	"log/slog"
	"os"
	"users/src"
	"users/src/apis/http"
	"users/src/domain/service"
	"users/src/repos/cache"
	"users/src/repos/database"
)

func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	config, err := src.NewConfig()
	if err != nil {
		logger.Error("Failed creating configuration", "err", err.Error())
		return
	}

	databaseRepo, err := database.New(ctx, config)
	if err != nil {
		logger.Error("Failed creating database repository", "err", err.Error())
		return
	}

	cacheRepo, err := cache.New(ctx, config)
	if err != nil {
		logger.Error("Failed creating cache repository", "err", err.Error())
		return
	}

	svc := service.New(config, databaseRepo, cacheRepo)
	httpAPI := http.New(config, logger, svc)
}
