package writer

import (
	"io"

	"github.com/lightstep/opentelemetry-golang-prototype/exporter/observer"
)

type (
	Encoder struct {
		file io.Writer
	}
)

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		file: w,
	}
}

func (e *Encoder) Observe(data observer.Event) {

}
