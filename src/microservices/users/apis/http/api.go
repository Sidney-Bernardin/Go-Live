package http

import (
	"fmt"
	"net/http"
	"users/domain"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type api struct {
	service domain.Service
	router  *mux.Router
	logger  *zerolog.Logger
}

// New returns a new api.
func New(svc domain.Service, l *zerolog.Logger) *api {

	a := &api{
		service: svc,
		router:  mux.NewRouter(),
		logger:  l,
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
