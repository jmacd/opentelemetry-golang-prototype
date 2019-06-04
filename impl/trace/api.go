package trace

import (
	"context"

	"github.com/lightstep/sandbox/jmacd/otel/core"
	"github.com/lightstep/sandbox/jmacd/otel/log"
	"github.com/lightstep/sandbox/jmacd/otel/scope"
	"github.com/lightstep/sandbox/jmacd/otel/stats"
	"github.com/lightstep/sandbox/jmacd/otel/tag"
)

type (
	Tracer interface {
		Start(context.Context, string, ...core.KeyValue) (context.Context, Span)

		WithSpan(
			ctx context.Context,
			operation string,
			body func(ctx context.Context) error,
		) error

		WithService(name string) Tracer
		WithComponent(name string) Tracer
		WithResources(res ...core.KeyValue) Tracer

		// TODO: see https://github.com/opentracing/opentracing-go/issues/127
		Inject(context.Context, Span, Injector)

		// ScopeID returns the resource scope of this tracer.
		core.Scope
	}

	Span interface {
		scope.Mutable

		log.Logger

		stats.Recorder

		SetError(bool)

		Tracer() Tracer

		Finish()
	}

	Injector interface {
		Inject(core.SpanContext, tag.Map)
	}
)

func GlobalTracer() Tracer {
	if t := global.Load(); t != nil {
		return t.(Tracer)
	}
	return empty
}

func SetGlobalTracer(t Tracer) {
	global.Store(t)
}

func Start(ctx context.Context, name string, attrs ...core.KeyValue) (context.Context, Span) {
	return GlobalTracer().Start(ctx, name, attrs...)
}

func Active(ctx context.Context) Span {
	span, _ := scope.Active(ctx).(*span)
	// TODO make a real no-op
	return span
}

func WithSpan(ctx context.Context, name string, body func(context.Context) error) error {
	return GlobalTracer().WithSpan(ctx, name, body)
}

func SetError(ctx context.Context, v bool) {
	Active(ctx).SetError(v)
}

func Inject(ctx context.Context, injector Injector) {
	span := Active(ctx)
	if span == nil {
		return
	}

	span.Tracer().Inject(ctx, span, injector)
}
