package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/go-kit/kit/sd/consul"
	gokitlog "github.com/go-kit/log"
	"github.com/hashicorp/consul/api"
)

// -------- Business Logic --------
type Service interface {
	Hello(name string) string
}

type helloService struct{}

func (s helloService) Hello(name string) string {
	// n := rand.Intn(100)
	// if n < 30 {
	// 	time.Sleep(2 * time.Second) // simulate slowness
	// }
	return "Hello " + name + " from Service B!" + os.Getenv("INSTANCE_ID")
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

// Health check
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// -------- Main --------
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
	http.HandleFunc("/health", healthHandler)

	instanceId := os.Getenv("INSTANCE_ID")
	port_str := os.Getenv("PORT")

	//convert port string to int
	port, err := strconv.Atoi(port_str)
	if err != nil {
		log.Fatal("Invalid PORT:", err)
	}

	addr := ":" + port_str
	log.Println("🚀 Service B running on", addr)

	// ---- Register with Consul ----
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal("consul client error:", err)
	}
	consulClient := consul.NewClient(client)
	logger := gokitlog.NewLogfmtLogger(log.Writer())

	registrar := consul.NewRegistrar(consulClient, &api.AgentServiceRegistration{
		ID:      instanceId,
		Name:    "service-b",
		Address: "host.docker.internal", // for Docker Consul to reach host
		Port:    port,
		Check: &api.AgentServiceCheck{
			HTTP:     "http://host.docker.internal:8082/health",
			Interval: "10s",
			Timeout:  "1s",
		},
	}, logger)

	registrar.Register()
	defer registrar.Deregister()

	log.Fatal(http.ListenAndServe(addr, nil))
}
