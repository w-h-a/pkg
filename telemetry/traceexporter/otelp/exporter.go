package otelp

import (
	"context"

	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/pkg/telemetry/traceexporter"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

type otelpExporter struct {
	options traceexporter.ExporterOptions
	*otlptrace.Exporter
}

func (e *otelpExporter) Options() traceexporter.ExporterOptions {
	return e.options
}

func (e *otelpExporter) String() string {
	return "otelp"
}

func (e *otelpExporter) configure() error {
	var client otlptrace.Client

	var err error

	if e.options.Protocol == "grpc" {
		clientOpts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(e.options.Nodes[0]),
		}
		if !e.options.Secure {
			clientOpts = append(clientOpts, otlptracegrpc.WithInsecure())
		}
		client = otlptracegrpc.NewClient(clientOpts...)
	} else {
		log.Infof("MY HEADERS %+v", e.options.Headers)
		log.Infof("MY SECURE %v", e.options.Secure)
		clientOpts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(e.options.Nodes[0]),
			otlptracehttp.WithHeaders(e.options.Headers),
		}
		if !e.options.Secure {
			clientOpts = append(clientOpts, otlptracehttp.WithInsecure())
		}
		client = otlptracehttp.NewClient(clientOpts...)
	}

	e.Exporter, err = otlptrace.New(context.Background(), client)
	if err != nil {
		return err
	}

	return nil
}

func NewExporter(opts ...traceexporter.ExporterOption) traceexporter.TraceExporter {
	options := traceexporter.NewExporterOptions(opts...)

	e := &otelpExporter{
		options: options,
	}

	if err := e.configure(); err != nil {
		log.Fatal(err)
	}

	return e
}
