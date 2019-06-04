This is not a complete implementation of the OpenTracing and OpenCensus API surface areas. I'm posting this here, now, to have as a point of reference for several of the issues in the specification repo. Some of the high-lights of the approach taken here:

* Always associate current span context with stats and metrics events
* Introduce a low-level "observer" exporter
* Avoid excessive memory allocations
* Avoid buffering state with API objects
* Use `context.Context` to propagate context tags and active scope
* Introduce a "reader" implementation to interpret "observer"-exported events and build state
* Use a common `KeyValue` construct for span attributes, context tags, resource definitions, log fields, and metric fields
* Support logging API w/o a current span context
* Support for golang `plugin` package to load implementations
* Example use of golang's `net/http/httptrace` w/ @iredelmeier's [tracecontext.go](https://github.com/lightstep/tracecontext.go) package

The first bullet about associating current span context and stats/metrics events bridges the tracing data model with the metrics data model. The APIs here would make this association not an option, as the `stats.Record` API takes a context, which passes through to the observer, which could choose to use the span-context association. The prototype includes a stderr exporter that writes a debugging log of events to the console. One of the critical features here, enabled by the low-level observer exporter, is that Span start events can be logged in chronological order, not as the span finishes.

To run the examples, first build the stderr tracer plugin (requires Linux or OS X):

```
(cd ./exporter/stderr/plugin && make)
```

then set the `OPENTELEMETRY_LIB` environment variable to the .so file in that directory:

```
OPENTELEMETRY_LIB=./exporter/stderr/plugin/stderr.so go run ./example/client/client.go
```
