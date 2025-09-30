package hello

import (
	"context"
	"testing"
)

func TestHelloEndpoint(t *testing.T) {
	svc := NewService()
	endpoint := MakeHelloEndpoint(svc)

	req := HelloRequest{Name: "Eve"}
	resp, err := endpoint(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	helloResp := resp.(HelloResponse)
	if helloResp.Message != "Hello, Eve" {
		t.Errorf("unexpected response: %s", helloResp.Message)
	}
}

func TestGoodbyeEndpoint(t *testing.T) {
	svc := NewService()
	endpoint := MakeGoodbyeEndpoint(svc)

	req := GoodbyeRequest{Name: "Eve"}
	resp, err := endpoint(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	goodbyeResp := resp.(GoodbyeResponse)
	if goodbyeResp.Message != "Goodbye, Eve" {
		t.Errorf("unexpected response: %s", goodbyeResp.Message)
	}
}
