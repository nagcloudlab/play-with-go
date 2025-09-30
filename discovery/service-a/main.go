package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	gokitlog "github.com/go-kit/log"
	"github.com/hashicorp/consul/api"
)

// -------- Request/Response --------
type HelloRequest struct {
	Name string `json:"name"`
}
type HelloResponse struct {
	Message string `json:"message"`
}

// -------- Client Endpoint Factory --------
func makeHelloProxy(instance string) (endpoint.Endpoint, io.Closer, error) {
	// Replace host.docker.internal with localhost for Mac
	if strings.Contains(instance, "host.docker.internal") {
		instance = strings.Replace(instance, "host.docker.internal", "localhost", 1)
	}

	tgt := "http://" + instance + "/hello"
	u, err := url.Parse(tgt)
	if err != nil {
		return nil, nil, err
	}

	return httptransport.NewClient(
		"POST",
		u,
		encodeRequest,
		decodeResponse,
	).Endpoint(), nil, nil
}

func encodeRequest(_ context.Context, r *http.Request, request interface{}) error {
	r.Header.Set("Content-Type", "application/json")
	if request == nil {
		return nil
	}
	buf, err := json.Marshal(request)
	if err != nil {
		return err
	}
	r.Body = io.NopCloser(bytes.NewReader(buf))
	r.ContentLength = int64(len(buf))
	return nil
}

func decodeResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp HelloResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// -------- Transport for Service A --------
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

// -------- Main --------
func main() {
	cfg := api.DefaultConfig()
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	consulClient := consul.NewClient(client)
	logger := gokitlog.NewLogfmtLogger(log.Writer())

	// Discover service-b
	instancer := consul.NewInstancer(consulClient, logger, "service-b", []string{}, true)

	endpointer := sd.NewEndpointer(instancer, makeHelloProxy, logger)
	balancer := lb.NewRoundRobin(endpointer)
	helloEndpoint := lb.Retry(3, 2*time.Second, balancer)

	handler := httptransport.NewServer(
		helloEndpoint,
		decodeHelloRequest,
		encodeResponse,
	)

	http.Handle("/call-b", handler)

	addr := ":8081"
	log.Println("🚀 Service A running on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
