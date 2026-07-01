package logger

import (
	"context"
	"log/slog"
)

type ContextHandler struct {
	inner slog.Handler
}

// Handle gets called for every log and recovers correlation ID from context
func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if id, ok := CorrelationIDFromContext(ctx); ok {
		r.AddAttrs(slog.String(KeyCorrelationID, string(id)))
	}

	return h.inner.Handle(ctx, r)
}

func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{
		inner: h.inner.WithAttrs(attrs),
	}
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{
		inner: h.inner.WithGroup(name),
	}
}
