package hello

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Request format
type HelloRequest struct {
	Name string `json:"name"`
}

// Response format
type HelloResponse struct {
	Message string `json:"message"`
}

// MakeHelloEndpoint converts service method into a Go Kit endpoint
func MakeHelloEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(HelloRequest) // cast request
		msg := s.SayHello(req.Name)
		return HelloResponse{Message: msg}, nil
	}
}

// Goodbye
type GoodbyeRequest struct {
	Name string `json:"name"`
}
type GoodbyeResponse struct {
	Message string `json:"message"`
}

func MakeGoodbyeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GoodbyeRequest)
		return GoodbyeResponse{Message: s.SayGoodbye(req.Name)}, nil
	}
}
