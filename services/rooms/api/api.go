package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"rooms/configuration"
	"rooms/domain"
)

type api struct {
	server   *http.Server
	router   *mux.Router
	upgrader *websocket.Upgrader

	service domain.Service
	logger  *zerolog.Logger

	wsCloseTimeout time.Duration
}

func NewAPI(config *configuration.Config, l *zerolog.Logger, svc domain.Service) *api {

	// Create an api.
	a := &api{
		service:        svc,
		router:         mux.NewRouter(),
		logger:         l,
		upgrader:       &websocket.Upgrader{},
		wsCloseTimeout: config.WSCloseTimeout,
	}

	// Create a Server.
	svr := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%v", config.Port),
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}

	a.server = svr
	svr.Handler = a

	a.doRoutes()
	return a
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

// Start starts the api's server.
func (a *api) Serve() error {

	// Create a listener.
	ln, err := net.Listen("tcp", a.server.Addr)
	if err != nil {
		return errors.Wrap(err, "cannot create listener")
	}

	// Start the server.
	err = a.server.Serve(ln)
	return errors.Wrap(err, "cannot serve")
}

// Shutdown gracefully shuts down the api's server.
func (a *api) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
