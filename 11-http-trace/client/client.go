package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"

	/* stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace" */
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

func initTracer() *sdktrace.TracerProvider {
	// Create stdout exporter to be able to retrieve
	// the collected spans.
	//exporter, err := stdout.New(stdout.WithPrettyPrint())
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))

	if err != nil {
		log.Fatal(err)
	}

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func main() {
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	url := flag.String("server", "http://localhost:7777/hello", "server url")
	flag.Parse()

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	bag, _ := baggage.Parse("username=donuts")
	ctx := baggage.ContextWithBaggage(context.Background(), bag)

	var body []byte

	tr := otel.Tracer("example/client")
	err := func(ctx context.Context) error {
		ctx, span := tr.Start(ctx, "say hello", trace.WithAttributes(semconv.PeerServiceKey.String("ExampleService")))
		defer span.End()
		req, _ := http.NewRequestWithContext(ctx, "GET", *url, nil)

		fmt.Printf("Sending request...\n")
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		body, err = ioutil.ReadAll(res.Body)
		_ = res.Body.Close()

		return err
	}(ctx)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Response Received: %s\n\n\n", body)
	fmt.Printf("Waiting for few seconds to export spans ...\n\n")
	time.Sleep(10 * time.Second)
	fmt.Printf("Inspect traces on stdout\n")
}
