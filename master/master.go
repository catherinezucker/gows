package master

import(
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/robcarney/gows/config"
)

var baseDirectory string

// ServerJob represents a server worker process
type ServerJob struct {
	host string
	port int
	dir string
	command *exec.Cmd
}

// Starts a single server and adds it on to the channel
func startServer(port int, host string) (ServerJob, error) {
	// Start the server
	cmd := exec.Command("./server/server", baseDirectory, strconv.Itoa(port))
	err := cmd.Start()
	if err != nil {
		log.Printf("Server at port: %d failed to start\n", port)
	}
	log.Printf("Started server for directory %s at port :%d with PID: %d\n", 
		baseDirectory, port, cmd.Process.Pid)
	return ServerJob{
		host: host,
		port: port,
		dir: baseDirectory,
		command: cmd,
	}, err
}

// Starts servers and adds them on to the channel
func initServers(serverJobs chan ServerJob, workers []config.ServerWorker)  {
	for _, worker := range workers  {
		job, err := startServer(worker.Port, worker.Host)
		if (err == nil)  {
			serverJobs <- job
		}
	}
}

// Monitors the servers in serverJobs, and returns when it recieves on the quitChannel
func monitorServers(serverJobs chan ServerJob, quitChannel chan bool)  {
	defer close(quitChannel)
	for  {
		select {
		case <-quitChannel:
			// We recieved a signal to stop monitoring
			return
		default:
			// We only want to monitor on some time interval, so sleep first
			//   NOTE: This will make the master process take roughly this amount
			//   of time to tear down when it recieves an exit signal
			time.Sleep(5 * time.Second)
			// Get a job from the channel...
			currentJob := <- serverJobs
			log.Printf("Looking at server at %s:%d\n", currentJob.host, currentJob.port)
			// Monitor the job...
			if (!serverIsHealthy(currentJob)) {
				// ensure dead
				currentJob.command.Process.Kill()
				// removes zombie process
				currentJob.command.Process.Wait()
				log.Printf("Server process at port: %d with PID: %d exited\n", 
					currentJob.port, currentJob.command.Process.Pid)
				newJob, err := startServer(currentJob.port, currentJob.host)
				if err != nil {
					log.Printf("Unable to restart server process at port: %d\n%s\n", currentJob.port, err.Error())
				} else {
					log.Printf("Restarted Server process at port: %d with PID: %d\n",
						currentJob.port, currentJob.command.Process.Pid)
					serverJobs <- newJob
				}
			} else {
				// And put it back on the channel
				serverJobs <- currentJob
			}
		}
	}
}

// Checks on the health of a server job by sending a test request to the /healthcheck endpoint
func serverIsHealthy(serverJob ServerJob) bool {
	healthCheckEndpoint := "http://" + serverJob.host + ":" + strconv.Itoa(serverJob.port) + "/healthcheck"
	res, err := http.Head(healthCheckEndpoint)
	if err != nil {
		log.Printf("Error in healthcheck at %s:\n%s\n", healthCheckEndpoint, err.Error())
		return false
	} else if res.StatusCode != http.StatusOK {
		log.Printf("Healthcheck fail, status code %d\n", res.StatusCode)
		return false
	}
	log.Printf("Healthcheck to %s succeeded\n", healthCheckEndpoint)
    return true
}

// Cleans up the running child processes (servers) when the program recieves an error signal
func cleanUpOnSignal(signals chan os.Signal, serverJobs chan ServerJob, quitChannels []chan bool)  {
	defer close(signals)
	<-signals
	for _, quitChannel := range quitChannels {
		quitChannel <- true
	}
	close(serverJobs)
	for currentJob := range serverJobs {
		if (currentJob.command != nil) {
			currentJob.command.Process.Kill()
		}
	}
	os.Exit(1)
}

// Returns a http handler function requests to a port popped off the given channel for the given endpoint
func redirectOnChannel(serverJobs chan ServerJob, endpoint string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request)  {
		// Get a job from the channel, and immidiately add it back in
		currentJob := <-serverJobs
		serverJobs <- currentJob
		url, _ := url.Parse(fmt.Sprintf("http://%s:%d", currentJob.host, currentJob.port))
		proxy := httputil.NewSingleHostReverseProxy(url)

		// Set the new URL for the request
		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = url.Host

		// Serve this request via our reverse proxy
		proxy.ServeHTTP(w, r)
	}
}

// Returns a function to be called for each file and directory in the base directory with filepath.Walk
//   The function returned sets up an endpoint for the given file
func getFileVisitor(serverJobs chan ServerJob) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// No endpoints for directories
		if info.IsDir()  {
			return nil
		}
		// Get the relative path to be used in the endpoint
		relativePath := strings.Replace(path, baseDirectory, "", 1)
		if relativePath[:0] != "/"  {
			relativePath = "/" + relativePath
		}
		log.Printf("Setting up endpoint for %s\n", relativePath)
		// Bind the http handler
		http.HandleFunc(relativePath, redirectOnChannel(serverJobs, relativePath))
		return nil
	}
}

// Sets up endpoints for all files in the base directory
func setUpFileEndpoints(serverJobs chan ServerJob)  {
	err := filepath.Walk(baseDirectory, getFileVisitor(serverJobs))
	if err != nil  {
		log.Fatal("setUpFileEndpoints: ", err)
	}
}

// Run starts the master server based on the given configuration parameters
func Run(serverConfig config.Config)  {
	baseDirectory = serverConfig.BaseDirectory
	serverJobs := make(chan ServerJob, len(serverConfig.Workers))
	var quitChannels []chan bool
	ports := make(chan int, len(serverConfig.Workers))
	for _, worker := range serverConfig.Workers  {
		ports <- worker.Port
	}
	initServers(serverJobs, serverConfig.Workers)
	monitorServersQuitChannel := make(chan bool)
	quitChannels = append(quitChannels, monitorServersQuitChannel)
	go monitorServers(serverJobs, monitorServersQuitChannel)
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt, os.Kill)
	go cleanUpOnSignal(signals, serverJobs, quitChannels)
	setUpFileEndpoints(serverJobs)
	http.ListenAndServe(fmt.Sprintf(":%d", serverConfig.MasterPort), nil)
}






