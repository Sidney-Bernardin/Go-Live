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

type Api struct {
	config *src.Config
	logger *slog.Logger
	svc    *service.Service

	server *http.Server
}

func New(config *src.Config, logger *slog.Logger, svc *service.Service) *Api {

	server := &http.Server{
		Addr: config.HttpAddr,
	}

	api := &Api{config, logger, svc, server}
	api.routes()

	return api
}

func (api *Api) write(w http.ResponseWriter, r *http.Request, statusCode int, data any) {
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

func (api *Api) err(w http.ResponseWriter, r *http.Request, err error) {

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

func (api *Api) requestAttr(r *http.Request) slog.Attr {
	return slog.Group("request",
		"method", r.Method,
		"path", r.URL.Path,
	)
}

func (api *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.server.Handler.ServeHTTP(w, r)
}
