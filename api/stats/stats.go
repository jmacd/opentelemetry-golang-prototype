package stats

import (
	"context"

	"github.com/lightstep/opentelemetry-golang-prototype/api/core"
	"github.com/lightstep/opentelemetry-golang-prototype/api/scope"
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/observer"
)

type (
	Interface interface {
		Record(ctx context.Context, m ...core.Measurement)
	}

	Recorder struct {
		scope.Scope
	}
)

func With(scope scope.Scope) Recorder {
	return Recorder{scope}
}

func Record(ctx context.Context, m ...core.Measurement) {
	With(scope.Active(ctx)).Record(ctx, m...)
}

func (s Recorder) Record(ctx context.Context, m ...core.Measurement) {
	observer.Record(observer.Event{
		Type:    observer.RECORD_STATS,
		Scope:   s.ScopeID(),
		Context: ctx,
		Stats:   m,
	})
}
