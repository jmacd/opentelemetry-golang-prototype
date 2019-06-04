package main

import (
	"context"

	"github.com/lightstep/opentelemetry-golang-prototype/api/log"
	"github.com/lightstep/opentelemetry-golang-prototype/api/tag"
	"github.com/lightstep/opentelemetry-golang-prototype/impl/metric"
	"github.com/lightstep/opentelemetry-golang-prototype/impl/stats"
	"github.com/lightstep/opentelemetry-golang-prototype/impl/trace"

	// This creates a debug log on the console.
	_ "github.com/lightstep/opentelemetry-golang-prototype/exporter/stderr"
)

var (
	tracer = trace.GlobalTracer().
		WithComponent("example").
		WithResources(
			tag.New("whatevs").String("yesss"),
		)

	fooKey     = tag.New("ex.com/foo", tag.WithDescription("A Foo var"))
	barKey     = tag.New("ex.com/bar", tag.WithDescription("A Bar var"))
	lemonsKey  = tag.New("ex.com/lemons", tag.WithDescription("A Lemons var"))
	anotherKey = tag.New("ex.com/another")

	oneMetric = metric.NewFloat64Gauge("ex.com/one",
		metric.WithKeys(fooKey, barKey, lemonsKey),
		metric.WithDescription("A gauge set to 1.0"),
	)

	measureTwo = tag.NewMeasure("ex.com/two")
)

func main() {
	ctx := context.Background()

	ctx = tag.NewContext(ctx,
		tag.Insert(fooKey.String("foo1")),
		tag.Insert(barKey.String("bar1")),
	)

	gauge := oneMetric.Gauge(
		fooKey.Value(ctx),
		barKey.Value(ctx),
		lemonsKey.Int(10),
	)

	err := tracer.WithSpan(ctx, "operation", func(ctx context.Context) error {

		trace.SetError(ctx, true)

		log.Log(ctx, "Nice operation!", tag.New("bogons").Int(100))

		trace.Active(ctx).SetAttributes(anotherKey.String("yes"))

		gauge.Set(ctx, 1)

		return tracer.WithSpan(
			ctx,
			"Sub operation...",
			func(ctx context.Context) error {
				trace.Active(ctx).SetAttribute(lemonsKey.String("five"))

				log.Logf(ctx, "Format schmormat %d!", 100)

				stats.Record(ctx, measureTwo.M(1.3))

				return nil
			},
		)
	})
	if err != nil {
		panic(err)
	}
}
