package http

import (
	"net/http"
	"time"
	"users/src/domain"
	"users/src/domain/service"

	"github.com/pkg/errors"
)

func (api *API) newSessionCookie(sessionID domain.UUID) *http.Cookie {
	return &http.Cookie{
		Name:     sessionIDCookieName,
		Value:    sessionID.String(),
		Domain:   api.config.HttpSessionCookieDomain,
		Path:     "/",
		Expires:  time.Now().Add(api.config.SessionDuration),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
}

func (api *API) handleSignup(w http.ResponseWriter, r *http.Request) {

	profilePicture, _, err := r.FormFile("profile_picture")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			api.err(w, &apiError[apiErrorType]{
				Type:    apiErrorTypeUnauthorized,
				Message: "Profile picture is required.",
				Details: map[string]any{},
			})
			return
		}

		api.err(w, errors.Wrap(err, "failed parsing profile_picture"))
		return
	}
	defer profilePicture.Close()

	form := &service.SignupForm{
		Username:       r.FormValue("username"),
		Email:          r.FormValue("email"),
		Password:       r.FormValue("password"),
		ProfilePicture: profilePicture,
	}

	session, err := api.service.Signup(r.Context(), form)
	if err != nil {
		api.err(w, errors.Wrap(err, "failed signing up"))
		return
	}

	http.SetCookie(w, api.newSessionCookie(session.ID))
}

func (api *API) handleSignin(w http.ResponseWriter, r *http.Request) {
	form := &service.SigninForm{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	session, err := api.service.Signin(r.Context(), form)
	if err != nil {
		api.err(w, errors.Wrap(err, "failed signing up"))
		return
	}

	http.SetCookie(w, api.newSessionCookie(session.ID))
}
