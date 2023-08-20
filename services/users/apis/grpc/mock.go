package grpc

import (
	"context"
	pb "users/proto/users.com"
)

type mockAPI struct {
	pb.UnimplementedUsersServer
}

func (a *mockAPI) AuthenticateUser(ctx context.Context, req *pb.AuthenticateUserRequest) (*pb.User, error) {
	return &pb.User{
		Username: "Hello, World!",
	}, nil
}
