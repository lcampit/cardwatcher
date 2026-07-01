// Package middleware exposes observability
// interceptors used in the server
package middleware

import (
	"context"

	"github.com/lcampit/cardwatcher/apps/server/internal/observability/logger"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const correlationHeader = "x-correlation-id"

func CorrelationInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		var corrID string

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get(correlationHeader); len(vals) > 0 {
				corrID = vals[0]
			}
		}

		if corrID == "" {
			corrID = uuid.NewString()
		}

		ctx = logger.WithCorrelationID(ctx, logger.CorrelationID(corrID))

		return handler(ctx, req)
	}
}
