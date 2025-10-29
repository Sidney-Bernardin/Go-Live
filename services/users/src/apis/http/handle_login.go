package http

import (
	"net/http"
	"time"
	"users/src/domain"

	"github.com/pkg/errors"
)

func (api *API) newSessionCookie(session *domain.Session) *http.Cookie {
	return &http.Cookie{
		Name:     "SESSION_ID",
		Value:    session.ID.String(),
		Domain:   api.config.HttpSessionCookieDomain,
		Expires:  time.Now().Add(api.config.SessionDuration),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func (api *API) handleSignup(w http.ResponseWriter, r *http.Request) {
	signupForm := &domain.SignupForm{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	session, err := api.svc.Signup(r.Context(), signupForm)
	if err != nil {
		api.err(w, errors.Wrap(err, "failed signing up"))
		return
	}

	http.SetCookie(w, api.newSessionCookie(session))
}
