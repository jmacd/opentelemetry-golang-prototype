package spandata

import (
	"io"

	"github.com/lightstep/opentelemetry-golang-prototype/exporter/observer"
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/reader"
)

type (
	Reader interface {
		Read(Span)
	}

	Span struct {
	}

	spanReader struct {
	}
)

func NewReaderObserver(w io.Writer) observer.Observer {
	return reader.NewReaderObserver(&spanReader{})
}

func (s *spanReader) Read(data reader.Event) {
	switch data.Type {
	case reader.START_SPAN:
	case reader.FINISH_SPAN:
	case reader.LOG_EVENT:
	case reader.LOGF_EVENT:
	case reader.MODIFY_ATTR:
	case reader.RECORD_STATS:
	}
}
