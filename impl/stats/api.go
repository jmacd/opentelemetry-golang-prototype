package stats

import (
	"context"

	"github.com/lightstep/sandbox/jmacd/otel/core"
	"github.com/lightstep/sandbox/jmacd/otel/scope"
)

type (
	Recorder interface {
		Record(ctx context.Context, m ...core.Measurement)
	}
)

func With(scope core.Scope) Recorder {
	return scopeRecorder{scope}
}

func Record(ctx context.Context, m ...core.Measurement) {
	With(scope.Active(ctx)).Record(ctx, m...)
}
