package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"users/configuration"
	"users/domain"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const svrErrMsg = "Server Error"

type api struct {
	server  *http.Server
	service domain.Service

	router      *mux.Router
	formDocoder *schema.Decoder

	logger *zerolog.Logger
}

func NewAPI(config *configuration.Config, l *zerolog.Logger, svc domain.Service) *api {

	// Create an api.
	a := &api{
		server: &http.Server{
			Addr:         fmt.Sprintf("0.0.0.0:%v", config.HTTPPort),
			ReadTimeout:  config.HTTPReadTimeout,
			WriteTimeout: config.HTTPWriteTimeout,
		},
		service:     svc,
		router:      mux.NewRouter(),
		formDocoder: schema.NewDecoder(),
		logger:      l,
	}

	a.server.Handler = a
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
