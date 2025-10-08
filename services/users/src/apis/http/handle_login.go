package http

import (
	"encoding/json"
	"net/http"
	"time"
	"users/src/domain"

	"github.com/pkg/errors"
)

func (api *API) newSessionCookie(session *domain.Session) *http.Cookie {
	return &http.Cookie{
		Name:     "SESSION_ID",
		Value:    session.ID.String(),
		Domain:   api.config.HTTPSessionCookieDomain,
		Expires:  time.Now().Add(api.config.SessionDuration),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func (api *API) handleSignup(w http.ResponseWriter, r *http.Request) {

	var signupForm *domain.SignupForm
	if err := json.NewDecoder(r.Body).Decode(&signupForm); err != nil {
		api.err(w, r, errors.Wrap(err, "failed decoding signup-form"))
		return
	}

	session, err := api.svc.Signup(r.Context(), signupForm)
	if err != nil {
		api.err(w, r, errors.Wrap(err, "failed signing up"))
		return
	}

	http.SetCookie(w, api.newSessionCookie(session))
}
