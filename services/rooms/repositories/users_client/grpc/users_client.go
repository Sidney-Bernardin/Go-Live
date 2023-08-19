package grpc

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"rooms/configuration"
	"rooms/domain"
	pb "rooms/proto/users.com"
)

type usersClientRepository struct {
	client pb.UsersClient
}

func NewUsersClientRepository(config *configuration.Config) (domain.UsersClientRepository, error) {

	// Connect to Users Service.
	conn, err := grpc.Dial(config.UsersGRPCUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrap(err, "cannot dial users service")
	}

	return &usersClientRepository{pb.NewUsersClient(conn)}, nil
}

// AuthenticateUser authenticates the Session-ID's User.
func (repo *usersClientRepository) AuthenticateUser(ctx context.Context, sessionID string, fields ...string) (*domain.User, error) {

	req := &pb.AuthenticateUserRequest{
		SessionID: sessionID,
		Fields:    fields,
	}

	// Authenticate the Session-ID's User.
	res, err := repo.client.AuthenticateUser(ctx, req)
	if err != nil {

		// Check if the error has a ProblemDetail.
		if pd := repo.getProblemDetail(err); pd != nil {
			return nil, *pd
		}

		return nil, errors.Wrap(err, "cannot authenticate user")
	}

	return &domain.User{
		ID:       res.ID,
		Username: res.Username,
		Email:    res.Email,
	}, nil
}
