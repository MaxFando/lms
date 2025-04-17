package v1

import (
	"context"

	drawresultservicev1 "github.com/MaxFando/lms/draw-result-service/api/grpc/gen/go/draw-result-service/v1"
)

type DrawResultServer struct {
	drawresultservicev1.DrawResultServiceServer
}

func NewDrawResultServer() *DrawResultServer {
	return &DrawResultServer{}
}

func (s *DrawResultServer) Echo(ctx context.Context, req *drawresultservicev1.EchoRequest) (*drawresultservicev1.EchoResponse, error) {
	return &drawresultservicev1.EchoResponse{
		Message: req.GetMessage(),
	}, nil
}
