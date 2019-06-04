package scope

import (
	"github.com/lightstep/opentelemetry-golang-prototype/api/core"
	"github.com/lightstep/opentelemetry-golang-prototype/impl/observer"
)

// TODO Rename scope.New
func Start(parent core.ScopeID, attributes ...core.KeyValue) core.Scope {
	eventID := observer.Record(observer.Event{
		Type:       observer.NEW_SCOPE,
		Scope:      parent,
		Attributes: attributes,
	})
	return &scopeIdent{
		id: core.ScopeID{
			EventID:     eventID,
			SpanContext: parent.SpanContext,
		},
	}
}
