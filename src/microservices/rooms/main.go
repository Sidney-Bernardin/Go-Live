package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"rooms/apis/http"
	"rooms/configuration"
	"rooms/domain"
	"rooms/repositories/cache/redis"
	"rooms/repositories/users_client/grpc"
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
	config, err := configuration.New("rooms")
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot load configuration")
	}

	// Create a mongo database repository.
	cacheRepo, err := redis.NewCacheRepository(config)
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot create redis cache database repository")
	}

	// Create a new users client repository.
	usersClientRepo, err := grpc.NewUsersClientRepository(config)
	if err != nil {
		logger.Fatal().Stack().Err(err).Msg("Cannot create grpc users client repository")
	}

	// Create a service.
	svc := domain.NewService(config, cacheRepo, usersClientRepo)

	// Setup a new wait-group and signal-channel.
	wg := sync.WaitGroup{}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Serve the api.
	wg.Add(1)
	go serveAPI(&wg, "HTTP", http.New(svc, config, &logger), config.HTTPPort, &logger)

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
