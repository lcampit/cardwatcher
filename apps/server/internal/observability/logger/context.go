package logger

import "context"

type CorrelationID string

func WithCorrelationID(ctx context.Context, id CorrelationID) context.Context {
	return context.WithValue(ctx, correlationIDKey, id)
}

func CorrelationIDFromContext(ctx context.Context) (CorrelationID, bool) {
	id, ok := ctx.Value(correlationIDKey).(CorrelationID)
	return id, ok
}
