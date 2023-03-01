package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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

		// Respond with StatusUnprocessableEntity.
		a.err(w, http.StatusUnprocessableEntity, domain.ProblemDetail{
			Type:   domain.PDTypeInvalidInput,
			Detail: fmt.Sprintf("Cannot decode request body: '%v'", err),
		})

		return
	}

	sessionID, err := a.service.Signup(info)
	if err != nil {

		// If err was caused by a problem-detail, respond with StatusBadRequest.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		// Respond with StatusInternalServerError.
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot signup user"))
		return
	}

	// Respond.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res := &LoginResponse{SessionID: sessionID}
	if err := json.NewEncoder(w).Encode(res); err != nil {

		// Respond with StatusInternalServerError.
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot write response"))
		return
	}
}

func (a *api) handleSignin(w http.ResponseWriter, r *http.Request) {

	// Decode the request's body.
	var info *domain.SigninInfo
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {

		// Respond with StatusUnprocessableEntity.
		a.err(w, http.StatusUnprocessableEntity, domain.ProblemDetail{
			Type:   domain.PDTypeInvalidInput,
			Detail: fmt.Sprintf("Cannot decode request body: '%v'", err),
		})

		return
	}

	sessionID, err := a.service.Signin(info)
	if err != nil {

		// If err was caused by a problem-detail, respond with StatusBadRequest.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		// Respond with StatusInternalServerError.
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot signin user"))
		return
	}

	// Respond.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res := &LoginResponse{SessionID: sessionID}
	if err := json.NewEncoder(w).Encode(res); err != nil {

		// Respond with StatusInternalServerError.
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot write response"))
		return
	}
}

func (a *api) handleLogout(w http.ResponseWriter, r *http.Request) {

	// Get middleware data from the request's context value.
	mwData := r.Context().Value("data_from_middleware").(map[string]any)
	sessionID := mwData["authorization_token"].(string)

	if err := a.service.Logout(sessionID); err != nil {

		// If err was caused by a problem-detail, respond with StatusBadRequest.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		// Respond with StatusInternalServerError.
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot logout user"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) handleGetUser(w http.ResponseWriter, r *http.Request) {

	// Get middleware data from the request's context value.
	mwData := r.Context().Value("data_from_middleware").(map[string]any)
	search := mwData["search"].(string)
	searchBy := mwData["search_by"].(string)
	fields := mwData["fields"].(string)

	user, err := a.service.GetUser(search, searchBy, strings.Split(fields, ",")...)
	if err != nil {

		// If err was caused by a problem-detail, respond with StatusBadRequest.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		// Respond with StatusInternalServerError.
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot get user"))
		return
	}

	// Respond.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {

		// Respond with StatusInternalServerError.
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot write response"))
		return
	}
}

func (a *api) handleSearchUsers(w http.ResponseWriter, r *http.Request) {

	// Get middleware data from the request's context value.
	mwData := r.Context().Value("data_from_middleware").(map[string]any)
	offset := mwData["offset"].(int)
	limit := mwData["limit"].(int)
	fields := mwData["fields"].(string)
	username := mwData["username"].(string)

	users, err := a.service.SearchUsers(username, offset, limit, strings.Split(fields, ",")...)
	if err != nil {

		// If err was caused by a problem-detail, respond with StatusBadRequest.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		// Respond with StatusInternalServerError.
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot get users"))
		return
	}

	// Respond.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {

		// Respond with StatusInternalServerError.
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot write response"))
		return
	}
}

func (a *api) handleGetSelf(w http.ResponseWriter, r *http.Request) {

	// Get middleware data from the request's context value.
	mwData := r.Context().Value("data_from_middleware").(map[string]any)
	sessionID := mwData["authorization_token"].(string)
	fields := mwData["fields"].(string)

	user, err := a.service.GetSelf(sessionID, strings.Split(fields, ",")...)
	if err != nil {

		// If err was caused by a problem-detail, respond with StatusBadRequest.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		// Respond with StatusInternalServerError.
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot get user"))
		return
	}

	// Respond.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {

		// Respond with StatusInternalServerError.
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot write response"))
		return
	}
}

func (a *api) handleSetProfilePicture(w http.ResponseWriter, r *http.Request) {

	// Get middleware data from the request's context value.
	mwData := r.Context().Value("data_from_middleware").(map[string]any)
	sessionID := mwData["authorization_token"].(string)
	profilePicture := mwData["profile_picture"].([]byte)

	if err := a.service.SetProfilePicture(sessionID, profilePicture); err != nil {

		// If err was caused by a problem-detail, respond with StatusBadRequest.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		// Respond with StatusInternalServerError.
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot set profile picture"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) handleGetProfilePicture(w http.ResponseWriter, r *http.Request) {

	// Get middleware data from the request's context value.
	mwData := r.Context().Value("data_from_middleware").(map[string]any)
	userID := mwData["user_id"].(string)

	profilePicture, err := a.service.GetProfilePicture(userID)
	if err != nil {

		// If err was caused by a problem-detail, respond with StatusBadRequest.
		if pd, ok := errors.Cause(err).(domain.ProblemDetail); ok {
			a.err(w, http.StatusBadRequest, pd)
			return
		}

		// Respond with StatusInternalServerError.
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot get profile picture"))
		return
	}

	// Respond.
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, profilePicture); err != nil {

		// Respond with StatusInternalServerError.
		a.err(w, http.StatusInternalServerError, errors.Wrap(err, "cannot write response"))
		return
	}
}
