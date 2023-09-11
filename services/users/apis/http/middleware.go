package http

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

const (
	mwBearerToken = iota
	mwFormValues
	mwFormFiles
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

func (a *api) getFormData(formValues any, formFiles ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Parse the request's form-data.
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				a.err(w, errors.Wrap(err, "cannot parse form data"))
				return
			}

			// Decode the request's baisc form-data.
			if err := a.formDocoder.Decode(formValues, r.MultipartForm.Value); err != nil {
				if err, ok := err.(schema.MultiError); !ok {
					a.err(w, errors.Wrap(err, "cannot decode form data"))
					return
				}
			}

			// Decode the request's file form-data.
			files := map[string][]byte{}
			for _, name := range formFiles {

				h, ok := r.MultipartForm.File[name]
				if !ok {
					continue
				}

				// Open the file.
				file, err := h[0].Open()
				if err != nil {
					a.err(w, errors.Wrap(err, "cannot open file"))
					return
				}
				defer file.Close()

				// Read the file and add it's bytes to the files-map.
				if files[name], err = io.ReadAll(file); err != nil {
					a.err(w, errors.Wrap(err, "cannot read file"))
					return
				}
			}

			// Add the all of the form-data to the request's context and call the next handler.
			ctx := context.WithValue(r.Context(), mwFormValues, formValues)
			ctx = context.WithValue(ctx, mwFormFiles, files)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
