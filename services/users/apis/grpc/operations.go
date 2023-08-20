package grpc

import (
	"context"

	pb "users/proto/users.com"
)

func (a *api) AuthenticateUser(ctx context.Context, req *pb.AuthenticateUserRequest) (*pb.User, error) {

	// Authenticate the Session-ID's User.
	user, err := a.service.AuthenticateUser(ctx, req.GetSessionID(), req.GetFields()...)
	if err != nil {
		return nil, a.err(err)
	}

	return &pb.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
