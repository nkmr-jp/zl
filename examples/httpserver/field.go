package main

import (
	"fmt"
	"go.uber.org/zap"
	"time"
)

const (
	TraceIDFieldKey  = "trace_id"
	DurationFieldKey = "duration"
)

// DurationField returns a zap field with the given key and a human-readable
func DurationField(start time.Time) zap.Field {
	return zap.String(
		DurationFieldKey,
		fmt.Sprintf("%.4fs", float64(time.Since(start))/float64(time.Second)),
	)
}

// TraceIDField returns a zap field with the given key and value.
func TraceIDField(traceID string) zap.Field {
	return zap.String(
		TraceIDFieldKey,
		traceID,
	)
}
