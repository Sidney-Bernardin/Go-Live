package grpc

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"rooms/configuration"
	"rooms/domain"
	pb "rooms/proto/users.com"
)

type usersClientRepository struct {
	client pb.UsersClient
}

func NewUsersClientRepository(config *configuration.Configuration) (domain.UsersClientRepository, error) {

	// Connect to the users service.
	conn, err := grpc.Dial(config.UsersGRPCUrl, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "cannot dial users service")
	}

	return &usersClientRepository{pb.NewUsersClient(conn)}, nil
}

func (repo *usersClientRepository) GetSelf(sessionID string, fields []string) (*domain.User, error) {

	req := &pb.GetSelfRequest{
		SessionID: sessionID,
		Fields:    fields,
	}

	res, err := repo.client.GetSelf(context.Background(), req)
	if err != nil {

		// Check if the error has a problem-detail.
		if pd, ok := repo.getProblemDetail(err); ok {
			return nil, pd
		}

		return nil, errors.Wrap(err, "cannot get self")
	}

	return &domain.User{
		ID:       res.ID,
		Username: res.Username,
		Email:    res.Email,
	}, nil
}
