package http

import "github.com/go-chi/chi/v5"

func (api *API) routes() {
	r := chi.NewRouter()
	api.server.Handler = r

	r.Use(api.mwLog, api.mwInitDetails)

	r.Route("/login", func(r chi.Router) {
		r.Post("/signup", api.handleSignup)
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/", api.handleGetUser)
	})
}
