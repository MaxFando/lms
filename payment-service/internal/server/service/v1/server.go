package v1

import (
	"context"

	paymentservicev1 "github.com/MaxFando/lms/payment-service/api/grpc/gen/go/payment-service/v1"
)

type Server struct {
	paymentservicev1.PaymentServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Echo(ctx context.Context, req *paymentservicev1.EchoRequest) (*paymentservicev1.EchoResponse, error) {
	return &paymentservicev1.EchoResponse{
		Message: req.GetMessage(),
	}, nil
}
