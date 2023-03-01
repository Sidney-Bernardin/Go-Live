package http

import "reflect"

func (a *api) doRoutes() {

	a.router.HandleFunc(
		"/",
		a.adapt(a.handleIndex, a.logRequest()),
	)

	a.router.HandleFunc(
		"/all",
		a.adapt(
			a.handleSearchUsers,
			a.logRequest(),
			a.getQueryParams(reflect.String, "username", "fields"),
			a.getQueryParams(reflect.Int, "offset", "limit"),
		),
	).Methods("GET")

	a.router.HandleFunc(
		"/all/{search}",
		a.adapt(
			a.handleGetUser,
			a.logRequest(),
			a.getPathParams(reflect.String, "search"),
			a.getQueryParams(reflect.String, "search_by", "fields"),
		),
	).Methods("GET")

	a.router.HandleFunc(
		"/all/{user_id}/picture",
		a.adapt(
			a.handleGetProfilePicture,
			a.logRequest(),
			a.getPathParams(reflect.String, "user_id"),
		),
	).Methods("GET")

	a.router.HandleFunc(
		"/self",
		a.adapt(
			a.handleGetSelf,
			a.logRequest(),
			a.getAuthToken(),
			a.getQueryParams(reflect.String, "fields"),
		),
	).Methods("GET")

	a.router.HandleFunc(
		"/self/signup",
		a.adapt(a.handleSignup, a.logRequest()),
	).Methods("POST")

	a.router.HandleFunc(
		"/self/signin",
		a.adapt(a.handleSignin, a.logRequest()),
	).Methods("POST")

	a.router.HandleFunc(
		"/self/logout",
		a.adapt(
			a.handleLogout,
			a.logRequest(),
			a.getAuthToken(),
		),
	).Methods("GET")

	a.router.HandleFunc(
		"/self/picture",
		a.adapt(
			a.handleSetProfilePicture,
			a.logRequest(),
			a.getAuthToken(),
			a.getFormFile("profile_picture", 10<<20),
		),
	).Methods("POST")
}
