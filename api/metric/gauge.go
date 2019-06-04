package metric

import (
	"context"

	"github.com/lightstep/opentelemetry-golang-prototype/api/core"
	"github.com/lightstep/opentelemetry-golang-prototype/api/trace"
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/observer"
)

type (
	Float64Gauge struct {
		baseMetric
	}

	Float64Entry struct {
		baseEntry
	}
)

func NewFloat64Gauge(name string, mos ...Option) *Float64Gauge {
	m := initBaseMetric(name, GaugeFloat64, mos, &Float64Gauge{}).(*Float64Gauge)
	return m
}

func (g *Float64Gauge) Gauge(values ...core.KeyValue) Float64Entry {
	var entry Float64Entry
	entry.init(g, values)
	return entry
}

func (g *Float64Gauge) DefinitionID() core.EventID {
	return g.eventID
}

func (g Float64Entry) Set(ctx context.Context, val float64) {
	observer.Record(observer.Event{
		Type: observer.SET_GAUGE,
		Scope: core.ScopeID{
			EventID:     g.eventID,
			SpanContext: trace.Active(ctx).ScopeID().SpanContext,
		},
		Metric:  g.metric.DefinitionID(),
		Context: ctx,
		Float64: val,
	})
}

func (g Float64Entry) Add(ctx context.Context, val float64) {
	observer.Record(observer.Event{
		Type:    observer.ADD_GAUGE,
		Scope:   g.eventID.Scope(),
		Metric:  g.metric.DefinitionID(),
		Context: ctx,
		Float64: val,
	})
}
