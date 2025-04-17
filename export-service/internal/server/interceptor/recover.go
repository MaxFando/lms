package interceptor

import (
	"context"
	"fmt"

	"github.com/MaxFando/lms/platform/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func PanicRecoveryUnaryInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic: %v", r)

				logger.With("method", info.FullMethod).Error(ctx, "panic recovered", "error", err)

				err = status.Errorf(codes.Internal, "internal server error: %v", r)
			}
		}()

		return handler(ctx, req)
	}
}
