package http

import (
	"net/http"
	"strings"

	"rooms/domain"
)

type adapter func(http.HandlerFunc) http.HandlerFunc

func (a *api) adapt(h http.HandlerFunc, adapters ...adapter) http.HandlerFunc {
	for _, a := range adapters {
		h = a(h)
	}
	return h
}

func (a *api) logRequest() adapter {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			// Log the request.
			a.logger.Info().
				Str("method", r.Method).
				Str("uri", r.RequestURI).
				Msg("New request")

			next.ServeHTTP(w, r)
		}
	}
}

func (a *api) getAuthToken() adapter {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			// Get the Authorization header from the request.
			parts := strings.Split(r.Header.Get("Authorization"), "Bearer")
			if len(parts) != 2 {

				// Respond with StatusUnauthorized.
				a.err(w, http.StatusUnauthorized, domain.ProblemDetail{
					Type: domain.PDTypeUnauthorized,
				})

				return
			}

			// Get the second part without white space.
			token := strings.TrimSpace(parts[1])

			// Add token to the request and call the next handler.
			r = a.addDataToRequest(r, map[string]any{"authorization_token": token})
			next.ServeHTTP(w, r)
		}
	}
}
