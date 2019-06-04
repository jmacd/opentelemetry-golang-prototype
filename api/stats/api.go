package stats

import (
	"context"

	"github.com/lightstep/opentelemetry-golang-prototype/api/core"
	"github.com/lightstep/opentelemetry-golang-prototype/api/scope"
)

type (
	Interface interface {
		Record(ctx context.Context, m ...core.Measurement)
	}
)

func With(scope scope.Scope) Recorder {
	return Recorder{scope}
}

func Record(ctx context.Context, m ...core.Measurement) {
	With(scope.Active(ctx)).Record(ctx, m...)
}
