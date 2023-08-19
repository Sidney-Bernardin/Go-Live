package http

func (a *api) doRoutes() {
	a.router.Use(a.logRequest)
	a.router.HandleFunc("/", a.handleIndex)

	a.router.HandleFunc("/signup", a.handleSignup).Methods("POST")
	a.router.HandleFunc("/signin", a.handleSignin).Methods("POST")

	a.router.HandleFunc("/all", a.handleSearchUsers).Methods("GET")
	a.router.HandleFunc("/all/{user_id}", a.handleGetUser).Methods("GET")
	a.router.HandleFunc("/all/{user_id}/picture", a.handleGetProfilePicture).Methods("GET")

	auth := a.router.PathPrefix("/auth").Subrouter()
	auth.Use(a.getBearerToken)
	auth.HandleFunc("/", a.handleAuthenticateUser).Methods("GET")
	auth.HandleFunc("/logout", a.handleLogout).Methods("DELETE")

	postPicture := auth.PathPrefix("/picture").Subrouter()
	postPicture.Use(a.getFormFile("profile_picture", 10<<20))
	postPicture.HandleFunc("/", a.handleSetProfilePicture).Methods("PUT")
}
