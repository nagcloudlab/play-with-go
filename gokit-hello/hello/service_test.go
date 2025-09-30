package hello

import "testing"

func TestSayHello(t *testing.T) {
	svc := NewService()
	got := svc.SayHello("Bob")
	want := "Hello, Bob"

	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestSayGoodbye(t *testing.T) {
	svc := NewService()
	got := svc.SayGoodbye("Bob")
	want := "Goodbye, Bob"

	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
