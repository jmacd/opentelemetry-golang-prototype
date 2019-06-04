package main

import (
	"io"
	"net/http"

	"github.com/lightstep/opentelemetry-golang-prototype/api/core"
	"github.com/lightstep/opentelemetry-golang-prototype/api/log"
	"github.com/lightstep/opentelemetry-golang-prototype/api/tag"
	"github.com/lightstep/opentelemetry-golang-prototype/api/trace"
	"github.com/lightstep/opentelemetry-golang-prototype/plugin/httptrace"

	// This creates a debug log on the console.
	_ "github.com/lightstep/opentelemetry-golang-prototype/exporter/stderr"
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
