package v1

import (
	"context"
	"errors"
	ticketservicev1 "github.com/MaxFando/lms/ticket-service/api/grpc/gen/go/ticket-service/v1"
	"github.com/MaxFando/lms/ticket-service/internal/converter"
	"github.com/MaxFando/lms/ticket-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

type Server struct {
	ticketservicev1.UnimplementedTicketServiceServer
	uc *usecase.TicketUsecase
}

func NewServer(uc *usecase.TicketUsecase) *Server {
	return &Server{uc: uc}
}

func (s *Server) GetTicket(ctx context.Context, req *ticketservicev1.GetTicketRequest) (*ticketservicev1.Ticket, error) {
	t, err := s.uc.GetTicket(ctx, req.TicketId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetTicket: %v", err)
	}

	return converter.ToTicketServiceFromEntity(t), nil
}

func (s *Server) CreateTicket(ctx context.Context, req *ticketservicev1.CreateTicketRequest) (*ticketservicev1.Ticket, error) {
	for i, sNum := range req.Numbers {
		_, err := strconv.Atoi(sNum)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "number[%d] invalid: %v", i, err)
		}
	}

	t, err := s.uc.CreateTicket(ctx, req.UserId, req.DrawId, req.Numbers)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrDrawNotActive):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		case errors.Is(err, usecase.ErrInvalidNumbers):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "CreateTicket: %v", err)
		}
	}

	return converter.ToTicketServiceFromEntity(t), nil
}

func (s *Server) ReserveTicket(ctx context.Context, req *ticketservicev1.ReserveTicketRequest) (*ticketservicev1.Ticket, error) {
	t, err := s.uc.ReserveTicket(ctx, req.TicketId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ReserveTicket: %v", err)
	}

	return converter.ToTicketServiceFromEntity(t), nil
}

func (s *Server) ListUserTickets(ctx context.Context, req *ticketservicev1.ListUserTicketsRequest) (*ticketservicev1.ListUserTicketsResponse, error) {
	tickets, err := s.uc.ListUserTickets(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ListUserTickets: %v", err)
	}

	resp := &ticketservicev1.ListUserTicketsResponse{}
	for _, t := range tickets {
		resp.Tickets = append(resp.Tickets, converter.ToTicketWithDrawServiceFromEntity(t))
	}
	return resp, nil
}

func (s *Server) ListAvailableTickets(ctx context.Context, req *ticketservicev1.ListAvailableTicketsRequest) (*ticketservicev1.ListAvailableTicketsResponse, error) {
	tickets, err := s.uc.ListAvailableTickets(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ListAvailableTickets: %v", err)
	}
	resp := &ticketservicev1.ListAvailableTicketsResponse{}
	for _, t := range tickets {
		resp.Tickets = append(resp.Tickets, converter.ToTicketServiceFromEntity(t))
	}
	return resp, nil
}

func (s *Server) SetWinningTickets(ctx context.Context, req *ticketservicev1.SetWinningTicketsRequest) (*ticketservicev1.SetWinningTicketsResponse, error) {
	tickets, err := s.uc.SetWinningTickets(ctx, req.TicketIds)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "SetWinningTickets: %v", err)
	}
	resp := &ticketservicev1.SetWinningTicketsResponse{}
	for _, t := range tickets {
		resp.Tickets = append(resp.Tickets, converter.ToTicketServiceFromEntity(t))
	}
	return resp, nil
}

func (s *Server) CheckResult(ctx context.Context, req *ticketservicev1.CheckResultRequest) (*ticketservicev1.CheckResultResponse, error) {
	st, err := s.uc.CheckResult(ctx, req.TicketId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CheckResult: %v", err)
	}
	return &ticketservicev1.CheckResultResponse{
		Status: st,
	}, nil
}
