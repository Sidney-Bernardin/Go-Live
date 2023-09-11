package http

import "users/domain"

func (a *api) doRoutes() {

	a.router.Use(a.logRequest)
	a.router.HandleFunc("/", a.handleIndex)

	// Get Users
	a.router.HandleFunc("/{username}/{offset}/{limit}", a.handleSearchUsers).Methods("GET")
	a.router.HandleFunc("/{user_id}", a.handleGetUser).Methods("GET")
	a.router.HandleFunc("/{user_id}/picture", a.handleGetProfilePicture).Methods("GET")

	// Signup
	signup := a.router.PathPrefix("/signup").Subrouter()
	signup.Use(a.getFormData(&domain.SignupInfo{}, "profile_picture"))
	signup.HandleFunc("", a.handleSignup).Methods("POST")

	// Signin
	signin := a.router.PathPrefix("/signin").Subrouter()
	signin.Use(a.getFormData(&domain.SigninInfo{}))
	signin.HandleFunc("", a.handleSignin).Methods("POST")

	// Authorization
	auth := a.router.PathPrefix("/auth").Subrouter()
	auth.Use(a.getBearerToken)
	auth.HandleFunc("/", a.handleAuthenticateUser).Methods("GET")
	auth.HandleFunc("", a.handleDeleteAccount).Methods("DELETE")
	auth.HandleFunc("/logout", a.handleLogout).Methods("DELETE")
}
