package hello

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
)

// Decode/Encode for Hello
func DecodeHelloRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req HelloRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// Decode/Encode for Goodbye
func DecodeGoodbyeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req GoodbyeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// NewHTTPHandler with multiple endpoints
func NewHTTPHandler(helloEndpoint, goodbyeEndpoint endpoint.Endpoint) http.Handler {

	mux := http.NewServeMux()

	mux.Handle("/hello", kithttp.NewServer(
		helloEndpoint,
		DecodeHelloRequest,
		EncodeResponse,
	))

	mux.Handle("/goodbye", kithttp.NewServer(
		goodbyeEndpoint,
		DecodeGoodbyeRequest,
		EncodeResponse,
	))

	return mux
}
