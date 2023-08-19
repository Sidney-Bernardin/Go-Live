package api

func (a *api) doRoutes() {

	a.router.Use(a.logRequest)
	a.router.HandleFunc("/", a.handleIndex)

	a.router.HandleFunc("/create", a.handleCreateRoom).Methods("POST")
	a.router.HandleFunc("/delete", a.handleDeleteRoom).Methods("POST")

	a.router.HandleFunc("/join", a.handleJoinRoom).Methods("GET")
	a.router.HandleFunc("/all/{room_id}", a.handleGetRoom).Methods("GET")
}
