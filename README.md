# gows
CS5600 Final Project, Fast Fault-Tolerant File Server Written in Go

## How to run
There is a Makefile included for building and running a simple server. 
To run a server:
1. Clone the repository to `{GOPATH}/github.com/robcarney/`
2. Edit the `conf/config.json` file. Choose a directory to serve and approriate ports to use. You can also add or remove workers.
3. Run `make run` on the command line to run the default configurations.
4. Navigate to `http://localhost:{PORT}/{RELATIVE FILE PATH}` to view a file in the browser.

## Master Server
The master server's role is to spawn multiple worker servers as child processes, and then monitor the health of the servers and redirect requests to the various servers. The master monitors the health of the servers via the `/healthcheck` endpoint. If a worker server is deemed to be unhealthy, the master kills the server and starts up another (this adds a level of fault-tolerance). The master redirects the requests to worker servers by setting up a remote proxy and forwarding the request to a worker server, which is chosen via a roud-robin strategy.

## Worker Servers
The worker server implements the basic fileserver functionality. It sets up endpoints to retrieve and return every file in a directory. In addition, the worker servers cache files when they are retrieved, making them able to serve popular files more quickly. 

## Configuration
The master server's parameters can be configured within `conf/config.go`. An example configration is as follows:
```
{
  "workers" :[
    {
      "host": "127.0.0.1",
      "port" : 9000
    },
    {
      "host": "127.0.0.1",
      "port" : 9001
    }
  ],
  "baseDirectory": "/Users/robertcarney/tmp/",
  "masterPort": 9090
}
```
* `workers` represents configurations for each worker 
  * `host` represents the host address of the worker, currently only localhost (127.0.0.1) is supported as we do not yet support a distrubuted architecture
  * `port` represents the port for the worker to listen on
* `baseDirectory` is the directory for the servers to serve files from
* `masterPort` is the port number for the master to listen on.

## Benchmarking
ApacheBench was used to benchmark the server. Benchmarks were taken with 2, 10, and 100 workers using 5000 requests with 500 concurrent requests. The results were then compared against the standard Apache2 server. Both 2 and 10 workers performed noticeably worse than the Apache2 server for the majority of the test until the requests in the upper 4000s, at which point the Apache server's time increased exponentially bringing the overall performance of it significantly down. The Go web server on the other hand increased at a much lower rate over time and maintained good performance throughout. The Go server with 100 workers performed similarly to the Apache2 server through the early requests and only increased in response time slightly towards the end, making it overall the best performer. Of all four configurations tested only the Apache2 server had any failed requests while benchmarking.

Response time for individual requests was also measured with the Go server responding about 11% slower. This suggests that the Apache2 server is better built for individual requests while the Go web server can better handle large numbers of concurrent requests.

## Future Improvements
1. Extend to a distributed protocol: Make it so that the master and client servers can be on different machines.
2. Extended Configurations: Add more configurations to allow users more control over the system. 
