package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"rooms/api"
	"rooms/configuration"
	"rooms/domain"
	"rooms/repositories/cache/redis"
	"rooms/repositories/users_client/grpc"
)

func main() {

	// Create a logger.
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	logger := zerolog.New(os.Stdout)
	logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Create a Config.
	config, err := configuration.NewConfig("rooms")
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot create configuration")
	}

	// Create a redis cache repository.
	cacheRepo := redis.NewCacheRepository(config)

	// Create a users-client repository.
	usersClientRepo, err := grpc.NewUsersClientRepository(config)
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot create grpc users-client repository")
	}

	// Create a service.
	svc := domain.NewService(config, cacheRepo, usersClientRepo)

	var (
		apiErrChan = make(chan error)
		signalChan = make(chan os.Signal)

		// Create an API.
		httpAPI = api.NewAPI(config, &logger, svc)
	)

	// Serve the API.
	go func() { apiErrChan <- httpAPI.Serve() }()
	logger.Info().Int("port", config.Port).Msg("Serving API")

	// Send interrupt and termination signals to signalChan.
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Wait for a signal or error from the channels.
	select {
	case <-signalChan:
	case err := <-apiErrChan:
		logger.Error().Stack().Err(err).Msg("API crashed")
	}

	logger.Info().Msg("Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer cancel()

	// Shutdown the API.
	if err := httpAPI.Shutdown(ctx); err != nil {
		logger.Error().Stack().Err(err).Msg("Error during shutdown")
	}
}
