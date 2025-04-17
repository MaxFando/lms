package v1

import (
	"context"

	drawresultservicev1 "github.com/MaxFando/lms/draw-service/api/grpc/gen/go/draw-service/v1"
)

type Server struct {
	drawresultservicev1.DrawServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Echo(ctx context.Context, req *drawresultservicev1.EchoRequest) (*drawresultservicev1.EchoResponse, error) {
	return &drawresultservicev1.EchoResponse{
		Message: req.GetMessage(),
	}, nil
}
