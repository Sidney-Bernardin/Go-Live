package domain

import (
	"io"
	"time"
	"users/src"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        UUID
	CreatedAt time.Time
	UpdatedAt time.Time

	Username     Username
	Email        string
	PasswordHash PasswordHash
	PasswordSalt PasswordSalt
}

type Username string

func NewUsername(u string) (Username, error) {
	if len(u) < 3 || 32 < len(u) {
		return "", &DomainError{
			Type:    DomainErrorTypeUsernameInvalid,
			Message: "Username must be between 3 and 32 characters.",
			Details: map[string]any{
				"username": u,
			},
		}
	}
	return Username(u), nil
}

type Password string

func NewPassword(pw string) (Password, error) {
	if len(pw) < 8 || 100 < len(pw) {
		return "", &DomainError{
			Type:    DomainErrorTypePasswordInvalid,
			Message: "Password must be between 8 and 100 characters.",
			Details: map[string]any{
				"password": pw,
			},
		}
	}
	return Password(pw), nil
}

type PasswordSalt string

func NewPasswordSalt() PasswordSalt {
	return PasswordSalt(src.MustRandomString(32))
}

type PasswordHash []byte

func NewPasswordHash(pw Password, pwSalt PasswordSalt) (PasswordHash, error) {
	pwHash := []byte(string(pw) + string(pwSalt))
	pwHash, err := bcrypt.GenerateFromPassword(pwHash, 12)
	if err != nil {
		return nil, errors.Wrap(err, "failed hashing passwaord")
	}
	return PasswordHash(pwHash), nil
}

func (pwHash PasswordHash) Compare(pw Password, pwSalt PasswordSalt) bool {
	return nil == bcrypt.CompareHashAndPassword(pwHash, []byte(string(pw)+string(pwSalt)))
}

type ProfilePicture io.Reader

func NewProfilePicture(config *src.Config, pp io.Reader) (ProfilePicture, error) {
	// b := make([]byte, config.ProfilePictureMaxBytes/10)
	// _, err := io.LimitReader(pp, int64(config.ProfilePictureMaxBytes)).Read(b)
	// if err != nil {
	// 	if errors.Is(err, io.ErrUnexpectedEOF) {
	// 		return nil, &DomainError{
	// 			Type:    DomainErrorTypeProfilePictureInvalid,
	// 			Message: fmt.Sprintf("The profile picture is too large. Must be below %v bytes", config.ProfilePictureMaxBytes),
	// 		}
	// 	}
	// 	return nil, errors.Wrap(err, "failed reading")
	// }
	return ProfilePicture(pp), nil
}
