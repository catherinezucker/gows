# gows
CS5600 Final Project, Fast Fault-Tolerant File Server Written in Go

## Master Server
The master server's role is to spawn multiple worker servers as child processes, and then monitor the health of the servers and redirect requests to the various servers. The master monitors the health of the servers via the `/healthcheck` endpoint. The master forwards the requests to worker servers by issuing a HTTP redirect to a worker server, which is chosen via a roud-robin strategy.

## Worker Servers













