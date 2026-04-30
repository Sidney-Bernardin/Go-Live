package http

import (
	"log/slog"
	"net/http"
)

func (api *API) mwLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.logger.Info("New request", slog.Group("request",
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.Query().Encode(),
		))
		next.ServeHTTP(w, r)
	})
}
