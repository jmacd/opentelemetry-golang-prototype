package metric

import (
	"github.com/lightstep/sandbox/jmacd/otel/core"
	"github.com/lightstep/sandbox/jmacd/otel/tag"
	"github.com/lightstep/sandbox/jmacd/otel/unit"
)

type (
	Metric interface {
		Measure() core.Measure

		DefinitionID() core.EventID

		Type() MetricType
		Fields() []core.Key
		Err() error

		base() *baseMetric
	}

	MetricType int
)

const (
	Invalid MetricType = iota
	GaugeInt64
	GaugeFloat64
	DerivedGaugeInt64
	DerivedGaugeFloat64
	CumulativeInt64
	CumulativeFloat64
	DerivedCumulativeInt64
	DerivedCumulativeFloat64
)

type (
	Options func(*baseMetric, *[]tag.Options)
)

// WithDescription applies provided description.
func WithDescription(desc string) Options {
	return func(_ *baseMetric, to *[]tag.Options) {
		*to = append(*to, tag.WithDescription(desc))
	}
}

// WithUnit applies provided unit.
func WithUnit(unit unit.Unit) Options {
	return func(_ *baseMetric, to *[]tag.Options) {
		*to = append(*to, tag.WithUnit(unit))
	}
}

// WithKeys applies the provided dimension keys.
func WithKeys(keys ...core.Key) Options {
	return func(bm *baseMetric, _ *[]tag.Options) {
		bm.keys = keys
	}
}
