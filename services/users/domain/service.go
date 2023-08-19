package domain

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"time"
	"users/configuration"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
)

type Service interface {
	AuthenticateUser(ctx context.Context, sessionID string, fields ...string) (*User, error)

	Signup(ctx context.Context, info *SignupInfo) (sessionID string, err error)
	Signin(ctx context.Context, info *SigninInfo) (ssessionID string, err error)
	Logout(ctx context.Context, sessionID string) error

	GetUser(ctx context.Context, userID string, fields ...string) (*User, error)
	SearchUsers(ctx context.Context, username string, offset, limit int, fields ...string) ([]*User, error)

	UpdateProfilePicture(ctx context.Context, sessionID string, profilePicture []byte) error
	GetProfilePicture(ctx context.Context, userID string) (*bytes.Buffer, error)
}

type service struct {
	config       *configuration.Config
	databaseRepo DatabaseRepository
}

func NewService(config *configuration.Config, dbRepo DatabaseRepository) Service {
	return &service{config, dbRepo}
}

// AuthenticateUser gets sessionID's User from the database.
func (svc *service) AuthenticateUser(ctx context.Context, sessionID string, fields ...string) (*User, error) {

	if i := slices.Index(fields, "password"); i != -1 {
		fields[i] = ""
	}

	// Get the Session-ID's User from the database.
	user, err := svc.databaseRepo.GetUserBySession(ctx, sessionID, fields...)
	if err != nil {

		// If the error was caused by a ProblemDetail, replace it.
		if _, ok := errors.Cause(err).(ProblemDetail); ok {
			return nil, ProblemDetail{Problem: ProblemUnauthorized}
		}

		return nil, errors.Wrap(err, "cannot get user")
	}

	return user, nil
}

// Signup uses info to create a User and a Session, then inserts them into the database.
func (svc *service) Signup(ctx context.Context, info *SignupInfo) (string, error) {

	// Check the length of the SignupInfo's username.
	if len(info.Username) < 3 || len(info.Username) > 20 {
		return "", ProblemDetail{
			Problem: ProblemInvalidSignupInfo,
			Detail:  "Username must be between 3 and 20 characters long.",
		}
	}

	// Check if the SignupInfo has a taken username or email.
	userExists, takenField, err := svc.databaseRepo.UserExists(ctx,
		"username", info.Username,
		"email", info.Email,
	)

	if err != nil {
		return "", errors.Wrap(err, "cannot check if user exists")
	}

	if userExists {
		return "", ProblemDetail{
			Problem: ProblemInvalidSignupInfo,
			Detail:  fmt.Sprintf("The %s has been taken by another user.", takenField),
		}
	}

	// Hash the SignupInfo's password.
	hashedPasw, err := bcrypt.GenerateFromPassword([]byte(info.Password), 14)
	if err != nil {
		return "", errors.Wrap(err, "cannot hash password")
	}

	// Create a User with the SignupInfo and insert it into the database.
	userID, err := svc.databaseRepo.InsertUser(ctx, &User{
		Username: info.Username,
		Email:    info.Email,
		Password: string(hashedPasw),
	})

	if err != nil {
		return "", errors.Wrap(err, "cannot insert user")
	}

	// Create a Session for the User and insert it into the database.
	sessionID, err := svc.databaseRepo.InsertSession(ctx, &Session{
		UserID:   userID,
		ExpireAt: time.Now().Add(svc.config.SessionLength),
	})

	return sessionID, errors.Wrap(err, "cannot insert session")
}

// Signin creates a Session for the User with info's username, and inserts it
// into the database.
func (svc *service) Signin(ctx context.Context, info *SigninInfo) (string, error) {

	// Get the User with the SigninInfo's username.
	user, err := svc.databaseRepo.GetUserByUsername(ctx, info.Username, "password")
	if err != nil {

		// If the error was caused by a ProblemDetail, replace it.
		if _, ok := err.(ProblemDetail); ok {
			return "", ProblemDetail{
				Problem: ProblemInvalidSigninInfo,
				Detail:  "Incorrect username or password.",
			}
		}

		return "", errors.Wrap(err, "cannot get user")
	}

	// Check if the SigninInfo's password matches the User's password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(info.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return "", ProblemDetail{
				Problem: ProblemInvalidSigninInfo,
				Detail:  "Incorrect username or password.",
			}
		}

		return "", errors.Wrap(err, "cannot compare passwords")
	}

	// Create a Session for the User and insert it into the database.
	sessionID, err := svc.databaseRepo.InsertSession(ctx, &Session{
		UserID:   user.ID,
		ExpireAt: time.Now().Add(svc.config.SessionLength),
	})

	return sessionID, errors.Wrap(err, "cannot insert session")
}

// Delete deletes sessionID's Session from the database.
func (svc *service) Logout(ctx context.Context, sessionID string) error {
	if err := svc.databaseRepo.DeleteSession(ctx, sessionID); err != nil {

		// If the error was caused by a ProblemDetail, replace it.
		if _, ok := err.(ProblemDetail); ok {
			return ProblemDetail{Problem: ProblemUnauthorized}
		}

		return errors.Wrap(err, "cannot delete session")
	}

	return nil
}

// GetUser gets userID's User form the database.
func (svc *service) GetUser(ctx context.Context, userID string, fields ...string) (user *User, err error) {

	if i := slices.Index(fields, "password"); i != -1 {
		fields[i] = ""
	}

	// Get the user-ID's User from the database.
	user, err = svc.databaseRepo.GetUser(ctx, userID, fields...)
	return user, errors.Wrap(err, "cannot get user")
}

// SearchUsers gets Users with usernames similar to the one given, from the database.
func (svc *service) SearchUsers(ctx context.Context, username string, offset, limit int, fields ...string) ([]*User, error) {

	if i := slices.Index(fields, "password"); i != -1 {
		fields[i] = ""
	}

	// Get absolute values of the offset and limit.
	absOffset := int64(math.Abs(float64(offset)))
	absLimit := int64(math.Abs(float64(limit)))

	// Get Users with usernames similar to the username.
	users, err := svc.databaseRepo.SearchUsers(ctx, username, absOffset, absLimit, fields...)
	return users, errors.Wrap(err, "cannot search fields")
}

// UpdateProfilePicture updates the profile-picture of sessionID's User in the database.
func (svc *service) UpdateProfilePicture(ctx context.Context, sessionID string, profilePicture []byte) error {

	user, err := svc.AuthenticateUser(ctx, sessionID)
	if err != nil {
		return errors.Wrap(err, "cannot authenticate user")
	}

	// Update the profile-picture of the User in the database.
	err = svc.databaseRepo.UpdateProfilePicture(ctx, user.ID, profilePicture)
	return errors.Wrap(err, "cannot insert profile-picture")
}

// GetProfilePicture gets the profile picture of userID's User from the database.
func (svc *service) GetProfilePicture(ctx context.Context, userID string) (*bytes.Buffer, error) {
	buf, err := svc.databaseRepo.GetProfilePicture(ctx, userID)
	return buf, errors.Wrap(err, "cannot get profile-picture")
}
