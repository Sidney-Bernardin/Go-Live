package http

import (
	"encoding/json"
	"fmt"
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

	// Decode the request's body.
	var info *domain.SignupInfo
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		a.err(w, http.StatusUnprocessableEntity, domain.ProblemDetail{
			Type:   domain.PDTypeInvalidInput,
			Detail: fmt.Sprintf("Cannot decode request body: '%v'", err),
		})
		return
	}

	sessionID, err := a.service.Signup(info)
	if err != nil {

		// Check if the error was caused by a problem-detail.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot signup user"))
		return
	}

	// Respond.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res := &domain.LoginResponse{SessionID: sessionID}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot write response"))
		return
	}
}

func (a *api) handleSignin(w http.ResponseWriter, r *http.Request) {

	// Decode the request's body.
	var info *domain.SigninInfo
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		a.err(w, http.StatusUnprocessableEntity, domain.ProblemDetail{
			Type:   domain.PDTypeInvalidInput,
			Detail: fmt.Sprintf("Cannot decode request body: '%v'", err),
		})
		return
	}

	sessionID, err := a.service.Signin(info)
	if err != nil {

		// Check if the error was caused by a problem-detail.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot signin user"))
		return
	}

	// Respond.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res := &domain.LoginResponse{SessionID: sessionID}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot write response"))
		return
	}
}

func (a *api) handleLogout(w http.ResponseWriter, r *http.Request) {

	// Get the requests middleware-data.
	mwData := r.Context().Value(middlewareCtxKey).(*middlewareData)

	if err := a.service.Logout(mwData.bearerToken); err != nil {

		// Check if the error was caused by a problem-detail.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot logout user"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) handleGetUser(w http.ResponseWriter, r *http.Request) {

	search := mux.Vars(r)["search"]
	by := r.URL.Query().Get("by")
	fields := r.URL.Query().Get("fields")

	user, err := a.service.GetUser(search, by, strings.Split(fields, ",")...)
	if err != nil {

		// Check if the error was caused by a problem-detail.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot get user"))
		return
	}

	// Respond.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot write response"))
		return
	}
}

func (a *api) handleSearchUsers(w http.ResponseWriter, r *http.Request) {

	username := mux.Vars(r)["username"]
	offset := r.URL.Query().Get("offset")
	limit := r.URL.Query().Get("limit")
	fields := r.URL.Query().Get("fields")

	// Convert offset to an int.
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		a.err(w, http.StatusUnprocessableEntity, domain.ProblemDetail{
			Type:   domain.PDTypeInvalidInput,
			Detail: fmt.Sprintf("The given offset doesn't resemble an integer: '%v'", err),
		})
		return
	}

	// Convert limit to an int.
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		a.err(w, http.StatusUnprocessableEntity, domain.ProblemDetail{
			Type:   domain.PDTypeInvalidInput,
			Detail: fmt.Sprintf("The given limit doesn't resemble an integer: '%v'", err),
		})
		return
	}

	users, err := a.service.SearchUsers(username, offsetInt, limitInt, strings.Split(fields, ",")...)
	if err != nil {

		// Check if the error was caused by a problem-detail.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot get users"))
		return
	}

	// Respond.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot write response"))
		return
	}
}

func (a *api) handleGetSelf(w http.ResponseWriter, r *http.Request) {

	// Get the requests middleware-data.
	mwData := r.Context().Value(middlewareCtxKey).(*middlewareData)

	fields := r.URL.Query().Get("fields")

	user, err := a.service.GetSelf(mwData.bearerToken, strings.Split(fields, ",")...)
	if err != nil {

		// Check if the error was caused by a problem-detail.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot get user"))
		return
	}

	// Respond.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot write response"))
		return
	}
}

func (a *api) handleSetProfilePicture(w http.ResponseWriter, r *http.Request) {

	// Get the requests middleware-data.
	mwData := r.Context().Value(middlewareCtxKey).(*middlewareData)

	if err := a.service.SetProfilePicture(mwData.bearerToken, mwData.formFile); err != nil {

		// Check if the error was caused by a problem-detail.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot set profile picture"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) handleGetProfilePicture(w http.ResponseWriter, r *http.Request) {

	profilePicture, err := a.service.GetProfilePicture(mux.Vars(r)["user_id"])
	if err != nil {

		// Check if the error was caused by a problem-detail.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot get profile picture"))
		return
	}

	// Respond.
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, profilePicture); err != nil {
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot write response"))
		return
	}
}
