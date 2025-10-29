package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	"users/src"
	"users/src/domain"
	"users/src/domain/service"

	"github.com/pkg/errors"
)

type API struct {
	config  *src.Config
	logger  *slog.Logger
	service *service.Service

	server *http.Server
}

func New(config *src.Config, logger *slog.Logger, svc *service.Service) *API {
	svr := &http.Server{
		Addr: config.HttpAddr,
	}

	api := &API{config, logger, svc, svr}
	api.routes()

	return api
}

func (api *API) Run() {
	api.logger.Info("Running...", "addr", api.config.HttpAddr)
	err := api.server.ListenAndServe()
	api.logger.Error("Failed running server", src.ErrGroup(err))
}

func (api *API) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	api.logger.Info("Shutting down...", "addr", api.config.HttpAddr)
	err := api.server.Shutdown(ctx)
	api.logger.Error("Failed shutting down server", src.ErrGroup(err))
}

func (api *API) write(w http.ResponseWriter, statusCode int, data any) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		api.logger.Error("Failed writting response", src.ErrGroup(err))
	}
}

func (api *API) err(w http.ResponseWriter, e error) {
	switch e := errors.Cause(e).(type) {

	case *apiError[apiErrorType]:
		api.write(w, apiCodes[e.Type], e)

	case *domain.DomainError:
		api.write(w, domainCodes[e.Type], &apiError[domain.DomainErrorType]{
			Type:    e.Type,
			Message: e.Message,
			Details: e.Details,
		})

	default:
		api.logger.Error("Internal server error", src.ErrGroup(e))
		api.write(w, apiCodes[apiErrorTypeInternalServerError], &apiError[apiErrorType]{
			Type:    apiErrorTypeInternalServerError,
			Message: "Internal Server Error",
		})
	}
}
