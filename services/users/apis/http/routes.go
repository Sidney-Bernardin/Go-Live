package http

func (a *api) doRoutes() {

	a.router.Use(a.logRequest)
	a.router.HandleFunc("/", a.handleIndex)

	a.router.HandleFunc("/signup", a.handleSignup).Methods("POST")
	a.router.HandleFunc("/signin", a.handleSignin).Methods("POST")

	a.router.HandleFunc("/{username}/{offset}/{limit}", a.handleSearchUsers).Methods("GET")
	a.router.HandleFunc("/{user_id}", a.handleGetUser).Methods("GET")
	a.router.HandleFunc("/{user_id}/picture", a.handleGetProfilePicture).Methods("GET")

	auth := a.router.PathPrefix("/auth").Subrouter()
	auth.Use(a.getBearerToken)
	auth.HandleFunc("/", a.handleAuthenticateUser).Methods("GET")
	auth.HandleFunc("/logout", a.handleLogout).Methods("DELETE")

	updatePicture := auth.PathPrefix("/picture").Subrouter()
	updatePicture.Use(a.getFormFile("profile_picture", 10<<20))
	updatePicture.HandleFunc("/", a.handleUpdateProfilePicture).Methods("PUT")
}
