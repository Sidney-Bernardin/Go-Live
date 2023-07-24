package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"users/domain"
	pb "users/proto/users.com"
)

type api struct {
	service domain.Service
	pb.UnimplementedUsersServer
	logger *zerolog.Logger
}

func New(svc domain.Service, l *zerolog.Logger) *api {
	return &api{
		service: svc,
		logger:  l,
	}
}

func (a *api) Serve(port int) error {

	// Create a GRPC server.
	svr := grpc.NewServer()
	pb.RegisterUsersServer(svr, a)
	reflection.Register(svr)

	// Create a TCP listener for the server.
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return errors.Wrap(err, "cannot create TCP listener")
	}

	// Serve the GRPC server with the listener.
	err = svr.Serve(l)
	return errors.Wrap(err, "cannot serve")
}

// newProblemDetailStatus returns a new GRPC status containing the problem-detail.
func (a *api) newProblemDetailStatus(code codes.Code, pd domain.ProblemDetail) (*status.Status, error) {

	// Create a GRPC status with details that contain the problem-detail.
	s, err2 := status.Newf(code, pd.Error()).WithDetails(&pb.ProblemDetail{
		Type:   pd.Type,
		Detail: pd.Detail,
	})

	return s, errors.Wrap(err2, "cannot create status")
}

func (a *api) GetSelf(ctx context.Context, req *pb.GetSelfRequest) (*pb.GetSelfResponse, error) {

	user, err := a.service.GetSelf(req.GetSessionID(), req.GetFields()...)
	if err != nil {

		// Check if the error was caused by a problem-detail.
		problemDetail, ok := errors.Cause(err).(domain.ProblemDetail)
		if ok {

			// Create a GRPC status that contains the problem-detail.
			s, err2 := a.newProblemDetailStatus(codes.FailedPrecondition, problemDetail)
			if err2 != nil {
				return nil, status.Newf(codes.Internal, err2.Error()).Err()
			}

			return nil, s.Err()
		}

		return nil, status.Newf(codes.Internal, err.Error()).Err()
	}

	return &pb.GetSelfResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
