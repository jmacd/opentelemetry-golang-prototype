package trace

import (
	"context"

	"github.com/lightstep/opentelemetry-golang-prototype/api/core"
	"github.com/lightstep/opentelemetry-golang-prototype/api/log"
	"github.com/lightstep/opentelemetry-golang-prototype/api/observer"
	"github.com/lightstep/opentelemetry-golang-prototype/api/stats"
)

func (sp *span) ScopeID() core.ScopeID {
	if sp == nil {
		return core.ScopeID{}
	}
	sp.lock.Lock()
	sid := core.ScopeID{
		EventID:     sp.eventID,
		SpanContext: sp.spanContext,
	}
	sp.lock.Unlock()
	return sid
}

func (sp *span) updateScope() (core.ScopeID, core.EventID) {
	next := observer.NextEventID()

	sp.lock.Lock()
	sid := core.ScopeID{
		EventID:     sp.eventID,
		SpanContext: sp.spanContext,
	}
	sp.eventID = next
	sp.lock.Unlock()

	return sid, next
}

func (sp *span) SetError(v bool) {
	sp.SetAttribute(ErrorKey.Bool(v))
}

func (sp *span) SetAttribute(attribute core.KeyValue) {
	if sp == nil {
		return
	}

	sid, next := sp.updateScope()

	observer.Record(observer.Event{
		Type:      observer.MODIFY_ATTR,
		Scope:     sid,
		Sequence:  next,
		Attribute: attribute,
	})
}

func (sp *span) SetAttributes(attributes ...core.KeyValue) {
	if sp == nil {
		return
	}

	sid, next := sp.updateScope()

	observer.Record(observer.Event{
		Type:       observer.MODIFY_ATTR,
		Scope:      sid,
		Sequence:   next,
		Attributes: attributes,
	})
}

func (sp *span) ModifyAttribute(mutator core.Mutator) {
	if sp == nil {
		return
	}

	sid, next := sp.updateScope()

	observer.Record(observer.Event{
		Type:     observer.MODIFY_ATTR,
		Scope:    sid,
		Sequence: next,
		Mutator:  mutator,
	})
}

func (sp *span) ModifyAttributes(mutators ...core.Mutator) {
	if sp == nil {
		return
	}

	sid, next := sp.updateScope()

	observer.Record(observer.Event{
		Type:     observer.MODIFY_ATTR,
		Scope:    sid,
		Sequence: next,
		Mutators: mutators,
	})
}

func (sp *span) Finish() {
	if sp == nil {
		return
	}
	recovered := recover()
	sp.finishOnce.Do(func() {
		observer.Record(observer.Event{
			Type:      observer.FINISH_SPAN,
			Scope:     sp.ScopeID(),
			Recovered: recovered,
		})
	})
	if recovered != nil {
		panic(recovered)
	}
}

func (sp *span) Tracer() Tracer {
	return sp.tracer
}

func (sp *span) Log(ctx context.Context, msg string, args ...core.KeyValue) {
	log.With(sp).Log(ctx, msg, args...)
}

func (sp *span) Logf(ctx context.Context, fmt string, args ...interface{}) {
	log.With(sp).Logf(ctx, fmt, args...)
}

func (sp *span) Record(ctx context.Context, args ...core.Measurement) {
	stats.With(sp).Record(ctx, args...)
}
