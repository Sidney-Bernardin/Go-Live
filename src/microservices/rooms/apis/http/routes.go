package http

import "reflect"

func (a *api) doRoutes() {

	a.router.HandleFunc(
		"/",
		a.adapt(a.handleIndex, a.logRequest()),
	).Methods("GET")

	a.router.HandleFunc(
		"/callbacks/{callback}",
		a.adapt(
			a.handleCallbacks,
			a.logRequest(),
			a.getPathParams(reflect.String, "callback"),
			a.getFormValues(reflect.String, "session_id", "name"),
		),
	).Methods("POST")

	a.router.HandleFunc(
		"/diagnostics",
		a.adapt(
			a.handleDiagnostics,
			a.logRequest(),
			a.getQueryParams(reflect.String, "session_id", "room_id"),
		),
	).Methods("GET")

	a.router.HandleFunc(
		"/all/{room_id}",
		a.adapt(
			a.handleGetRoom,
			a.logRequest(),
			a.getPathParams(reflect.String, "room_id"),
		),
	).Methods("GET")

}
