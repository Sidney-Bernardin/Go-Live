package grpc

import (
	"rooms/domain"

	"google.golang.org/grpc/status"

	pb "rooms/proto/users.com"
)

// getProblemDetail checks if err as a Status has a ProblemDetail, if so returns it.
func (repo *usersClientRepository) getProblemDetail(err error) *domain.ProblemDetail {

	// Create a Status from the error.
	s, ok := status.FromError(err)
	if !ok {
		return nil
	}

	// Check the Status's details for a ProblemDetail.
	for _, detail := range s.Details() {
		if pd, ok := detail.(*pb.ProblemDetail); ok {
			return &domain.ProblemDetail{
				Problem: pd.Problem,
				Detail:  pd.Detail,
			}
		}
	}

	return nil
}
