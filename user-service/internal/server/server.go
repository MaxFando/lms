package server

import (
	"context"
	"net"
	"time"

	"github.com/MaxFando/lms/platform/logger"
	userservicev1 "github.com/MaxFando/lms/user-service/api/grpc/gen/go/user-service/v1"
	"github.com/MaxFando/lms/user-service/internal/server/interceptor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	defaultGRPCPort         = "50051"
	defaultKeepAliveTime    = 30 * time.Second
	defaultKeepAliveTimeout = 20 * time.Second
	defaultMaxRecvMsgSize   = 50 << 20
	defaultMaxSendMsgSize   = 50 << 20
)

type Server struct {
	grpcServer *grpc.Server
	logger     logger.Logger
	errors     chan error
	grpcPort   string
}

func NewServer(logger logger.Logger, svcServer userservicev1.UserServiceServer) *Server {
	s := &Server{
		logger:   logger,
		errors:   make(chan error, 1),
		grpcPort: defaultGRPCPort,
	}
	s.grpcServer = initGRPCServer(logger, svcServer)
	return s
}

func initGRPCServer(logger logger.Logger, svcServer userservicev1.UserServiceServer) *grpc.Server {
	grpcSrv := grpc.NewServer(
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

	userservicev1.RegisterUserServiceServer(grpcSrv, svcServer)
	reflection.Register(grpcSrv)
	return grpcSrv
}

func (s *Server) Serve(ctx context.Context) {
	lis, err := net.Listen("tcp", ":"+s.grpcPort)
	if err != nil {
		s.sendError(ctx, err)
		return
	}
	s.logger.Info(ctx, "gRPC server listening on port "+s.grpcPort)
	if err := s.grpcServer.Serve(lis); err != nil {
		s.sendError(ctx, err)
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	s.logger.Info(ctx, "Shutting down gRPC server")
	s.grpcServer.GracefulStop()
	close(s.errors)
}

func (s *Server) Notify() <-chan error {
	out := make(chan error)
	go func() {
		for err := range s.errors {
			out <- err
		}
		close(out)
	}()
	return out
}

func (s *Server) sendError(ctx context.Context, err error) {
	select {
	case s.errors <- err:
	default:
		s.logger.Error(ctx, "error channel full, dropping", "error", err)
	}
}
