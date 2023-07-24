package http

import (
	"context"
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

const middlewareCtxKey = 0

type middlewareData struct {
	bearerToken string
	formFile    []byte
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

func (a *api) getBearerToken() adapter {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			// Get the Authorization header from the request, and split it into
			// two parts.
			parts := strings.Split(r.Header.Get("Authorization"), "Bearer")
			if len(parts) != 2 {
				a.err(w, http.StatusUnauthorized, domain.ProblemDetail{
					Type: domain.PDTypeUnauthorized,
				})
				return
			}

			// Get the requests middleware-data.
			mwData := &middlewareData{}
			mwData = r.Context().Value(middlewareCtxKey).(*middlewareData)

			// Add the second part to the middleware-data with it's white space trimed.
			mwData.bearerToken = strings.TrimSpace(parts[1])

			// Call the next handler with the updated middleware-data.
			ctx := context.WithValue(r.Context(), middlewareCtxKey, &mwData)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
