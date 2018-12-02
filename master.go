package main

import(
	"fmt"
	"net/http"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"time"
)

// ServerJob represents a server worker process
type ServerJob struct {
	port int
	dir string
	command *exec.Cmd
}

// Starts a single server and adds it on to the channel
func startServer(dir string, port int, serverJobs chan ServerJob)  {
	// Start the server
	cmd := exec.Command("./server/server", dir, strconv.Itoa(port))
	err := cmd.Start()
	if err != nil {
		log.Printf("Server at port: %d failed to start\n", port)
		return
	}
	// Add server to the channel of server jobs
	serverJobs <- ServerJob{port, dir, cmd}
	log.Printf("Started server for directory %s at port :%d with PID: %d\n", dir, port, cmd.Process.Pid)
	err = cmd.Wait()
	log.Printf("Server at port :%d exited with error code %v\n", port, err)
}

// Starts servers and adds them on to the channel
func initServers(serverJobs chan ServerJob)  {
	go startServer("/Users/robertcarney", 9000, serverJobs)
	go startServer("/Users/robertcarney", 9001, serverJobs)
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
			log.Printf("Looking at server on port: %d\n", currentJob.port)
			// Monitor the job...
			if (!serverIsHealthy(currentJob)) {
				log.Printf("Server process at port: %d with PID: %d exited\n", 
					currentJob.port, currentJob.command.Process.Pid)
				continue
			}
			// And put it back on the channel
			serverJobs <- currentJob
		}
	}
}

// Checks on the health of a server job by sending a test request to the /healthcheck endpoint
func serverIsHealthy(serverJob ServerJob) bool {
	// TODO
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


func redirectOnChannel(ports chan int) func(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In redirectOnChannel")
	return func(w http.ResponseWriter, r *http.Request)  {
		fmt.Println("In return func")
		currentPort := <-ports
		fmt.Printf("Port recieved from channel was %d\n", currentPort)
		ports <- currentPort
		http.Redirect(w, r, fmt.Sprintf("http://localhost:%d", currentPort), 301)
	}
}

func main()  {
	log.SetOutput(os.Stderr)
	serverJobs := make(chan ServerJob, 2)
	var quitChannels []chan bool
	currentChannel := make(chan int, 2)
	currentChannel <- 9000
	currentChannel <- 9001
	go initServers(serverJobs)
	monitorServersQuitChannel := make(chan bool)
	quitChannels = append(quitChannels, monitorServersQuitChannel)
	go monitorServers(serverJobs, monitorServersQuitChannel)
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt, os.Kill)
	go cleanUpOnSignal(signals, serverJobs, quitChannels)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request)  {
		currentPort := <-currentChannel
		fmt.Printf("Port recieved from channel was %d\n", currentPort)
		currentChannel <- currentPort
		http.Redirect(w, r, fmt.Sprintf("http://localhost:%d", currentPort), 301)
	})
	http.ListenAndServe(":9090", nil)
}