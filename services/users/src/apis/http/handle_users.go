package http

import (
	"net/http"
	"users/src/domain"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type userView struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

func newUserView(u *domain.User) *userView {
	return &userView{
		ID:       u.ID.String(),
		Username: string(u.Username),
		Email:    u.Email,
	}
}

func (api *API) handleUserGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := domain.NewUUIDFromString(chi.URLParam(r, "user_id"))
	if err != nil {
		api.err(w, errors.Wrap(err, "failed creating user-ID"))
		return
	}

	user, err := api.service.GetUserByID(ctx, userID)
	if err != nil {
		api.err(w, errors.Wrap(err, "failed getting user"))
		return
	}

	api.write(w, http.StatusOK, newUserView(user))
}

func (api *API) handleUserGetSelf(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionIDCookie, err := r.Cookie(sessionIDCookieName)
	if err != nil {
		api.err(w, &apiError[apiErrorType]{
			Type: apiErrorTypeUnauthorized,
		})
		return
	}

	sessionID, err := domain.NewUUIDFromString(sessionIDCookie.Value)
	if err != nil {
		api.err(w, errors.Wrap(err, "failed creating session-ID"))
		return
	}

	session, err := api.service.Authenticate(ctx, sessionID)
	if err != nil {
		api.err(w, errors.Wrap(err, "failed authenticating"))
		return
	}

	user, err := api.service.GetUserByID(ctx, session.UserID)
	if err != nil {
		api.err(w, errors.Wrap(err, "failed getting user"))
		return
	}

	api.write(w, http.StatusOK, newUserView(user))
}
