package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"users/src"
	"users/src/domain"
	"users/src/domain/service"

	"github.com/pkg/errors"
)

type API struct {
	config *src.Config
	logger *slog.Logger
	svc    *service.Service

	server *http.Server
}

func NewAPI(logger *slog.Logger, config *src.Config, svc *service.Service) *API {

	server := &http.Server{
		Addr: config.HTTPAddr,
	}

	api := &API{config, logger, svc, server}
	api.routes()

	return api
}

func (api *API) write(w http.ResponseWriter, r *http.Request, statusCode int, data any) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		api.logger.Error("Internal Server Error", "err", err.Error())
	}
}

var domainErrorCodes = map[domain.DomainErrorType]int{
	domain.DomainErrorTypeUserDoesNotExist:    http.StatusNotFound,
	domain.DomainErrorTypeSessionDoesNotExist: http.StatusNotFound,

	domain.DomainErrorTypeUUIDInvalid:     http.StatusBadRequest,
	domain.DomainErrorTypeUsernameInvalid: http.StatusBadRequest,
	domain.DomainErrorTypePasswordInvalid: http.StatusBadRequest,
}

func (api *API) err(w http.ResponseWriter, r *http.Request, err error) {

	type domainErrorView struct {
		Type    string         `json:"type"`
		Message string         `json:"message,omitempty"`
		Details map[string]any `json:"details,omitempty"`
	}

	var domainErr *domain.DomainError
	if !errors.As(err, &domainErr) {
		api.write(w, r, http.StatusInternalServerError, err)
		api.logger.LogAttrs(r.Context(), slog.LevelError, "Internal Server Error",
			slog.Group("error", "msg", err.Error()),
			api.requestAttr(r))
		return
	}

	api.write(w, r, domainErrorCodes[domainErr.Type], &domainErrorView{
		Type:    string(domainErr.Type),
		Message: domainErr.Message,
		Details: domainErr.Details,
	})
}

func (api *API) requestAttr(r *http.Request) slog.Attr {
	return slog.Group("request",
		"method", r.Method,
		"path", r.URL.Path,
	)
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.server.Handler.ServeHTTP(w, r)
}
