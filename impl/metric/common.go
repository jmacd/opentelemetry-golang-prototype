package metric

import (
	"github.com/lightstep/sandbox/jmacd/otel/core"
	"github.com/lightstep/sandbox/jmacd/otel/observer"
	"github.com/lightstep/sandbox/jmacd/otel/scope"
	"github.com/lightstep/sandbox/jmacd/otel/tag"
)

type (
	baseMetric struct {
		measure core.Measure

		mtype   MetricType
		keys    []core.Key
		eventID core.EventID
		status  error // Indicates registry conflict
	}

	baseEntry struct {
		base    *baseMetric
		metric  Metric
		eventID core.EventID
	}
)

func initBaseMetric(name string, mtype MetricType, opts []Options, init Metric) Metric {
	var tagOpts []tag.Options
	bm := init.base()

	for _, opt := range opts {
		opt(bm, &tagOpts)
	}

	bm.measure = tag.NewMeasure(name, tagOpts...)
	bm.mtype = mtype

	bm.eventID = observer.Record(observer.Event{
		Type:  observer.NEW_METRIC,
		Scope: bm.measure.DefinitionID().Scope(),
	})

	other, err := GetRegistry().RegisterMetric(init)
	if err != nil {
		bm.status = err
	}
	return other
}

func (bm *baseMetric) base() *baseMetric {
	return bm
}

func (bm *baseMetric) Measure() core.Measure {
	return bm.measure
}

func (bm *baseMetric) Type() MetricType {
	return bm.mtype
}

func (bm *baseMetric) Fields() []core.Key {
	return bm.keys
}

func (bm *baseMetric) Err() error {
	return bm.status
}

func (e *baseEntry) init(m Metric, values []core.KeyValue) {
	e.base = m.base()
	e.metric = m
	e.eventID = scope.Start(core.ScopeID{}, values...).ScopeID().EventID
}
