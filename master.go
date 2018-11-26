package main

import(
	"fmt"
	"log"
	"os/exec"
)

func startServer(dir, port string)  {
	for true  {
		cmd := exec.Command("./server", dir, port)
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Started server for directory %s at port :%s\n", dir, port)
		err = cmd.Wait()
		fmt.Println("Server at port :%s exited with error code %v, attempting to start again", 
			port, err)
	}
}

func main()  {
	go startServer("/Users/robertcarney/tmp", "9000")
	go startServer("/Users/robertcarney/tmp", "9001")
	for true {  }
}