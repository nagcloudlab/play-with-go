package hello

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type metricsMiddleware struct {
	next           Service
	helloCounter   metric.Int64Counter
	goodbyeCounter metric.Int64Counter
	latencyHist    metric.Float64Histogram
}

func NewMetricsMiddleware(svc Service) Service {
	meter := otel.Meter("hello-service")

	helloCounter, _ := meter.Int64Counter("hello_service_hello_requests_total")
	goodbyeCounter, _ := meter.Int64Counter("hello_service_goodbye_requests_total")
	latencyHist, _ := meter.Float64Histogram("hello_service_request_duration_seconds")

	return &metricsMiddleware{
		next:           svc,
		helloCounter:   helloCounter,
		goodbyeCounter: goodbyeCounter,
		latencyHist:    latencyHist,
	}
}

func (mw *metricsMiddleware) SayHello(name string) string {
	start := time.Now()
	defer mw.latencyHist.Record(context.Background(), time.Since(start).Seconds())

	mw.helloCounter.Add(context.Background(), 1)
	return mw.next.SayHello(name)
}

func (mw *metricsMiddleware) SayGoodbye(name string) string {
	start := time.Now()
	defer mw.latencyHist.Record(context.Background(), time.Since(start).Seconds())

	mw.goodbyeCounter.Add(context.Background(), 1)
	return mw.next.SayGoodbye(name)
}
