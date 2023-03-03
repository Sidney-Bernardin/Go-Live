package grpc

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"rooms/configuration"
	"rooms/domain"
	pb "rooms/proto/users.com"
)

type usersClientRepository struct {
	client pb.UsersClient
}

func NewUsersClientRepository(config *configuration.Configuration) (domain.UsersClientRepository, error) {

	// Connect to the users service.
	conn, err := grpc.Dial(config.UsersGRPCAddr, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "cannot dial users service")
	}

	return &usersClientRepository{pb.NewUsersClient(conn)}, nil
}

func (repo *usersClientRepository) getProblemDetail(err error) (domain.ProblemDetail, bool) {

	s, ok := status.FromError(err)
	if !ok {
		return domain.ProblemDetail{}, false
	}

	for _, detail := range s.Details() {
		if problemDetail, ok := detail.(*pb.ProblemDetail); ok {
			return domain.ProblemDetail{
				Type:   problemDetail.Type,
				Detail: problemDetail.Detail,
			}, true
		}
	}

	return domain.ProblemDetail{}, false
}

func (repo *usersClientRepository) GetSelf(sessionID string, fields []string) (*domain.User, error) {

	req := &pb.GetSelfRequest{
		SessionID: sessionID,
		Fields:    fields,
	}

	res, err := repo.client.GetSelf(context.Background(), req)
	if err != nil {

		if pd, ok := repo.getProblemDetail(err); ok {
			return nil, pd
		}

		return nil, errors.Wrap(err, "cannot get self")
	}

	user := &domain.User{
		ID:       res.ID,
		Username: res.Username,
		Email:    res.Email,
	}

	return user, nil
}
