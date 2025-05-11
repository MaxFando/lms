package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	drawresultservicev1 "github.com/MaxFando/lms/draw-service/api/grpc/gen/go/draw-service/v1"
	"github.com/MaxFando/lms/draw-service/internal/entity"
	"github.com/MaxFando/lms/draw-service/internal/usecase"
)

type Server struct {
	drawresultservicev1.UnimplementedDrawServiceServer
	usecase *usecase.DrawUseCase
}

func NewServer(usecase *usecase.DrawUseCase) *Server {
	return &Server{usecase: usecase}
}

// CreateDraw CreateDraws создает новый тираж
func (s *Server) CreateDraw(ctx context.Context, req *drawresultservicev1.CreateDrawRequest) (*drawresultservicev1.DrawResponse, error) {
	startTime := req.GetStartTime().AsTime()
	endTime := req.GetEndTime().AsTime()

	draw := entity.Draw{
		LotteryType: entity.LotteryType(req.GetLotteryType()),
		StartTime:   startTime,
		EndTime:     endTime,
	}

	createdDraw, err := s.usecase.CreateDraws(ctx, draw)
	if err != nil {
		return nil, status.Error(status.Code(err), "create draw")
	}

	return &drawresultservicev1.DrawResponse{
		Id:          createdDraw.ID,
		LotteryType: string(createdDraw.LotteryType),
		StartTime:   timestamppb.New(createdDraw.StartTime),
		EndTime:     timestamppb.New(createdDraw.EndTime),
		Status:      string(createdDraw.Status),
	}, nil
}

// GetDrawsList возвращает список активных тиражей
func (s *Server) GetDrawsList(ctx context.Context, req *emptypb.Empty) (*drawresultservicev1.GetDrawsListResponse, error) {
	draws, err := s.usecase.GetDrawsList(ctx)
	if err != nil {
		return nil, fmt.Errorf("get draws list: %w", err)
	}

	respDraws := make([]*drawresultservicev1.DrawResponse, 0, len(draws))
	for _, d := range draws {
		respDraws = append(respDraws, &drawresultservicev1.DrawResponse{
			Id:          d.ID,
			LotteryType: string(d.LotteryType),
			StartTime:   timestamppb.New(d.StartTime),
			EndTime:     timestamppb.New(d.EndTime),
			Status:      string(d.Status),
		})
	}

	return &drawresultservicev1.GetDrawsListResponse{Draws: respDraws}, nil
}

// CancelDraw отменяет тираж по ID
func (s *Server) CancelDraw(ctx context.Context, req *drawresultservicev1.CancelDrawRequest) (*emptypb.Empty, error) {
	err := s.usecase.CancelDraw(ctx, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("cancel draw: %w", err)
	}

	return &emptypb.Empty{}, nil
}

// GetCompletedDrawsList - получения списка завершенных тиражей
func (s *Server) GetCompletedDrawsList(ctx context.Context, req *emptypb.Empty) (*drawresultservicev1.GetDrawsListResponse, error) {
	draws, err := s.usecase.GetCompletedDraws(ctx)
	if err != nil {
		return nil, fmt.Errorf("get completed draws list: %w", err)
	}

	var respDraws []*drawresultservicev1.DrawResponse
	for _, d := range draws {
		respDraws = append(respDraws, &drawresultservicev1.DrawResponse{
			Id:          d.ID,
			LotteryType: string(d.LotteryType),
			StartTime:   timestamppb.New(d.StartTime),
			EndTime:     timestamppb.New(d.EndTime),
			Status:      string(d.Status),
		})
	}

	return &drawresultservicev1.GetDrawsListResponse{Draws: respDraws}, nil
}

// GetDrawResult - получение результатов тиража
func (s *Server) GetDrawResult(ctx context.Context, req *drawresultservicev1.GetDrawResultRequest) (*drawresultservicev1.GetDrawResultResponse, error) {
	draw, err := s.usecase.GetDrawResult(ctx, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("get draw result: %w", err)
	}

	return &drawresultservicev1.GetDrawResultResponse{
		Id:                 draw.ID,
		DrawId:             draw.DrawID,
		ResultTime:         timestamppb.New(draw.ResultTime),
		WinningCombination: draw.WinningCombination,
	}, nil
}
