package hello

import (
	"context"

	"go.opentelemetry.io/otel"
)

type tracingMiddleware struct {
	next Service
}

func NewTracingMiddleware(svc Service) Service {
	return &tracingMiddleware{next: svc}
}

func (mw *tracingMiddleware) SayHello(name string) string {
	_, span := otel.Tracer("hello-service").Start(context.Background(), "SayHello")
	defer span.End()
	return mw.next.SayHello(name)
}

func (mw *tracingMiddleware) SayGoodbye(name string) string {
	_, span := otel.Tracer("hello-service").Start(context.Background(), "SayGoodbye")
	defer span.End()
	return mw.next.SayGoodbye(name)
}
