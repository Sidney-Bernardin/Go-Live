package domain

import (
	"bytes"
	"fmt"
	"math"
	"time"
	"users/configuration"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
)

type Service interface {
	Signup(info *SignupInfo) (sessionID string, err error)
	Signin(info *SigninInfo) (ssessionID string, err error)
	Logout(sessionID string) error

	GetUser(search, by string, fields ...string) (*User, error)
	SearchUsers(username string, offset, limit int, fields ...string) ([]*User, error)
	GetSelf(sessionID string, fields ...string) (*User, error)

	SetProfilePicture(sessionID string, picture []byte) error
	GetProfilePicture(userID string) (*bytes.Buffer, error)
}

type service struct {
	config       *configuration.Configuration
	databaseRepo DatabaseRepository
}

func NewService(config *configuration.Configuration, dbRepo DatabaseRepository) Service {
	return &service{config, dbRepo}
}

func (svc *service) Signup(info *SignupInfo) (string, error) {

	// Check the length of the username.
	if len(info.Username) < 3 || len(info.Username) > 20 {
		return "", ProblemDetail{
			Type:   PDTypeInvalidSignupInfo,
			Detail: "Username must be between 3 and 20 characters long.",
		}
	}

	// Check for taken user fields.
	fields := map[string]any{"username": info.Username, "email": info.Email}
	if err := svc.databaseRepo.CheckForTakenUserFields(fields); err != nil {

		// If the error was caused by a problem-detail that has one of the
		// following types, replace it's type.
		if pd, ok := errors.Cause(err).(ProblemDetail); ok && pd.hasType(PDTypeFieldTaken) {
			pd.Type = PDTypeInvalidSignupInfo
			return "", pd
		}

		return "", errors.Wrap(err, "cannot check for taken user fields")
	}

	// Hash the password.
	hashedPasw, err := bcrypt.GenerateFromPassword([]byte(info.Password), 14)
	if err != nil {
		return "", errors.Wrap(err, "cannot hash password")
	}

	// Create a user.
	user := &User{
		Username: info.Username,
		Email:    info.Email,
		Password: string(hashedPasw),
	}

	// Insert the user with a new session.
	sessionID, err := svc.databaseRepo.InsertUserWithSession(user, &Session{
		ExpireAt: time.Now().Add(svc.config.SessionLength),
	})

	return sessionID, errors.Wrap(err, "cannot insert user with session")
}

func (svc *service) Signin(info *SigninInfo) (string, error) {

	// Get the user from the database.
	user, err := svc.databaseRepo.GetUserByUsername(info.Username, "password")
	if err != nil {

		// If the error was caused by a problem-detail that has one of the
		// following types, replace it's type and detail.
		if pd, ok := errors.Cause(err).(ProblemDetail); ok && pd.hasType(PDTypeUserDoesntExist) {
			pd.Type = PDTypeInvalidSigninInfo
			pd.Detail = "Incorrect username or password."
			return "", pd
		}

		return "", errors.Wrap(err, "cannot get user by username")
	}

	// Compare the passwords.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(info.Password)); err != nil {

		// Check if the passwords don't match.
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return "", ProblemDetail{
				Type:   PDTypeInvalidSigninInfo,
				Detail: "Incorrect username or password.",
			}
		}

		return "", errors.Wrap(err, "cannot compare hash and password")
	}

	// Create a session for the user and insert it into the database.
	sessionID, err := svc.databaseRepo.InsertSession(&Session{
		UserID:   user.ID,
		ExpireAt: time.Now().Add(svc.config.SessionLength),
	})

	return sessionID, errors.Wrap(err, "cannot insert session")
}

func (svc *service) Logout(sessionID string) error {

	// Delete the session from the database.
	err := svc.databaseRepo.DeleteSession(sessionID)

	// If the error was caused by a problem-detail that has one of the
	// following types, replace it's type and detail.
	if pd, ok := errors.Cause(err).(ProblemDetail); ok && pd.hasType(PDTypeInvalidID) {
		pd.Type = PDTypeUnauthorized
		return pd
	}

	return errors.Wrap(err, "cannot delete session")
}

func (svc *service) GetUser(search, by string, fields ...string) (user *User, err error) {

	// Make sure the user's password isn't returned, by removing the password field.
	if i := slices.Index(fields, "password"); i != -1 {
		fields[i] = ""
	}

	// Get the user by ID or username.
	switch by {
	case "id":
		user, err = svc.databaseRepo.GetUser(search, fields...)
	case "username":
		user, err = svc.databaseRepo.GetUserByUsername(search, fields...)
	default:
		return nil, ProblemDetail{
			Type:   PDTypeInvalidInput,
			Detail: fmt.Sprintf("Cannot get user by %s.", by),
		}
	}

	return user, errors.Wrapf(err, "cannot get user by %s", by)
}

func (svc *service) SearchUsers(username string, offset, limit int, fields ...string) ([]*User, error) {

	// Make sure the user passwords aren't returned, by removing the password field.
	if i := slices.Index(fields, "password"); i != -1 {
		fields[i] = ""
	}

	// Get offset and limit's absolute and int64 values.
	offset64 := int64(math.Abs(float64(offset)))
	limit64 := int64(math.Abs(float64(limit)))

	// Search for the users from the database.
	users, err := svc.databaseRepo.SearchUsers(username, offset64, limit64, fields...)
	return users, errors.Wrap(err, "cannot search fields")
}

func (svc *service) GetSelf(sessionID string, fields ...string) (*User, error) {

	// Make sure the user's password isn't returned, by removing the password field.
	if i := slices.Index(fields, "password"); i != -1 {
		fields[i] = ""
	}

	// Get the session's user from the database.
	user, err := svc.databaseRepo.GetSessionsUser(sessionID, fields...)

	// If the error was caused by a problem-detail that has one of the following
	// types, replace it's type.
	if pd, ok := errors.Cause(err).(ProblemDetail); ok && pd.hasType(PDTypeInvalidID, PDTypeSessionDoesntExist, PDTypeUserDoesntExist) {
		pd.Type = PDTypeUnauthorized
		return nil, pd
	}

	return user, errors.Wrap(err, "cannot get sessions user")
}

func (svc *service) SetProfilePicture(sessionID string, profilePicture []byte) error {

	// Get the session's user from the database.
	user, err := svc.databaseRepo.GetSessionsUser(sessionID)
	if err != nil {

		// If the error was caused by a problem-detail that has one of the
		// following types, replace it's type.
		if pd, ok := errors.Cause(err).(ProblemDetail); ok && pd.hasType(PDTypeInvalidID, PDTypeSessionDoesntExist, PDTypeUserDoesntExist) {
			pd.Type = PDTypeUnauthorized
			return pd
		}

		return errors.Wrap(err, "cannot get sessions user")
	}

	// Insert the profile-picture into the database.
	err = svc.databaseRepo.InsertProfilePicture(user.ID, profilePicture)
	return errors.Wrap(err, "cannot insert profile-picture")
}

func (svc *service) GetProfilePicture(userID string) (*bytes.Buffer, error) {

	// Get the profile-picture from the database.
	buf, err := svc.databaseRepo.GetProfilePicture(userID)
	return buf, errors.Wrap(err, "cannot get profile-picture")
}
