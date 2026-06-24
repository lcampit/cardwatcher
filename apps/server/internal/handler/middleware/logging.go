package middleware

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/lcampit/cardwatcher/apps/server/internal/logger"
)

func LoggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		service, method := splitFullMethodName(info.FullMethod)
		logger.Info("server method invoked",
			slog.String("service", service),
			slog.String("method", method))

		logger.Debug("request details", slog.Any("request", req))
		// Call handler
		resp, err := handler(ctx, req)

		// Log request
		duration := time.Since(start).Milliseconds()
		statusCode := "OK"
		if err != nil {
			statusCode = status.Code(err).String()
			logger.Error("error in method",
				slog.String("service", service),
				slog.String("method", method),
				slog.Any("error", err))
		} else {
			logger.Debug("response details", slog.Any("response", resp))
		}

		logger.Info("response returned",
			slog.Time("startTime", start),
			slog.String("service", service),
			slog.String("method", method),
			slog.String("statusCode", statusCode),
			slog.Int64("durationMs", duration))

		return resp, err
	}
}

// splitFullMethodName splits the service and method portions of the
// fullMethod written as /service/method
func splitFullMethodName(fullMethod string) (string, string) {
	fullMethod = strings.TrimPrefix(fullMethod, "/") // remove leading slash
	if before, after, ok := strings.Cut(fullMethod, "/"); ok {
		return before, after
	}
	return "unknown", "unknown"
}

func CorrelationIDInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		id := uuid.NewString()
		ctx = logger.WithCorrelationID(ctx, id)
		return handler(ctx, req)
	}
}
