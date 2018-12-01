# GOWS Brainstorm, Implementation Blueprint

## Master Pseudocode

### Init
* Command line arg is path to config file
* Read config file into conf struct
  * Abort on failures
    * File doesn't exist
    * JSON given doesn't match config spec.
* Parse each worker config into worker struct
* Store workers info in some data structure
  * channel, stack, queue?

### Deploy Workers
* Set up master/worker communication pipeline
  * Options:
    * SSH tunnels (look into API or just exec shell commands)
    * Maybe something else exists... TODO
* Execute workers
    * On failure ideas:
      * Abort entire thing
      * Retry 'x' times (can be a config param)
      * Mark as 'bad' and continue with all other workers

### Post-Setup / Body
* Start listening on IP:Port we will be serving
  * Again, handle failures with abort/retries
* Do forever:
  * Listen on port, delegate to workers
  * Could be multithreaded, look into atomicity if multiple listeners
    (maybe master puts requests into shared queue?)
* Health check
  * Thread for each worker
    * Do forever:
      * Issue GET request to worker's health check endpoint
      * Wait <timeout> number of seconds for response
      * If response is "Status OK":
        * Sleep <HeathCheckInterval>
      * Else:
        * Kill worker
          * If process is still running, kill:
            * Maybe workers have signal handler for graceful shutdown to finish any incomplete requests?
            * Could also log incomplete requests in tmp file which replacement could read on startup
            * Otherwise Hard Kill 
        * Replace worker
        * Sleep <HeathCheckInterval>

## Ideas
* Look into logging in Go for master so we can get live status updates, should be simple
* Scale workers up and down based on load
* Implement caching

## TO-DO
* Donald:
  * Figure out Go constructs for:
    * packaging
