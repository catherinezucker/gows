package main

import (
	"fmt"
	"log"
	"net/http"
)

func healthcheck(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w, "Healthcheck succeeded")
}

func main() {
	directory := "/Users/robertcarney/tmp/server-test";
	http.HandleFunc("/healthcheck", healthcheck)
	http.Handle("/", http.FileServer(http.Dir(directory)))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Listen and serve: ", err)
	}
}