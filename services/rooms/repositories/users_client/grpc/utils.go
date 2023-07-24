package grpc

import (
	"rooms/domain"

	"google.golang.org/grpc/status"

	pb "rooms/proto/users.com"
)

// getProblemDetail returns true if the given error is a GRPC status and has a problem-detail.
func (repo *usersClientRepository) getProblemDetail(err error) (domain.ProblemDetail, bool) {

	// Check if the error is a GRPC status.
	s, ok := status.FromError(err)
	if !ok {
		return domain.ProblemDetail{}, false
	}

	for _, detail := range s.Details() {

		// Check if the detail is a problem-detail. If so return a copy of it.
		if problemDetail, ok := detail.(*pb.ProblemDetail); ok {
			return domain.ProblemDetail{
				Type:   problemDetail.Type,
				Detail: problemDetail.Detail,
			}, true
		}
	}

	return domain.ProblemDetail{}, false
}
