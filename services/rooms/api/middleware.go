package api

import (
	"net/http"
)

func (a *api) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Log the request.
		a.logger.Info().
			Str("method", r.Method).
			Str("uri", r.RequestURI).
			Msg("New request")

		next.ServeHTTP(w, r)
	})
}
