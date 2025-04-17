package v1

import (
	"context"

	userservicev1 "github.com/MaxFando/lms/user-service/api/grpc/gen/go/user-service/v1"
)

type Server struct {
	userservicev1.UserServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Echo(ctx context.Context, req *userservicev1.EchoRequest) (*userservicev1.EchoResponse, error) {
	return &userservicev1.EchoResponse{
		Message: req.GetMessage(),
	}, nil
}
