package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/lightstep/sandbox/jmacd/otel/plugin/httptrace"
	"github.com/lightstep/sandbox/jmacd/otel/tag"
	"github.com/lightstep/sandbox/jmacd/otel/trace"

	// This creates a debug log on the console.
	_ "github.com/lightstep/sandbox/jmacd/otel/observer/debug"
)

var (
	tracer = trace.GlobalTracer().
		WithService("client").
		WithComponent("main").
		WithResources(
			tag.New("whatevs").String("yesss"),
		)
)

func main() {
	client := http.DefaultClient

	ctx := tag.NewContext(context.Background(),
		tag.Insert(tag.New("username").String("donuts")),
	)

	var body []byte

	err := tracer.WithSpan(ctx, "say hello",
		func(ctx context.Context) error {
			req, _ := http.NewRequest("GET", "http://localhost:7777/hello", nil)

			ctx, req, inj := httptrace.W3C(ctx, req)

			trace.Inject(ctx, inj)

			res, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			body, err = ioutil.ReadAll(res.Body)
			res.Body.Close()

			return err
		})

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", body)
}
