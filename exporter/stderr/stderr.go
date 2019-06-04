package stderr

import (
	"fmt"
	"strings"

	"github.com/lightstep/sandbox/jmacd/otel/core"
	"github.com/lightstep/sandbox/jmacd/otel/observer"
	"github.com/lightstep/sandbox/jmacd/otel/observer/reader"
	"github.com/lightstep/sandbox/jmacd/otel/trace"
)

type (
	debugLog struct{}
)

var (
	logger = newDebugLog()
)

func newDebugLog() *debugLog {
	logger := &debugLog{}
	observer.RegisterObserver(reader.NewReaderObserver(logger))
	return logger
}

func (l *debugLog) Read(data reader.Event) {
	var buf strings.Builder

	f := func(skipIf bool) func(kv core.KeyValue) bool {
		return func(kv core.KeyValue) bool {
			if skipIf && data.Attributes.HasValue(kv.Key) {
				return true
			}
			buf.WriteString(" " + kv.Key.Name() + "=" + kv.Value.Emit())
			return true
		}
	}

	buf.WriteString(fmt.Sprint(data.Sequence, ": "))
	buf.WriteString(data.Time.Format("2006/01/02 15-04-05.000000"))
	buf.WriteString(" ")

	switch data.Type {
	case reader.START_SPAN:
		buf.WriteString("start ")
		buf.WriteString(data.Name)

		if !data.Parent.HasSpanID() {
			buf.WriteString(", a root span")
		} else {
			buf.WriteString(" <")
			if data.Parent.HasSpanID() {
				f(false)(trace.ParentSpanIDKey.String(data.SpanContext.SpanIDString()))
			}
			data.ParentAttributes.Foreach(f(false))
			buf.WriteString(" >")
		}

	case reader.FINISH_SPAN:
		buf.WriteString("finish ")
		buf.WriteString(data.Name)

		buf.WriteString(" (")
		buf.WriteString(data.Duration.String())
		buf.WriteString(")")

	case reader.LOG_EVENT:
		buf.WriteString(data.Message)

	case reader.LOGF_EVENT:
		buf.WriteString(data.Message)

	case reader.SET_GAUGE:
		buf.WriteString("set gauge ")
		buf.WriteString(data.Name)
	case reader.ADD_GAUGE:
		buf.WriteString("add gauge")
		buf.WriteString(data.Name)
	case reader.MODIFY_ATTR:
		buf.WriteString("modify attr")
	case reader.RECORD_STATS:
		buf.WriteString("record")

		buf.WriteString(" <")
		for _, s := range data.Stats {
			f(false)(s.Measure.V(s.Value))
		}
		buf.WriteString(" >")
	default:
		buf.WriteString(fmt.Sprintf("WAT? %d", data.Type))
	}

	// Attach the scope (span) attributes and context tags.
	buf.WriteString(" [")
	data.Attributes.Foreach(f(false))
	data.Tags.Foreach(f(true))
	if data.SpanContext.HasSpanID() {
		f(false)(trace.SpanIDKey.String(data.SpanContext.SpanIDString()))
	}
	if data.SpanContext.HasTraceID() {
		f(false)(trace.TraceIDKey.String(data.SpanContext.TraceIDString()))
	}

	buf.WriteString(" ]")

	fmt.Println(buf.String())
}
