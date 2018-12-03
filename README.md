# gows
CS5600 Final Project, Fast Fault-Tolerant File Server Written in Go

## How to run

## Master Server
The master server's role is to spawn multiple worker servers as child processes, and then monitor the health of the servers and redirect requests to the various servers. The master monitors the health of the servers via the `/healthcheck` endpoint. If a worker server is deemed to be unhealthy, the master kills the server and starts up another (this adds a level of fault-tolerance). The master redirects the requests to worker servers by issuing a HTTP redirect to a worker server, which is chosen via a roud-robin strategy.

## Worker Servers
The worker server implements the basic fileserver functionality. It sets up endpoints to retrieve and return every file in a directory. In addition, the worker servers cache files when they are retrieved, making them able to serve popular files more quickly. 

## Configuration
The master server's parameters can be configured within `conf/config.go`. An example configration is as follows:
```
{
  "workers" :[
    {
      "port" : 9000
    },
    {
      "port" : 9001
    }
  ],
  "baseDirectory": "/Users/robertcarney/tmp/",
  "masterPort": 9090
}
```
Here, `workers` represents configurations for each worker, `baseDirectory` is the directory for the servers to serve files from, and `masterPort` is the port number for the master to listen on.

## Benchmarking

## Potential Improvements










