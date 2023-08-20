package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"users/apis/grpc"
	"users/apis/http"
	"users/configuration"
	"users/domain"
	"users/repositories/database/mongo"
)

func main() {

	// Create a logger.
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	logger := zerolog.New(os.Stdout)
	logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Create a Config.
	config, err := configuration.NewConfig("users")
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot create configuration")
	}

	var svc domain.Service
	if !config.Mock {

		// Create a mongo database repository.
		databaseRepo, err := mongo.NewDatabaseRepository(config)
		if err != nil {
			logger.Fatal().Stack().Err(err).Msg("Cannot create mongo database repository")
		}

		// Create a Service.
		svc = domain.NewService(config, databaseRepo)
	}

	var (
		apiErrChan = make(chan error)
		signalChan = make(chan os.Signal)

		// Create APIs.
		httpAPI = http.NewAPI(config, &logger, svc)
		grpcAPI = grpc.NewAPI(config, &logger, svc)
	)

	// Serve APIs.
	go func() { apiErrChan <- httpAPI.Serve() }()
	go func() { apiErrChan <- grpcAPI.Serve() }()

	logger.Info().
		Int("http_port", config.HTTPPort).
		Int("grpc_port", config.GRPCPort).
		Msg("Serving APIs")

	// Send interrupt and termination signals to the signal-channel.
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Wait for a signal or non-nil error from the channels.
	for {
		select {
		case <-signalChan:
		case err := <-apiErrChan:
			if err == nil {
				continue
			}
			logger.Error().Stack().Err(err).Msg("An API crashed")
		}

		break
	}

	logger.Info().Msg("Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer cancel()

	// Shutdown APIs.
	grpcAPI.Shutdown()
	if err := httpAPI.Shutdown(ctx); err != nil {
		logger.Error().Stack().Err(err).Msg("Error during shutdown")
	}
}
