package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func healthcheck(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w, "Healthcheck succeeded")
}

func main() {
	directory := os.Args[1]
	port := ":" + os.Args[2]
	http.HandleFunc("/healthcheck", healthcheck)
	http.Handle("/", http.FileServer(http.Dir(directory)))
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Listen and serve: ", err)
	}
}
