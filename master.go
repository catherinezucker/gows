package main

import(
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

func startServer(dir string, port int, serverJobs chan ServerJob)  {
	// Start the server
	cmd := exec.Command("./server", dir, strconv.Itoa(port))
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

func initServers(serverJobs chan ServerJob)  {
	go startServer("/Users/robertcarney", 9090, serverJobs)
	go startServer("/Users/robertcarney", 9091, serverJobs)
}

func monitorServers(serverJobs chan ServerJob, quitChannel chan bool)  {
	for  {
		select {
		case <-quitChannel:
			close(quitChannel)
			return
		default:
			time.Sleep(5 * time.Second)
			currentJob := <- serverJobs
			log.Printf("Looking at server on port: %d\n", currentJob.port)
			if (currentJob.command == nil || currentJob.command.ProcessState != nil) {
				log.Printf("Server process at port: %d with PID: %d exited\n", 
					currentJob.port, currentJob.command.Process.Pid)
			}
			serverJobs <- currentJob
		}
	}
}

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

func main()  {
	log.SetOutput(os.Stderr)
	serverJobs := make(chan ServerJob, 2)
	var quitChannels []chan bool
	go initServers(serverJobs)
	monitorServersQuitChannel := make(chan bool)
	quitChannels = append(quitChannels, monitorServersQuitChannel)
	go monitorServers(serverJobs, monitorServersQuitChannel)
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt, os.Kill)
	go cleanUpOnSignal(signals, serverJobs, quitChannels)
	for {  }
}