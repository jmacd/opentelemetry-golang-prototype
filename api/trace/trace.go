package trace

import (
	"context"
	"math/rand"
	"sync"
	"sync/atomic"

	"github.com/lightstep/opentelemetry-golang-prototype/api/core"
	"github.com/lightstep/opentelemetry-golang-prototype/api/log"
	"github.com/lightstep/opentelemetry-golang-prototype/api/observer"
	"github.com/lightstep/opentelemetry-golang-prototype/api/scope"
	"github.com/lightstep/opentelemetry-golang-prototype/api/tag"
)

type (
	span struct {
		tracer      *tracer
		spanContext core.SpanContext
		lock        sync.Mutex
		eventID     core.EventID
		finishOnce  sync.Once
	}

	tracer struct {
		resources core.EventID
	}
)

var (
	ServiceKey      = tag.New("service")
	ComponentKey    = tag.New("component")
	ErrorKey        = tag.New("error")
	SpanIDKey       = tag.New("span_id")
	TraceIDKey      = tag.New("trace_id")
	ParentSpanIDKey = tag.New("parent_span_id")
	MessageKey      = tag.New("message",
		tag.WithDescription("message text: info, error, etc"),
	)

	// The process global tracer could have process-wide resource
	// tags applied directly, or we can have a SetGlobal tracer to
	// install a default tracer w/ resources.
	global atomic.Value
	empty  = &tracer{}
)

func (t *tracer) ScopeID() core.ScopeID {
	return t.resources.Scope()
}

func (t *tracer) WithResources(attributes ...core.KeyValue) Tracer {
	s := scope.Start(t.resources.Scope(), attributes...)
	return &tracer{
		resources: s.ScopeID().EventID,
	}
}

func (g *tracer) WithComponent(name string) Tracer {
	return g.WithResources(ComponentKey.String(name))
}

func (g *tracer) WithService(name string) Tracer {
	return g.WithResources(ServiceKey.String(name))
}

func (t *tracer) WithSpan(ctx context.Context, name string, body func(context.Context) error) error {
	// TODO: use runtime/trace.WithRegion for execution tracer support
	ctx, span := t.Start(ctx, name)
	defer span.Finish()

	if err := body(ctx); err != nil {
		span.SetAttribute(ErrorKey.Bool(true))
		log.Log(ctx, "span error", MessageKey.String(err.Error()))
		return err
	}
	return nil
}

// TODO options, should be Options
func (t *tracer) Start(ctx context.Context, name string, attributes ...core.KeyValue) (context.Context, Span) {
	var child core.SpanContext

	parentScope := Active(ctx).ScopeID()

	child.SpanID = rand.Uint64()

	if parentScope.HasTraceID() {
		parent := parentScope.SpanContext
		child.TraceIDHigh = parent.TraceIDHigh
		child.TraceIDLow = parent.TraceIDLow
	} else {
		child.TraceIDHigh = rand.Uint64()
		child.TraceIDLow = rand.Uint64()
	}

	childScope := core.ScopeID{
		SpanContext: child,
		EventID:     t.resources,
	}

	span := &span{
		spanContext: child,
		tracer:      t,
		eventID: observer.Record(observer.Event{
			Type:    observer.START_SPAN,
			Scope:   scope.Start(childScope, attributes...).ScopeID(),
			Context: ctx,
			Parent:  parentScope,
			String:  name,
		}),
	}
	return scope.SetActive(ctx, span), span
}

func (t *tracer) Inject(ctx context.Context, span Span, injector Injector) {
	injector.Inject(span.ScopeID().SpanContext, tag.FromContext(ctx))
}
