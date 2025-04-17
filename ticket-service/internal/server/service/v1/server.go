package v1

import (
	"context"

	ticketservicev1 "github.com/MaxFando/lms/ticket-service/api/grpc/gen/go/ticket-service/v1"
)

type Server struct {
	ticketservicev1.TicketServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Echo(ctx context.Context, req *ticketservicev1.EchoRequest) (*ticketservicev1.EchoResponse, error) {
	return &ticketservicev1.EchoResponse{
		Message: req.GetMessage(),
	}, nil
}
