package stdout

import (
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/buffer"
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/observer"
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/stderr"
)

func New() observer.Observer {
	return buffer.NewBuffer(1000, stderr.New())
}
