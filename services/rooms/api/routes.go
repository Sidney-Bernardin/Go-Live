package api

func (a *api) doRoutes() {

	a.router.Use(a.logRequest)
	a.router.HandleFunc("/", a.handleIndex)

	a.router.HandleFunc("/create", a.handleCreateRoom).Methods("POST")
	a.router.HandleFunc("/delete", a.handleDeleteRoom).Methods("POST")

	a.router.HandleFunc("/{room_id}", a.handleGetRoom).Methods("GET")
	a.router.HandleFunc("/join/{room_id}", a.handleJoinRoom).Methods("GET")
}
