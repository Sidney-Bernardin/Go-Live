package domain

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
)

const (
	ProblemServerError  = "server_error"
	ProblemUnauthorized = "unauthorized"
	ProblemTaken        = "taken"

	ProblemInvalidInput = "invalid_input"
	ProblemInvalidID    = "invalid_id"

	ProblemInvalidSignupInfo = "invalid_signup_info"
	ProblemInvalidSigninInfo = "invalid_signin_info"

	ProblemUserDoesntExist           = "user_doesnt_exist"
	ProblemSessionDoesntExist        = "session_doesnt_exist"
	ProblemProfilePictureDoesntExist = "profile_picture_doesnt_exist"
)

type ProblemDetail struct {
	Problem string `json:"problem,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

func (pd ProblemDetail) Error() string {
	return fmt.Sprintf("%s: %s", pd.Problem, pd.Detail)
}

func (pd ProblemDetail) HTTPStatusCode() int {
	switch pd.Problem {
	case ProblemServerError:
		return http.StatusInternalServerError
	case ProblemUnauthorized:
		return http.StatusUnauthorized
	case ProblemInvalidInput, ProblemInvalidID:
		return http.StatusUnprocessableEntity
	case ProblemUserDoesntExist, ProblemSessionDoesntExist, ProblemProfilePictureDoesntExist:
		return http.StatusNotFound
	case ProblemTaken:
		return http.StatusConflict
	default:
		return http.StatusBadRequest
	}
}

func (pd ProblemDetail) GRPCStatusCode() codes.Code {
	switch pd.Problem {
	case ProblemServerError:
		return codes.Internal
	case ProblemUnauthorized:
		return codes.Unauthenticated
	case ProblemInvalidInput, ProblemInvalidID:
		return codes.InvalidArgument
	case ProblemUserDoesntExist, ProblemSessionDoesntExist, ProblemProfilePictureDoesntExist:
		return codes.NotFound
	case ProblemTaken:
		return codes.Unavailable
	default:
		return codes.Unknown
	}
}
