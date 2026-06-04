package middleware

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	apperrors "github.com/lcampit/cardwatcher/apps/server/internal/errors"
)

// UnaryErrorInterceptor translates application errors into gRPC statuses.
func UnaryErrorInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	resp, err := handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	if appErr, ok := errors.AsType[*apperrors.AppError](err); ok {
		if appErr.CausedBy != nil {
			log.Printf("ERROR: %s, Original cause: %v", appErr.Title, appErr.CausedBy)
		}
		return nil, appErr.ToGRPCStatus().Err()
	}

	if _, ok := status.FromError(err); ok {
		return nil, err // Already a gRPC status
	}

	log.Printf("UNEXPECTED ERROR: %v", err)
	return nil, apperrors.NewInternal("An unexpected error occurred", "", err).ToGRPCStatus().Err()
}
