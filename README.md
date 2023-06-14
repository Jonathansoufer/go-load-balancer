# Load Balancer

This is a simple load balancer written in Go. It allows you to distribute incoming requests across multiple servers.

## Usage

To use the load balancer, you need to create a `LoadBalancer` instance and add one or more `Server` instances to it. Here's an example:

```go
package main

import (
    "net/http"
)

func main() {
    // Create a load balancer
    lb := &LoadBalancer{
        port: "8080",
    }

    // Add some servers
    lb.AddServer(newSimpleServer("http://localhost:8000"))
    lb.AddServer(newSimpleServer("http://localhost:8001"))
    lb.AddServer(newSimpleServer("http://localhost:8002"))

    // Start the load balancer
    lb.Start()
}
```

This will start a load balancer on port 8080 that distributes incoming requests across three servers running on ports 8000, 8001, and 8002.

Contributing
Contributions are welcome! If you find a bug or have a feature request, please open an issue on GitHub.
