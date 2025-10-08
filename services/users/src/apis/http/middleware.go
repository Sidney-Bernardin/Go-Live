package http

import (
	"context"
	"net/http"
	"users/src/domain"
)

func (api *Api) mwInitDetails(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*r = *r.WithContext(context.WithValue(r.Context(), domain.DomainErrorDetailsKey, map[string]any{}))
		next.ServeHTTP(w, r)
	})
}

func (api *Api) mwLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.logger.Info("New Request", api.requestAttr(r))
		next.ServeHTTP(w, r)
	})
}
