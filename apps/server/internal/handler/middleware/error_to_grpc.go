// Package middleware provides interceptors
// to wrap all requests on the server
package middleware

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	apperrors "github.com/lcampit/cardwatcher/apps/server/internal/errors"
)

// ErrorInterceptor translates application errors into gRPC statuses.
func ErrorInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	resp, err := handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	if appErr, ok := errors.AsType[*apperrors.AppError](err); ok {
		return nil, appErr.ToGRPCStatus().Err()
	}

	if _, ok := status.FromError(err); ok {
		return nil, err // Already a gRPC status
	}

	return nil, apperrors.NewInternal("An unexpected error occurred", "", err).ToGRPCStatus().Err()
}
