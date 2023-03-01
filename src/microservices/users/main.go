package main

import (
	"os"
	"os/signal"
	"sync"
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

type API interface {
	Serve(port int) error
}

func main() {

	// Create a logger.
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	logger := zerolog.New(os.Stdout)
	logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Create a configuration.
	config, err := configuration.New("users")
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot load configuration")
	}

	// Create a mongo database repository.
	databaseRepo, err := mongo.NewDatabaseRepository(config)
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot create mongo database repository")
	}

	// Create a service.
	svc := domain.NewService(config, databaseRepo)

	// Setup a new wait-group and signal-channel.
	wg := sync.WaitGroup{}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Serve the apis.
	wg.Add(2)
	go serveAPI(&wg, "HTTP", http.New(svc, &logger), config.HTTPPort, &logger)
	go serveAPI(&wg, "GRPC", grpc.New(svc, &logger), config.GRPCPort, &logger)

	// Wait for the wait-group in another go-routine.
	go func() {
		wg.Wait()
		close(signalChan)
	}()

	// Stop until a signal is received or the channel closes.
	<-signalChan
}

// serveAPI decrements the wait-group after the API is done serving.
func serveAPI(wg *sync.WaitGroup, name string, api API, port int, l *zerolog.Logger) {
	defer wg.Done()

	l.Info().Int("port", port).Msgf("Serving %s", name)

	// Serve the API.
	if err := api.Serve(port); err != nil {
		l.Error().Stack().Err(err).Msgf("Cannot serve %s", name)
		return
	}

	l.Info().Msgf("Done serving %s", name)
}
