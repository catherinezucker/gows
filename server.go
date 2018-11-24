package main

import (
	"log"
	"net/http"
)

func main() {
	directory := "/Users/robertcarney/tmp/server-test";
	err := http.ListenAndServe(":8080", http.FileServer(http.Dir(directory)))
	if err != nil {
		log.Fatal("Listen and serve: ", err)
	}
}