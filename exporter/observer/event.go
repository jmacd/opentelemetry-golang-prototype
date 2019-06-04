package observer

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lightstep/sandbox/jmacd/otel/core"
)

type (
	EventType int

	Event struct {
		// Automatic fields
		Sequence core.EventID // Auto-filled
		Time     time.Time    // Auto-filled

		// Type, Scope, Context
		Type    EventType       // All events
		Scope   core.ScopeID    // All events
		Context context.Context // core.FromContext() and scope.Active()
		Metric  core.EventID

		// Arguments (type-specific)
		Attribute  core.KeyValue   // SET_ATTRIBUTE
		Attributes []core.KeyValue // SET_ATTRIBUTES, LOG_EVENT
		Mutator    core.Mutator    // SET_ATTRIBUTE
		Mutators   []core.Mutator  // SET_ATTRIBUTES
		Arguments  []interface{}   // LOGF_EVENT
		Recovered  interface{}     // FINISH_SPAN

		// Values
		String  string // START_SPAN, EVENT, ...
		Float64 float64
		Parent  core.ScopeID // START_SPAN
		Stats   []core.Measurement
	}

	Observer interface {
		Observe(data Event)
	}

	exportersMap map[Observer]struct{}
)

//go:generate stringer -type=EventType
const (
	// TODO: rename these NOUN_VERB
	INVALID EventType = iota
	START_SPAN
	FINISH_SPAN
	LOG_EVENT
	LOGF_EVENT
	SET_GAUGE
	ADD_GAUGE
	NEW_SCOPE
	NEW_MEASURE
	NEW_METRIC
	MODIFY_ATTR
	RECORD_STATS
)

var (
	exporterMu sync.Mutex
	exporters  atomic.Value

	sequenceNum uint64
)

// RegisterObserver adds to the list of Observers that will receive sampled
// trace spans.
//
// Binaries can register exporters, libraries shouldn't register exporters.
func RegisterObserver(e Observer) {
	exporterMu.Lock()
	new := make(exportersMap)
	if old, ok := exporters.Load().(exportersMap); ok {
		for k, v := range old {
			new[k] = v
		}
	}
	new[e] = struct{}{}
	exporters.Store(new)
	exporterMu.Unlock()
}

// UnregisterObserver removes from the list of Observers the Observer that was
// registered with the given name.
func UnregisterObserver(e Observer) {
	exporterMu.Lock()
	new := make(exportersMap)
	if old, ok := exporters.Load().(exportersMap); ok {
		for k, v := range old {
			new[k] = v
		}
	}
	delete(new, e)
	exporters.Store(new)
	exporterMu.Unlock()
}

func Record(event Event) core.EventID {
	if event.Sequence == 0 {
		event.Sequence = NextEventID()
	}
	if event.Time.IsZero() {
		event.Time = time.Now()
	}

	exp, _ := exporters.Load().(exportersMap)
	for e := range exp {
		e.Observe(event)
	}
	return event.Sequence
}

func NextEventID() core.EventID {
	return core.EventID(atomic.AddUint64(&sequenceNum, 1))
}
