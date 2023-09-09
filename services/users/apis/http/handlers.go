package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"users/domain"
)

func (a *api) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("index"))
}

func (a *api) handleSignup(w http.ResponseWriter, r *http.Request) {

	var (
		info      = r.Context().Value(mwFormValues).(*domain.SignupInfo)
		formFiles = r.Context().Value(mwFormFiles).(map[string][]byte)
	)

	// Sign-Up a new User.
	sessionID, err := a.service.Signup(r.Context(), info, formFiles["profile_picture"])
	if err != nil {
		a.err(w, errors.Wrap(err, "cannot signup user"))
		return
	}

	// Respond.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res := &domain.LoginResponse{SessionID: sessionID}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		a.err(w, errors.Wrap(err, "cannot write response"))
	}
}

func (a *api) handleSignin(w http.ResponseWriter, r *http.Request) {

	info := r.Context().Value(mwFormValues).(*domain.SigninInfo)

	sessionID, err := a.service.Signin(r.Context(), info)
	if err != nil {
		a.err(w, errors.Wrap(err, "cannot signin user"))
		return
	}

	// Respond.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res := &domain.LoginResponse{SessionID: sessionID}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		a.err(w, errors.Wrap(err, "cannot write response"))
	}
}

func (a *api) handleLogout(w http.ResponseWriter, r *http.Request) {

	// Get the Session-ID from request's context.
	sessionID := r.Context().Value(mwBearerToken).(string)

	// Logout the Session-ID's User.
	if err := a.service.Logout(r.Context(), sessionID); err != nil {
		a.err(w, errors.Wrap(err, "cannot logout user"))
		return
	}

	// Respond.
	w.WriteHeader(http.StatusOK)
}

func (a *api) handleGetUser(w http.ResponseWriter, r *http.Request) {

	var (
		userID = mux.Vars(r)["user_id"]
		fields = r.URL.Query().Get("fields")
	)

	// Get the User-ID's User.
	user, err := a.service.GetUser(r.Context(), userID, strings.Split(fields, ",")...)
	if err != nil {
		a.err(w, errors.Wrap(err, "cannot get user"))
		return
	}

	// Respond.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		a.err(w, errors.Wrap(err, "cannot write response"))
	}
}

func (a *api) handleSearchUsers(w http.ResponseWriter, r *http.Request) {

	var (
		username = mux.Vars(r)["username"]
		offset   = mux.Vars(r)["offset"]
		limit    = mux.Vars(r)["limit"]
		fields   = r.URL.Query().Get("fields")
	)

	// Convert offset and limit to integers.
	offsetInt, errA := strconv.Atoi(offset)
	limitInt, errB := strconv.Atoi(limit)
	if errA != nil || errB != nil {
		a.err(w, domain.ProblemDetail{
			Problem: domain.ProblemInvalidInput,
			Detail:  "The offset and limit must resemble an integer."})
		return
	}

	// Get Users with usernames similar to the username.
	users, err := a.service.SearchUsers(r.Context(), username, offsetInt, limitInt, strings.Split(fields, ",")...)
	if err != nil {
		a.err(w, errors.Wrap(err, "cannot search users"))
		return
	}

	// Respond.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		a.err(w, errors.Wrap(err, "cannot write response"))
	}
}

func (a *api) handleAuthenticateUser(w http.ResponseWriter, r *http.Request) {

	var (
		sessionID = r.Context().Value(mwBearerToken).(string)
		fields    = r.URL.Query().Get("fields")
	)

	// Authenticate the Session-ID's User.
	user, err := a.service.AuthenticateUser(r.Context(), sessionID, strings.Split(fields, ",")...)
	if err != nil {
		a.err(w, errors.Wrap(err, "cannot authenticate user"))
		return
	}

	// Respond.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		a.err(w, errors.Wrap(err, "cannot write response"))
	}
}

func (a *api) handleGetProfilePicture(w http.ResponseWriter, r *http.Request) {

	// Get the profile-picture of the User-ID's User.
	profilePicture, err := a.service.GetProfilePicture(r.Context(), mux.Vars(r)["user_id"])
	if err != nil {
		a.err(w, errors.Wrap(err, "cannot get profile picture"))
		return
	}

	// Respond.
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, profilePicture); err != nil {
		a.err(w, errors.Wrap(err, "cannot write response"))
	}
}
