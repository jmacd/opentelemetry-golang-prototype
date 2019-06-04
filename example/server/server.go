package main

import (
	"io"
	"net/http"

	"github.com/lightstep/sandbox/jmacd/otel/core"
	"github.com/lightstep/sandbox/jmacd/otel/log"
	"github.com/lightstep/sandbox/jmacd/otel/plugin/httptrace"
	"github.com/lightstep/sandbox/jmacd/otel/tag"
	"github.com/lightstep/sandbox/jmacd/otel/trace"

	// This creates a debug log on the console.
	_ "github.com/lightstep/sandbox/jmacd/otel/observer/debug"
)

var (
	tracer = trace.GlobalTracer().
		WithService("server").
		WithComponent("main").
		WithResources(
			tag.New("whatevs").String("nooooo"),
		)
)

func main() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		attrs, tags, spanCtx := httptrace.ServerHeaders(req)

		req = req.WithContext(tag.WithMap(req.Context(), tag.NewMap(core.KeyValue{}, tags, core.Mutator{}, nil)))

		// @@@ TODO spanCtx parent issues
		_ = spanCtx

		ctx, span := tracer.Start(
			req.Context(),
			"hello",
			attrs...,
		)
		defer span.Finish()

		log.Log(ctx, "handling this...")

		io.WriteString(w, "Hello, world!\n")
	}

	http.HandleFunc("/hello", helloHandler)
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		panic(err)
	}
}
