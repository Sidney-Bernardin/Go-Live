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

func New(config *src.Config, logger *slog.Logger, svc *service.Service) *API {
	server := &http.Server{
		Addr: config.HttpAddr,
	}

	api := &API{config, logger, svc, server}
	api.routes()

	return api
}

func Run() {

}

var domainErrorCodes = map[domain.DomainErrorType]int{
	domain.DomainErrorTypeUserDoesNotExist:    http.StatusNotFound,
	domain.DomainErrorTypeSessionDoesNotExist: http.StatusNotFound,

	domain.DomainErrorTypeUUIDInvalid:     http.StatusBadRequest,
	domain.DomainErrorTypeUsernameInvalid: http.StatusBadRequest,
	domain.DomainErrorTypePasswordInvalid: http.StatusBadRequest,
}

func (api *API) err(w http.ResponseWriter, r *http.Request, statusCode int, err error) {

	var view struct {
		Type    string         `json:"type,omitempty"`
		Msg     string         `json:"message,omitempty"`
		Details map[string]any `json:"details,omitempty"`
	}

	if domainErr := new(domain.DomainError); errors.As(err, &domainErr) {
		statusCode = domainErrorCodes[domainErr.Type]
		view.Type = string(domainErr.Type)
		view.Msg = domainErr.Msg
		view.Details = domainErr.Details
	} else if statusCode >= 500 {
		view.Msg = http.StatusText(statusCode)
		api.logger.LogAttrs(r.Context(), slog.LevelError, view.Msg,
			slog.Group("error", "msg", err.Error()),
			api.requestAttr(r))
	} else {
		view.Msg = err.Error()
	}

	api.write(w, statusCode, view)
}

func (api *API) write(w http.ResponseWriter, statusCode int, data any) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		api.logger.Error("Internal Server Error", "err", err.Error())
	}
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
