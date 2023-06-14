package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}

type simpleServer struct {
	addr	string
	proxy 	*httputil.ReverseProxy
}

type LoadBalancer struct {
	port	string
	roundRobinCount int
	servers []Server
}

func newSimpleServer(addr string) *simpleServer {
	serverURL, err := url.Parse(addr)
	HandleError(err)

	return &simpleServer{
		addr: addr,
		proxy: httputil.NewSingleHostReverseProxy(serverURL),
	}
}

func (s *simpleServer) Address() string { return s.addr }

func (s *simpleServer) IsAlive() bool {
	_, err := http.Get(s.addr)
	return err == nil
}

func (s *simpleServer) Serve(rw http.ResponseWriter, r *http.Request) { s.proxy.ServeHTTP(rw, r) }

func HandleError(err error) {
	if err != nil {
		fmt.Printf("error %v\n", err)
		os.Exit(1)
	}
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		roundRobinCount : 0,
		port: port,
		servers: servers,
	}
}

func (lb *LoadBalancer) getNextAvailableServer() Server{
	nextServer := lb.servers[lb.roundRobinCount]
	for !nextServer.IsAlive() {
		lb.roundRobinCount = (lb.roundRobinCount + 1) % len(lb.servers)
		nextServer = lb.servers[lb.roundRobinCount]
	}
		lb.roundRobinCount++
	return nextServer
}

func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, r *http.Request ) {
	targetServer := lb.getNextAvailableServer()
	targetServer.Serve(rw, r)
	fmt.Printf("Proxying request to server %s\n", targetServer.Address())
}

func main() {
	servers := []Server{
		newSimpleServer("http://localhost:3000"),
		newSimpleServer("http://localhost:3001"),
		newSimpleServer("http://localhost:3002"),
	}

	lb := NewLoadBalancer(":8080", servers)

	handleRedirect := func(rw http.ResponseWriter, r *http.Request) {
		lb.serveProxy(rw, r)
	}
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("Load Balancer started at %s\n", lb.port)

	if err := http.ListenAndServe(lb.port, nil); err != nil {
		panic(err)
	}
}