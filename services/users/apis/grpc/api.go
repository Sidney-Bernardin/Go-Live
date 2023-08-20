package grpc

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"users/configuration"
	"users/domain"
	pb "users/proto/users.com"
)

type API interface {
	Serve() error
	Shutdown()
}

type api struct {
	pb.UnimplementedUsersServer
	server  *grpc.Server
	service domain.Service
	logger  *zerolog.Logger

	address string
}

func NewAPI(config *configuration.Config, l *zerolog.Logger, svc domain.Service) API {

	// Create an api.
	a := &api{
		logger:  l,
		server:  grpc.NewServer(),
		service: svc,
		address: fmt.Sprintf("0.0.0.0:%v", config.GRPCPort),
	}

	if config.Mock {
		pb.RegisterUsersServer(a.server, &mockAPI{})
	} else {
		pb.RegisterUsersServer(a.server, a)
	}

	reflection.Register(a.server)
	return a
}

// Start starts the api's server.
func (a *api) Serve() error {

	// Create a listener.
	ln, err := net.Listen("tcp", a.address)
	if err != nil {
		return errors.Wrap(err, "cannot create listener")
	}

	// Start the server.
	err = a.server.Serve(ln)
	return errors.Wrap(err, "cannot serve")
}

// Shutdown gracefully shuts down the api's server.
func (a *api) Shutdown() {
	a.server.GracefulStop()
}
