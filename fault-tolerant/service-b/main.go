package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// -------- Business Logic --------
type Service interface {
	Hello(name string) string
}

type helloService struct{}

func (s helloService) Hello(name string) string {
	// Simulate random latency
	n := rand.Intn(100)
	if n < 30 {
		time.Sleep(3 * time.Second) // 30% chance of 3s delay
	}
	return "Hello " + name + " from Service B!"
}

// -------- Endpoint --------
type HelloRequest struct {
	Name string `json:"name"`
}
type HelloResponse struct {
	Message string `json:"message"`
}

func makeHelloEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(HelloRequest)
		msg := svc.Hello(req.Name)
		return HelloResponse{Message: msg}, nil
	}
}

// -------- Transport --------
func decodeHelloRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req HelloRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	svc := helloService{}
	helloEndpoint := makeHelloEndpoint(svc)

	handler := httptransport.NewServer(
		helloEndpoint,
		decodeHelloRequest,
		encodeResponse,
	)

	http.Handle("/hello", handler)

	log.Println("🚀 Service B running on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
