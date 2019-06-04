package log

import (
	"context"

	"github.com/lightstep/sandbox/jmacd/otel/core"
	"github.com/lightstep/sandbox/jmacd/otel/observer"
)

func (l Logger) Log(ctx context.Context, msg string, fields ...core.KeyValue) {
	observer.Record(observer.Event{
		Type:       observer.LOG_EVENT,
		Scope:      l.scope.ScopeID(),
		String:     msg,
		Attributes: fields,
		Context:    ctx,
	})
}

func (l Logger) Logf(ctx context.Context, fmt string, args ...interface{}) {
	observer.Record(observer.Event{
		Type:      observer.LOGF_EVENT,
		Scope:     l.scope.ScopeID(),
		String:    fmt,
		Arguments: args,
		Context:   ctx,
	})
}
