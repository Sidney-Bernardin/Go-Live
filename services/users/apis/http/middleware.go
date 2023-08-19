package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"users/domain"

	"github.com/gorilla/mux"
)

const (
	mwBearerToken = iota
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

func (a *api) getBearerToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the token from the request's Authorization header.
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer")

		// Add the token to the request's context and call the next handler.
		ctx := context.WithValue(r.Context(), mwBearerToken, strings.TrimSpace(token))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *api) getFormFile(name string, maxSize int64) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Parse the request's form data.
			if err := r.ParseMultipartForm(maxSize); err != nil {
				a.err(w, domain.ProblemDetail{
					Problem: domain.ProblemInvalidInput,
					Detail:  fmt.Sprintf("Cannot parse form data: %v", err)})
				return
			}

			// Get the file from the request's form data.
			file, _, err := r.FormFile(name)
			if err != nil {

				// Check if the file is missing.
				if err == http.ErrMissingFile {
					a.err(w, domain.ProblemDetail{
						Problem: domain.ProblemInvalidInput,
						Detail:  fmt.Sprintf("Missing the %s file from request form data.", name)})
					return
				}
			}
			defer file.Close()

			// Get the file's bytes.
			fileBytes, err := io.ReadAll(file)
			if err != nil {
			}

			// Add the file to the request's context and call the next handler.
			ctx := context.WithValue(r.Context(), name, fileBytes)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
