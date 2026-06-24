package middleware

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/lcampit/cardwatcher/apps/server/internal/observability/logger"
)

func ObservabilityInterceptor(
	log *slog.Logger,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start).Milliseconds()
		st := status.Code(err)

		if err != nil {
			log.ErrorContext(
				ctx, "grpc error",
				slog.String(logger.KeyMethod, info.FullMethod),
				slog.String(logger.KeyStatus, st.String()),
				slog.Int64(logger.KeyDuration, duration),
				logger.Err(err),
			)
		} else {
			log.InfoContext(
				ctx, "grpc response",
				slog.String(logger.KeyMethod, info.FullMethod),
				slog.String(logger.KeyStatus, st.String()),
				slog.Int64(logger.KeyDuration, duration),
			)
		}

		return resp, err
	}
}
