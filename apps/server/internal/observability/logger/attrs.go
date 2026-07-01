package logger

import (
	"log/slog"
)

// Err is a custom slog attr to make sure all errors
// in the application are logger using the right key
func Err(err error) slog.Attr {
	if err == nil {
		return slog.String(KeyError, "")
	}

	attrs := []slog.Attr{
		slog.String("message", err.Error()),
	}

	return slog.GroupAttrs(KeyError, attrs...)
}
