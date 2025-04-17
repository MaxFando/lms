package server

import (
	"context"
	"net"
	"time"

	"github.com/MaxFando/lms/platform/logger"

	drawservicev1 "github.com/MaxFando/lms/draw-service/api/grpc/gen/go/draw-service/v1"
	"github.com/MaxFando/lms/draw-service/internal/server/interceptor"
	v1 "github.com/MaxFando/lms/draw-service/internal/server/service/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	defaultGRPCPort         = "50051"
	defaultKeepAliveTime    = 30 * time.Second
	defaultKeepAliveTimeout = 20 * time.Second

	defaultMaxRecvMsgSize = 1024 * 1024 * 50 // 50 MB
	defaultMaxSendMsgSize = 1024 * 1024 * 50 // 50 MB
)

type Server struct {
	grpcServer *grpc.Server
	logger     logger.Logger

	grpcPort string

	errors chan error
}

func NewServer(logger logger.Logger, serviceServer *v1.Server) *Server {
	srv := new(Server)

	srv.grpcPort = defaultGRPCPort
	srv.errors = make(chan error, 1)
	srv.logger = logger

	srv.grpcServer = initGRPCServer(logger, serviceServer)

	return srv
}

func (s *Server) Serve(ctx context.Context) {
	grpcDone := make(chan struct{})
	httpDone := make(chan struct{})

	go func() {
		defer close(grpcDone)
		s.serveGRPC(ctx)
	}()

	go func() {
		select {
		case <-ctx.Done():
			s.Shutdown(ctx)
		case <-grpcDone:
		case <-httpDone:
		}
	}()

	<-grpcDone
	<-httpDone
}

func (s *Server) Shutdown(ctx context.Context) {
	s.logger.Info(ctx, "Завершение работы сервера")

	s.grpcServer.GracefulStop()
	close(s.errors)
}

func (s *Server) Notify() <-chan error {
	errorsCh := make(chan error)

	go func() {
		for err := range s.errors {
			errorsCh <- err
		}
		close(errorsCh)
	}()

	return errorsCh
}

func (s *Server) serveGRPC(ctx context.Context) {
	grpcListener, err := net.Listen("tcp", ":"+s.grpcPort)
	if err != nil {
		s.sendError(ctx, err)

		return
	}

	s.logger.Info(ctx, "Запуск gRPC сервера на порту "+s.grpcPort)
	if err := s.grpcServer.Serve(grpcListener); err != nil {
		s.sendError(ctx, err)
	}
}

func initGRPCServer(logger logger.Logger, serviceServer *v1.Server) *grpc.Server {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.PanicRecoveryUnaryInterceptor(logger),
		),
		grpc.MaxRecvMsgSize(defaultMaxRecvMsgSize),
		grpc.MaxSendMsgSize(defaultMaxSendMsgSize),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    defaultKeepAliveTime,
			Timeout: defaultKeepAliveTimeout,
		}),
	)

	drawservicev1.RegisterDrawServiceServer(server, serviceServer)
	reflection.Register(server)

	return server
}

func (s *Server) sendError(ctx context.Context, err error) {
	select {
	case s.errors <- err:
		s.logger.Error(ctx, "Error sent to channel", "error", err)
	default:
		s.logger.Error(ctx, "Error channel is full, error dropped", "error", err)
	}
}
