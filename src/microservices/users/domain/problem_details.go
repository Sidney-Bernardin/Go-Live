package domain

import "fmt"

const (
	PDTypeUnauthorized = "unauthorized"

	PDTypeInvalidInput = "invalid_input"
	PDTypeInvalidID    = "invalid_id"

	PDTypeFieldTaken = "field_taken"

	PDTypeInvalidSignupInfo = "invalid_signup_info"
	PDTypeInvalidSigninInfo = "invalid_signin_info"

	PDTypeUserDoesntExist           = "user_doesnt_exist"
	PDTypeSessionDoesntExist        = "session_doesnt_exist"
	PDTypeProfilePictureDoesntExist = "profile_picture_doesnt_exist"
)

type ProblemDetail struct {
	Type   string `json:"type,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func (pd ProblemDetail) Error() string {
	return fmt.Sprintf("%s: '%s'", pd.Type, pd.Detail)
}

// hasType returns true if the ProblemDetail's type matches on of the given types.
func (pd ProblemDetail) hasType(types ...string) bool {

	for _, t := range types {
		if pd.Type == t {
			return true
		}
	}

	return false
}
