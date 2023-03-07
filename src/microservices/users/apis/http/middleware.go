package http

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"users/domain"

	"github.com/pkg/errors"
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
			mwData, ok := r.Context().Value(middlewareCtxKey).(*middlewareData)
			if !ok {
				mwData = &middlewareData{}
			}

			// Add the second part to the middleware-data with it's white space trimed.
			mwData.bearerToken = strings.TrimSpace(parts[1])

			// Call the next handler with the updated middleware-data.
			ctx := context.WithValue(r.Context(), middlewareCtxKey, mwData)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

func (a *api) getFormFile(name string, maxSize int64) adapter {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			// Parse the request's form data.
			if err := r.ParseMultipartForm(maxSize); err != nil {
				a.err(w, http.StatusUnprocessableEntity, domain.ProblemDetail{
					Type:   domain.PDTypeInvalidInput,
					Detail: fmt.Sprintf("Cannot parse form data: '%v'", err),
				})
				return
			}

			// Get the file from the request.
			file, _, err := r.FormFile(name)
			if err != nil {

				// Check if the file is missing.
				if err == http.ErrMissingFile {
					a.err(w, http.StatusUnprocessableEntity, domain.ProblemDetail{
						Type:   domain.PDTypeInvalidInput,
						Detail: fmt.Sprintf("Missing the '%s' file from the request.", name),
					})
					return
				}

				a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot get file from request"))
				return
			}
			defer file.Close()

			// Get the requests middleware-data.
			mwData, ok := r.Context().Value(middlewareCtxKey).(*middlewareData)
			if !ok {
				mwData = &middlewareData{}
			}

			// Read the file and add it to the middleware-data.
			mwData.formFile, err = ioutil.ReadAll(file)
			if err != nil {
				a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot read file"))
				return
			}

			// Call the next handler with the updated middleware-data.
			ctx := context.WithValue(r.Context(), middlewareCtxKey, mwData)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
