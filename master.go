package main

import(
	"fmt"
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
	cmd := exec.Command("./server", dir, strconv.Itoa(port))
	err := cmd.Start()
	if err != nil {
		fmt.Println("Something went wrong")
	}
	serverJobs <- ServerJob{port, dir, cmd}
	fmt.Printf("Started server for directory %s at port :%d\n", dir, port)
	err = cmd.Wait()
	fmt.Printf("Server at port :%d exited with error code %v\n", 
		port, err)
}

func initServers(serverJobs chan ServerJob)  {
	go startServer("/Users/robertcarney", 9090, serverJobs)
	go startServer("/Users/robertcarney", 9091, serverJobs)
}

func monitorServers(serverJobs chan ServerJob, quitChannel chan bool)  {
	for  {
		select {
		case <-quitChannel:
			fmt.Println("Recieved quit signal")
			close(quitChannel)
			return
		default:
			time.Sleep(1 * time.Second)
			fmt.Println("Looking for a server to monitor")
			currentJob := <- serverJobs
			fmt.Printf("Looking at server on port: %d\n", currentJob.port)
			if (currentJob.command == nil) {
				fmt.Println("Server process exited")
				continue
			}
			serverJobs <- currentJob
			fmt.Printf("Done monitoring server at port %d\n", currentJob.port)
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