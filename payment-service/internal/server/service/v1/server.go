package v1

import (
	"context"

	api "github.com/MaxFando/lms/payment-service/api/grpc/gen/go/payment-service/v1"
	"github.com/MaxFando/lms/payment-service/internal/entity"
	"github.com/shopspring/decimal"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type service interface {
	CreateInvoice(ctx context.Context, userId int64, ticketId int64) (int64, decimal.Decimal, error)
	CreateInvoiceForBookedTicket(ctx context.Context, userId int64, ticketId int64) (int64, decimal.Decimal, error)
	Pay(ctx context.Context, userId int64, invoiceId int64, card *entity.Card) error
}

type Server struct {
	api.PaymentServiceServer
	service service
}

func NewServer(service service) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) CreateInvoice(ctx context.Context, req *api.CreateInvoiceRequest) (*api.CreateInvoiceResponse, error) {
	invoiceID, price, err := s.service.CreateInvoice(ctx, req.GetUserId(), req.GetTicketId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &api.CreateInvoiceResponse{
		Id:    invoiceID,
		Price: decimalToMoney(price),
	}, nil
}

func (s *Server) CreateInvoiceInternal(ctx context.Context, req *api.CreateInvoiceRequest) (*api.CreateInvoiceResponse, error) {
	invoiceID, price, err := s.service.CreateInvoiceForBookedTicket(ctx, req.GetUserId(), req.GetTicketId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &api.CreateInvoiceResponse{
		Id:    invoiceID,
		Price: decimalToMoney(price),
	}, nil
}

func (s *Server) Pay(ctx context.Context, req *api.PayRequest) (*emptypb.Empty, error) {
	card := &entity.Card{
		Number:  req.GetCardNumber(),
		ExpDate: req.GetExpDate(),
		CVV:     req.GetCVV(),
	}

	if err := card.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid card data: %s", err.Error())
	}

	err := s.service.Pay(ctx, req.GetUserId(), req.GetInvoiceId(), card)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func decimalToMoney(d decimal.Decimal) *money.Money {
	units := d.Truncate(0).IntPart()
	nanosDecimal := d.Sub(decimal.NewFromInt(units))
	nanos := nanosDecimal.Mul(decimal.NewFromInt(1_000_000_000)).IntPart()

	return &money.Money{
		CurrencyCode: "RUB",
		Units:        units,
		Nanos:        int32(nanos),
	}
}
