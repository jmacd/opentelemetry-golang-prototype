package log

import (
	"context"

	"github.com/lightstep/opentelemetry-golang-prototype/api/core"
	"github.com/lightstep/opentelemetry-golang-prototype/api/scope"
)

type (
	Interface interface {
		Log(ctx context.Context, msg string, fields ...core.KeyValue)
		Logf(ctx context.Context, fmt string, args ...interface{})
	}

	Logger struct {
		scope.Scope
	}
)

func With(scope scope.Scope) Logger {
	return Logger{scope}
}

func Log(ctx context.Context, msg string, fields ...core.KeyValue) {
	With(scope.Active(ctx)).Log(ctx, msg, fields...)
}

func Logf(ctx context.Context, fmt string, args ...interface{}) {
	With(scope.Active(ctx)).Logf(ctx, fmt, args...)
}

func (l Logger) Log(ctx context.Context, msg string, fields ...core.KeyValue) {
}
func (l Logger) Logf(ctx context.Context, fmt string, args ...interface{}) {
}
