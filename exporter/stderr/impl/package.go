package impl

import (
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/observer"
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/stderr"
)

func init() {
	observer.RegisterObserver(stderr.New())
}
