package blobstore

import (
	"testing"
	"users/src"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
)

func suite(t *testing.T) *repository {
	t.Helper()
	ctx := t.Context()

	localstackContainer, err := localstack.Run(ctx, "localstack/localstack:1.4.0")
	testcontainers.CleanupContainer(t, localstackContainer)
	require.NoError(t, err)

	endpoint, err := localstackContainer.PortEndpoint(ctx, nat.Port("4566/tcp"), "http")
	require.NoError(t, err)

	repo, err := New(ctx, &src.Config{
		AWSBaseEndpoint: endpoint,
		AWSBaseRegion:   "us-east-1"},
	)
	require.NoError(t, err)

	return repo.(*repository)
}
