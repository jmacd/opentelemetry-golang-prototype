package stats

import (
	"context"

	"github.com/lightstep/opentelemetry-golang-prototype/api/core"
	"github.com/lightstep/opentelemetry-golang-prototype/api/scope"
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/observer"
)

type (
	Recorder struct {
		scope.Scope
	}

	statsKey struct {
		key core.Measure
	}
)

func (s Recorder) Record(ctx context.Context, m ...core.Measurement) {
	observer.Record(observer.Event{
		Type:    observer.RECORD_STATS,
		Scope:   s.ScopeID(),
		Context: ctx,
		Stats:   m,
	})
}
