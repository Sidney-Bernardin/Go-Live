package blobstore

import (
	"context"
	"fmt"
	"users/src/domain"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

func (repo *repository) InsertProfilePicture(ctx context.Context, userID domain.UUID, profilePicture domain.ProfilePicture) error {
	key := fmt.Sprintf("profiles/pictures/%s.png", userID)

	// Put the profile-picture.
	_, err := repo.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &usersBucket,
		Key:    &key,
		Body:   profilePicture,
	})

	return errors.Wrap(err, "failed putting")
}
