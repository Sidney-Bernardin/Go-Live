package blobstore

import (
	"bytes"
	"fmt"
	"testing"
	"time"
	"users/src/domain"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertProfilePicture(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	repo := suite(t)

	t.Run("foo", func(t *testing.T) {
		userID := domain.NewUUID()
		profilePicture := domain.ProfilePicture(bytes.NewReader([]byte(`foobarbaz`)))

		err := repo.InsertProfilePicture(t.Context(), userID, profilePicture)
		require.NoError(t, err)

		err = s3.NewObjectExistsWaiter(repo.client).
			Wait(ctx, &s3.HeadObjectInput{Bucket: &usersBucket, Key: aws.String(fmt.Sprintf("%s.png", userID))}, 10*time.Second)
		assert.NoError(t, err)
	})
}
