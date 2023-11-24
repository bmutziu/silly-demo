package main

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"os"

	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var tp *sdktrace.TracerProvider

func initTracer() (*sdktrace.TracerProvider, error) {
	url := os.Getenv("JAEGER_ENDPOINT")
	if len(url) > 0 {
		return initFileTracer()
	} else {
		return initFileTracer()
	}
}

func initFileTracer() (*sdktrace.TracerProvider, error) {
	// exporter, err := stdout.New(stdout.WithPrettyPrint())
	f, err := os.Create("traces.json")
	if err != nil {
		return nil, err
	}
	exporter, err := stdout.New(
		stdout.WithWriter(f),
		stdout.WithPrettyPrint(),
	)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}
