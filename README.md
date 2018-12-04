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
The master server's role is to spawn multiple worker servers as child processes, and then monitor the health of the servers and redirect requests to the various servers. The master monitors the health of the servers via the `/healthcheck` endpoint. If a worker server is deemed to be unhealthy, the master kills the server and starts up another (this adds a level of fault-tolerance). The master redirects the requests to worker servers by issuing a HTTP redirect to a worker server, which is chosen via a roud-robin strategy.

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
TODO: Info on benchmarking

## Future Improvements
1. Extend to a distributed protocol: Make it so that the master and client servers can be on different machines.
2. Different port forwarding techinique: Right now the master uses HTTP redirects to send requests. It would be better to use a strategy similar to a reverse proxy, in which the master would act as a middle man between the client and the worker server.
3. Extended Configurations: Add more configurations to allow users more control over the system. 
4. HTTPS Support: For added security, implement an option for communication over HTTPS (could be a config parameter).
