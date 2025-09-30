package hello

import (
	"testing"
	"time"
)

type mockService struct {
	calledHello   bool
	calledGoodbye bool
}

func (m *mockService) SayHello(name string) string {
	m.calledHello = true
	return "mock hello"
}

func (m *mockService) SayGoodbye(name string) string {
	m.calledGoodbye = true
	return "mock goodbye"
}

func TestLoggingMiddleware(t *testing.T) {
	mock := &mockService{}
	svc := NewLoggingMiddleware(mock)

	resp := svc.SayHello("Bob")
	if !mock.calledHello {
		t.Error("expected SayHello to be called")
	}
	if resp != "mock hello" {
		t.Errorf("unexpected response: %s", resp)
	}
}

func TestMetricsMiddleware(t *testing.T) {
	mock := &mockService{}
	svc := NewMetricsMiddleware(mock)

	start := time.Now()
	resp := svc.SayGoodbye("Alice")
	if !mock.calledGoodbye {
		t.Error("expected SayGoodbye to be called")
	}
	if resp != "mock goodbye" {
		t.Errorf("unexpected response: %s", resp)
	}
	if time.Since(start) <= 0 {
		t.Error("expected latency measurement")
	}
}
