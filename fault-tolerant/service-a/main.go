package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

// -------- Service A Business Logic --------
type Service interface {
	CallB(name string) (string, error)
}

type callBService struct {
	callBEndpoint endpoint.Endpoint
	cb            *gobreaker.CircuitBreaker
	semaphore     chan struct{} // Bulkhead limiter
}

// CallB calls Service B with bulkhead + breaker + retry + timeout
func (s callBService) CallB(name string) (string, error) {
	var (
		resp interface{}
		err  error
	)

	// 🧱 Bulkhead: acquire a slot
	select {
	case s.semaphore <- struct{}{}:
		defer func() { <-s.semaphore }() // release slot
	default:
		return "", errors.New("bulkhead: too many concurrent requests")
	}

	// 🔁 Retry loop
	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		result, err := s.cb.Execute(func() (interface{}, error) {
			return s.callBEndpoint(ctx, HelloRequest{Name: name})
		})

		if err == nil {
			resp = result
			r := resp.(HelloResponse)
			return r.Message, nil
		}

		log.Printf("⚠️ attempt %d failed: %v", i+1, err)
		time.Sleep(200 * time.Millisecond)
	}

	return "", err
}

// -------- Shared Types with Service B --------
type HelloRequest struct {
	Name string `json:"name"`
}
type HelloResponse struct {
	Message string `json:"message"`
}

// -------- Endpoint for Service A --------
type CallBRequest struct {
	Name string `json:"name"`
}
type CallBResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func makeCallBEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CallBRequest)
		msg, err := svc.CallB(req.Name)
		if err != nil {
			return CallBResponse{Message: "", Error: err.Error()}, nil
		}
		return CallBResponse{Message: msg}, nil
	}
}

// -------- Transports --------
func decodeCallBRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req CallBRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// Client encoders/decoders for Service B
func encodeHelloRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = io.NopCloser(&buf)
	r.Header.Set("Content-Type", "application/json")
	return nil
}

func decodeHelloResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response HelloResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func main() {
	// Build client to Service B
	u, _ := url.Parse("http://localhost:8082/hello")
	callBEndpoint := httptransport.NewClient(
		"POST",
		u,
		encodeHelloRequest,
		decodeHelloResponse,
	).Endpoint()

	// Configure Circuit Breaker
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "ServiceB",
		MaxRequests: 1,
		Interval:    60 * time.Second,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
	})

	// Configure Bulkhead (allow max 5 concurrent calls)
	semaphore := make(chan struct{}, 5)

	// Build Service A
	svc := callBService{
		callBEndpoint: callBEndpoint,
		cb:            cb,
		semaphore:     semaphore,
	}

	// Build endpoint
	endpoint := makeCallBEndpoint(svc)

	// 🚦 Add Rate Limiter (2 requests per second max)
	rl := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Limit(2), 1))
	endpoint = rl(endpoint)

	// Expose HTTP handler
	handler := httptransport.NewServer(
		endpoint,
		decodeCallBRequest,
		encodeResponse,
	)

	http.Handle("/call-b", handler)

	log.Println("🟢 Service A (with rate limit + bulkhead + breaker + retries + timeout) running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
