package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"users/domain"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
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

func (a *api) getPathParams(t reflect.Kind, names ...string) adapter {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			for _, name := range names {

				// Get the parameter's value from the request.
				var value any = mux.Vars(r)[name]
				var err error

				if t == reflect.Int {

					// Convert the value to an int.
					value, err = strconv.Atoi(value.(string))
					if err != nil {

						// Respond with StatusUnprocessableEntity.
						a.err(w, http.StatusUnprocessableEntity, domain.ProblemDetail{
							Type:   domain.PDTypeInvalidInput,
							Detail: name + " doesn't resemble an int.",
						})

						return
					}
				}

				// Add the parameter to the request.
				r = a.addDataToRequest(r, map[string]any{name: value})
			}

			// Call the next handler.
			next.ServeHTTP(w, r)
		}
	}
}

func (a *api) getQueryParams(t reflect.Kind, names ...string) adapter {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			for _, name := range names {

				// Get the parameter's value from the request.
				var value any = r.URL.Query().Get(name)
				var err error

				if t == reflect.Int {

					// Convert the value to an int.
					value, err = strconv.Atoi(value.(string))
					if err != nil {

						// Respond with StatusUnprocessableEntity.
						a.err(w, http.StatusUnprocessableEntity, domain.ProblemDetail{
							Type:   domain.PDTypeInvalidInput,
							Detail: name + " doesn't resemble an int.",
						})

						return
					}
				}

				// Add the parameter to the request.
				r = a.addDataToRequest(r, map[string]any{name: value})
			}

			// Call the next handler.
			next.ServeHTTP(w, r)
		}
	}
}

func (a *api) getFormFile(name string, maxSize int64) adapter {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			// Parse the request's form data.
			if err := r.ParseMultipartForm(maxSize); err != nil {

				// Respond with StatusUnprocessableEntity.
				a.err(w, http.StatusUnprocessableEntity, domain.ProblemDetail{
					Type:   domain.PDTypeInvalidInput,
					Detail: fmt.Sprintf("Cannot parse form data: '%v'", err),
				})

				return
			}

			// Get the file from the request.
			file, _, err := r.FormFile(name)
			if err != nil {

				// If the file is missing, respond with StatusUnprocessableEntity.
				if err == http.ErrMissingFile {
					a.err(w, http.StatusUnprocessableEntity, domain.ProblemDetail{
						Type:   domain.PDTypeInvalidInput,
						Detail: fmt.Sprintf("Missing the '%s' file from the request.", name),
					})
					return
				}

				// Respond with StatusInternalServerError.
				a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot get file from request"))
				return
			}
			defer file.Close()

			// Read the file.
			bytes, err := ioutil.ReadAll(file)
			if err != nil {

				// Respond with StatusInternalServerError.
				a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot read file"))
				return
			}

			// Add the file to the request and call the next handler.
			r = a.addDataToRequest(r, map[string]any{name: bytes})
			next.ServeHTTP(w, r)
		}
	}
}
