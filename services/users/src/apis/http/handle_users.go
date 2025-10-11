package http

import (
	"net/http"
	"users/src"
	"users/src/domain"

	"github.com/pkg/errors"
)

type UserView struct {
	ID string `json:"id"`

	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

func newUserView(u *domain.User) *UserView {
	return &UserView{
		ID:       u.ID.String(),
		Username: string(u.Username),
		Email:    u.Email,
	}
}

func (api *Api) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	details := ctx.Value(domain.DomainErrorDetailsKey).(map[string]any)

	switch {

	case r.URL.Query().Has("id"):

		userID, err := domain.NewUUIDFromString(ctx, r.URL.Query().Get("id"))
		if err != nil {
			api.err(w, r, errors.Wrap(err, "failed creating UUID"))
			return
		}
		details["user_id"] = userID

		user, err := api.svc.GetUserByID(ctx, userID)
		if err != nil {
			api.err(w, r, errors.Wrap(err, "failed getting user"))
			return
		}

		api.write(w, r, http.StatusOK, newUserView(user))

	default:

		username, err := domain.NewUsername(ctx, r.URL.Query().Get("username"))
		if err != nil {
			api.err(w, r, errors.Wrap(err, "failed creating username"))
			return
		}
		details["username"] = username

		users, err := api.svc.SearchUsers(ctx, username)
		if err != nil {
			api.err(w, r, errors.Wrap(err, "failed getting users"))
			return
		}

		api.write(w, r, http.StatusOK, src.MapSlice(users, func(user *domain.User) *UserView {
			return newUserView(user)
		}))
	}
}
