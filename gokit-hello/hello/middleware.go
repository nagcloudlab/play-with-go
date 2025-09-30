package hello

import (
	"log"
	"time"
)

type loggingMiddleware struct {
	next Service
}

func (mw loggingMiddleware) SayHello(name string) string {
	defer func(start time.Time) {
		log.Printf("method=SayHello name=%s took=%s\n", name, time.Since(start))
	}(time.Now())

	return mw.next.SayHello(name)
}

func (mw loggingMiddleware) SayGoodbye(name string) string {
	defer func(start time.Time) {
		log.Printf("method=SayGoodbye name=%s took=%s\n", name, time.Since(start))
	}(time.Now())

	return mw.next.SayGoodbye(name)
}

func NewLoggingMiddleware(svc Service) Service {
	return loggingMiddleware{next: svc}
}
