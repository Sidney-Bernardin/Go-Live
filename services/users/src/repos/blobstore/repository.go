package blobstore

import (
	"context"
	"time"
	"users/src"
	"users/src/domain/service"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/pkg/errors"
)

var usersBucket = "users"

type repository struct {
	client *s3.Client
}

func New(ctx context.Context, config *src.Config) (service.BlobStoreRepository, error) {
	c := s3.New(s3.Options{
		Region:       config.AWSBaseRegion,
		BaseEndpoint: &config.AWSBaseEndpoint,
		UsePathStyle: true,
	})

	if err := initS3(ctx, c); err != nil {
		return nil, errors.Wrap(err, "failed initializing")
	}

	return &repository{c}, nil
}

func initS3(ctx context.Context, c *s3.Client) error {

	// Create users bucket.
	_, err := c.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &usersBucket,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint("us-east-1"),
		},
	})

	if err != nil {
		return errors.Wrap(err, "failed creating users bucket")
	}

	err = s3.NewBucketExistsWaiter(c).
		Wait(ctx, &s3.HeadBucketInput{Bucket: aws.String("users")}, 10*time.Second)

	if err != nil {
		return errors.Wrap(err, "failed waiting for users bucket")
	}

	_, err = c.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: &usersBucket,
		Policy: aws.String(`{
			"Effect": "Allow",
			"Principle" :"*",
			"Action": [
				"s3:GetObject"
			],
			"Resource": "arn:aws:s3:::users/*"
		}`),
	})

	if err != nil {
		return errors.Wrap(err, "failed waiting for users bucket")
	}

	return nil
}
