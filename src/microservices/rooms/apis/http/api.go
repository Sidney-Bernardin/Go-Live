package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"rooms/configuration"
	"rooms/domain"
)

const (
	PDTypeCannotProcessRequestData = "cannot_process_request_data"
	PDTypeCannotUpgradeRequest     = "cannot_upgrade_request"
	PDTypeInvalidCallback          = "invalid_callback"
)

type api struct {
	service  domain.Service
	router   *mux.Router
	upgrader *websocket.Upgrader
	logger   *zerolog.Logger

	pongTimeout time.Duration
}

func New(svc domain.Service, config *configuration.Configuration, l *zerolog.Logger) *api {

	a := &api{
		service: svc,
		router:  mux.NewRouter(),
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		logger:      l,
		pongTimeout: config.HTTPPongTimeout,
	}

	a.doRoutes()
	return a
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func (a *api) Serve(port int) error {
	addr := fmt.Sprintf(":%v", port)
	err := http.ListenAndServe(addr, a)
	return errors.Wrap(err, "cannot listen and serve")
}
