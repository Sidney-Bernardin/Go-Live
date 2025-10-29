package http

import (
	"net/http"
	"users/src"
	"users/src/domain"

	"github.com/pkg/errors"
)

type userView struct {
	ID string `json:"id"`

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

func (api *API) handleGetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch {

	case r.URL.Query().Has("id"):

		userID, err := domain.NewUUIDFromString(r.URL.Query().Get("id"))
		if err != nil {
			api.err(w, errors.Wrap(err, "failed creating UUID"))
			return
		}

		user, err := api.svc.GetUserByID(ctx, userID)
		if err != nil {
			api.err(w, errors.Wrap(err, "failed getting user"))
			return
		}

		api.write(w, http.StatusOK, newUserView(user))

	default:

		username, err := domain.NewUsername(r.URL.Query().Get("username"))
		if err != nil {
			api.err(w, errors.Wrap(err, "failed creating username"))
			return
		}

		users, err := api.svc.SearchUsers(ctx, username)
		if err != nil {
			api.err(w, errors.Wrap(err, "failed getting users"))
			return
		}

		api.write(w, http.StatusOK, src.MapSlice(users, func(user *domain.User) *userView {
			return newUserView(user)
		}))
	}
}
