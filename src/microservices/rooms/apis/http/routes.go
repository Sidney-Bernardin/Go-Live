package http

func (a *api) doRoutes() {

	a.router.HandleFunc(
		"/",
		a.adapt(
			a.handleIndex,
			a.logRequest(),
		),
	).Methods("GET")

	a.router.HandleFunc(
		"/create",
		a.adapt(
			a.handleCreateRoom,
			a.logRequest(),
		),
	).Methods("POST")

	a.router.HandleFunc(
		"/delete",
		a.adapt(
			a.handleDeleteRoom,
			a.logRequest(),
		),
	).Methods("POST")

	a.router.HandleFunc(
		"/authenticate_stream",
		a.adapt(
			a.handleAuthenticateStream,
			a.logRequest(),
		),
	).Methods("GET")

	a.router.HandleFunc(
		"/join",
		a.adapt(
			a.handleJoinRoom,
			a.logRequest(),
		),
	).Methods("GET")

	a.router.HandleFunc(
		"/all/{room_id}",
		a.adapt(
			a.handleGetRoom,
			a.logRequest(),
		),
	).Methods("GET")
}
