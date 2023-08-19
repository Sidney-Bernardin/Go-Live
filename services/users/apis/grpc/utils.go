package grpc

import (
	"users/domain"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "users/proto/users.com"
)

// err returns a Status error with e as a detail. If e wasn't caused by a
// ProblemDetail, it's is logged and a new server error ProblemDetail is used instead.
func (a *api) err(e error) error {

	// If the error wasn't caused by a ProblemDetail, treat it as a server error.
	pd, ok := errors.Cause(e).(domain.ProblemDetail)
	if !ok {
		a.logger.Error().Stack().Err(e).Msg("Server Error")
		pd = domain.ProblemDetail{
			Problem: domain.ProblemServerError,
			Detail:  "Server Error",
		}
	}

	// Create a Status with the ProblemDetail.
	s, err := status.New(pd.GRPCStatusCode(), pd.Error()).
		WithDetails(&pb.ProblemDetail{
			Problem: pd.Problem,
			Detail:  pd.Detail,
		})

	if err != nil {
		err = errors.Wrap(err, "cannot add details to GRPC status")
		a.logger.Error().Stack().Err(e).Msg("Server Error")
		s = status.New(codes.Internal, "Server Error")
	}

	return s.Err()
}
