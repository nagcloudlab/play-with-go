package main

import (
	"context"
	lg "log"
	"net/http"
	"time"

	"gokit-hello/hello"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func initResources(ctx context.Context) *resource.Resource {
	res, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("hello-service"),
			semconv.ServiceVersion("1.0.0"),
			attribute.String("environment", "dev"),
		),
	)
	return res
}

// ---- Metrics ----
func initMeterProvider(ctx context.Context, res *resource.Resource) *metric.MeterProvider {
	exporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint("localhost:4318"),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		lg.Fatalf("failed to create OTLP metrics exporter: %v", err)
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter, metric.WithInterval(2*time.Second))),
		metric.WithResource(res),
	)

	otel.SetMeterProvider(provider)
	return provider
}

// ---- Traces ----
func initTracerProvider(ctx context.Context, res *resource.Resource) *trace.TracerProvider {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		lg.Fatalf("failed to create OTLP trace exporter: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	return tp
}

// ---- Logs ----
// func initLoggerProvider(ctx context.Context, res *resource.Resource) *log.LoggerProvider {
// 	exporter, err := otlploghttp.New(ctx,
// 		otlploghttp.WithEndpoint("localhost:4318"),
// 		otlploghttp.WithInsecure(),
// 	)
// 	if err != nil {
// 		lg.Fatalf("failed to create OTLP log exporter: %v", err)
// 	}
// 	lp := log.NewLoggerProvider(
// 		log.WithProcessor(log.NewBatchProcessor(exporter)),
// 		log.WithResource(res),
// 	)
// 	otel.SetLoggerProvider(lp)
// 	return lp
// }

func main() {
	ctx := context.Background()
	res := initResources(ctx)

	// Init providers
	meterProvider := initMeterProvider(ctx, res)
	tracerProvider := initTracerProvider(ctx, res)
	//loggerProvider := initLoggerProvider(ctx, res)

	defer func() {
		_ = meterProvider.Shutdown(ctx)
		_ = tracerProvider.Shutdown(ctx)
		// _ = loggerProvider.Shutdown(ctx)
	}()

	// ---- Go Kit Service ----
	var svc hello.Service
	svc = hello.NewService()
	svc = hello.NewLoggingMiddleware(svc)
	svc = hello.NewMetricsMiddleware(svc)
	svc = hello.NewTracingMiddleware(svc) // 👈 NEW: traces
	//svc = hello.NewLoggingOTelMiddleware(svc) // 👈 NEW: OTel logs

	helloEndpoint := hello.MakeHelloEndpoint(svc)
	goodbyeEndpoint := hello.MakeGoodbyeEndpoint(svc)
	handler := hello.NewHTTPHandler(helloEndpoint, goodbyeEndpoint)

	http.Handle("/", handler)

	lg.Println("🚀 Server running on :8080")
	lg.Fatal(http.ListenAndServe(":8080", nil))
}
