package v1

import (
	"context"
	exportservicev1 "github.com/MaxFando/lms/export-service/api/grpc/gen/go/export-service/v1"
)

type Server struct {
	exportservicev1.ExportServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Echo(ctx context.Context, req *exportservicev1.EchoRequest) (*exportservicev1.EchoResponse, error) {
	return &exportservicev1.EchoResponse{
		Message: req.GetMessage(),
	}, nil
}
