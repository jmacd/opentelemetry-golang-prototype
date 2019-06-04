package stats

import (
	"context"

	"github.com/lightstep/sandbox/jmacd/otel/core"
	"github.com/lightstep/sandbox/jmacd/otel/observer"
)

type (
	scopeRecorder struct {
		scope core.Scope
	}

	statsKey struct {
		key core.Measure
	}
)

func (s scopeRecorder) Record(ctx context.Context, m ...core.Measurement) {
	observer.Record(observer.Event{
		Type:    observer.RECORD_STATS,
		Scope:   s.scope.ScopeID(),
		Context: ctx,
		Stats:   m,
	})
}
