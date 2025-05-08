package v1

import (
	"context"
	"fmt"
	"time"

	drawresultservicev1 "github.com/MaxFando/lms/draw-service/api/grpc/gen/go/draw-service/v1"
	"github.com/MaxFando/lms/draw-service/internal/entity"
)

type DrawServiceUsecase interface {
	CreateDraws(ctx context.Context, draw entity.Draw) (*entity.Draw, error)
	GetDrawsList(ctx context.Context) ([]*entity.Draw, error)
	CancelDraw(ctx context.Context, id int32) error
	GetCompletedDraws(ctx context.Context) ([]*entity.Draw, error)
	GetDrawResult(ctx context.Context, id int32) (*entity.DrawResult, error)
}

type Server struct {
	drawresultservicev1.UnimplementedDrawServiceServer
	usecase DrawServiceUsecase
}

func NewServer(usecase DrawServiceUsecase) *Server {
	return &Server{usecase: usecase}
}

// CreateDraws создает новый тираж
func (s *Server) CreateDraw(ctx context.Context, req *drawresultservicev1.CreateDrawRequest) (*drawresultservicev1.DrawResponse, error) {
	startTime, err := time.Parse(time.RFC3339, req.GetStartTime())
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	endTime, err := time.Parse(time.RFC3339, req.GetEndTime())
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	draw := entity.Draw{
		LotteryType: entity.LotteryType(req.GetLotteryType()),
		StartTime:   startTime,
		EndTime:     endTime,
	}

	createdDraw, err := s.usecase.CreateDraws(ctx, draw)
	if err != nil {
		return nil, fmt.Errorf("create draws: %w", err)
	}

	return &drawresultservicev1.DrawResponse{
		Id:          createdDraw.ID,
		LotteryType: string(createdDraw.LotteryType),
		StartTime:   createdDraw.StartTime.Format(time.RFC3339),
		EndTime:     createdDraw.EndTime.Format(time.RFC3339),
		Status:      string(createdDraw.Status),
	}, nil
}

// GetDrawsList возвращает список активных тиражей
func (s *Server) GetDrawsList(ctx context.Context, req *drawresultservicev1.GetDrawsListRequest) (*drawresultservicev1.GetDrawsListResponse, error) {
	draws, err := s.usecase.GetDrawsList(ctx)
	if err != nil {
		return nil, fmt.Errorf("get draws list: %w", err)
	}

	var respDraws []*drawresultservicev1.DrawResponse
	for _, d := range draws {
		respDraws = append(respDraws, &drawresultservicev1.DrawResponse{
			Id:          d.ID,
			LotteryType: string(d.LotteryType),
			StartTime:   d.StartTime.Format(time.RFC3339),
			EndTime:     d.EndTime.Format(time.RFC3339),
			Status:      string(d.Status),
		})
	}

	return &drawresultservicev1.GetDrawsListResponse{Draws: respDraws}, nil
}

// CancelDraw отменяет тираж по ID
func (s *Server) CancelDraw(ctx context.Context, req *drawresultservicev1.CancelDrawRequest) (*drawresultservicev1.CancelDrawResponse, error) {
	err := s.usecase.CancelDraw(ctx, int32(req.GetId()))
	if err != nil {
		return nil, fmt.Errorf("cancel draw: %w", err)
	}

	return &drawresultservicev1.CancelDrawResponse{}, nil
}

// GetCompletedDrawsList - получения списка завершенных тиражей
func (s *Server) GetCompletedDrawsList(ctx context.Context, req *drawresultservicev1.GetDrawsListRequest) (*drawresultservicev1.GetDrawsListResponse, error) {
	draws, err := s.usecase.GetCompletedDraws(ctx)
	if err != nil {
		return nil, fmt.Errorf("get completed draws list: %w", err)
	}

	var respDraws []*drawresultservicev1.DrawResponse
	for _, d := range draws {
		respDraws = append(respDraws, &drawresultservicev1.DrawResponse{
			Id:          d.ID,
			LotteryType: string(d.LotteryType),
			StartTime:   d.StartTime.Format(time.RFC3339),
			EndTime:     d.EndTime.Format(time.RFC3339),
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
		ResultTime:         draw.ResultTime.Format(time.RFC3339),
		WinningCombination: draw.WinningCombination,
	}, nil
}
