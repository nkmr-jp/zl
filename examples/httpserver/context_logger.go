package main

import (
	"context"
	"github.com/nkmr-jp/zl"
)

const (
	loggerKey = "logger"
)

var (
	ctxKeyLogger = &contextKey{loggerKey}
)

type contextKey struct {
	name string
}

// GetLogger retrieves the logger from the context or creates a new logger if it does not exist.
func GetLogger(ctx context.Context) *zl.Logger {
	if logger, ok := ctx.Value(ctxKeyLogger).(*zl.Logger); ok {
		return logger
	}
	return zl.New()
}

// SetNewLogger sets the logger and trace ID in context.
// Set a trace ID so that the processing of this request can be traced.
func SetNewLogger(ctx context.Context, traceID string) context.Context {
	logger := zl.New(TraceIDField(traceID))
	return context.WithValue(ctx, ctxKeyLogger, logger)
}
