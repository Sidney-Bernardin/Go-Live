package http

import "github.com/go-chi/chi/v5"

func (api *API) routes() {
	r := chi.NewRouter()
	api.server.Handler = r

	r.Use(api.mwLog)

	r.Route("/login", func(r chi.Router) {
		r.Post("/signup", api.handleSignup)
		r.Post("/signin", api.handleSignin)
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/{user_id}", api.handleUserGet)
		r.Get("/self", api.handleUserGetSelf)
	})
}
