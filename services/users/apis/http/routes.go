package http

func (a *api) doRoutes() {

	a.router.HandleFunc(
		"/",
		a.adapt(
			a.handleIndex,
			a.logRequest(),
		),
	)

	a.router.HandleFunc(
		"/signup",
		a.adapt(
			a.handleSignup,
			a.logRequest(),
		),
	).Methods("POST")

	a.router.HandleFunc(
		"/signin",
		a.adapt(
			a.handleSignin,
			a.logRequest(),
		),
	).Methods("POST")

	a.router.HandleFunc(
		"/self",
		a.adapt(
			a.handleGetSelf,
			a.logRequest(),
			a.getBearerToken(),
		),
	).Methods("GET")

	a.router.HandleFunc(
		"/self/logout",
		a.adapt(
			a.handleLogout,
			a.logRequest(),
			a.getBearerToken(),
		),
	).Methods("GET")

	a.router.HandleFunc(
		"/self/picture",
		a.adapt(
			a.handleSetProfilePicture,
			a.logRequest(),
			a.getBearerToken(),
			a.getFormFile("profile_picture", 10<<20),
		),
	).Methods("POST")

	a.router.HandleFunc(
		"/all",
		a.adapt(
			a.handleSearchUsers,
			a.logRequest(),
		),
	).Methods("GET")

	a.router.HandleFunc(
		"/all/{search}",
		a.adapt(
			a.handleGetUser,
			a.logRequest(),
		),
	).Methods("GET")

	a.router.HandleFunc(
		"/all/{user_id}/picture",
		a.adapt(
			a.handleGetProfilePicture,
			a.logRequest(),
		),
	).Methods("GET")
}
