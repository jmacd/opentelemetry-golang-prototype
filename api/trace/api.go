package trace

import (
	"context"

	"github.com/lightstep/opentelemetry-golang-prototype/api/core"
	"github.com/lightstep/opentelemetry-golang-prototype/api/log"
	"github.com/lightstep/opentelemetry-golang-prototype/api/scope"
	"github.com/lightstep/opentelemetry-golang-prototype/api/stats"
	"github.com/lightstep/opentelemetry-golang-prototype/api/tag"
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

		// Note: see https://github.com/opentracing/opentracing-go/issues/127
		Inject(context.Context, Span, Injector)

		// ScopeID returns the resource scope of this tracer.
		scope.Scope
	}

	Span interface {
		scope.Mutable

		log.Interface

		stats.Interface

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
