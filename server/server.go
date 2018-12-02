package main

import (
	"fmt"
	"github.com/robcarney/gows/cache"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Cache to be used
var fileCache *cache.Cache
// Base directory for the file server
var baseDirectory string

func healthcheck(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w, "Healthcheck succeeded")
}

// Endpoint handler for a file with the given path
func fileHandler(path string) func(w http.ResponseWriter, r *http.Request)  {
	return func(w http.ResponseWriter, r *http.Request)  {
		fmt.Printf("Trying to access file at %s\n", path)
		content := fileCache.Get(path)
		fmt.Printf("Cache try returned\n")
		if content != nil {
			fmt.Printf("Cache hit\n")
			w.Write(content)
		} else {
			content, err := ioutil.ReadFile(path)
			fmt.Printf("Content read\n")
			if err != nil {
				log.Printf("Could not read file at path %s\n", path)
				return
			}
			fileCache.Set(path, content, time.Second * 50)
			w.Write(content)
		}
	}
}

// Add an endpoint for a given file path relative to the base directory
func visitFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.IsDir() {
		relativePath := strings.Replace(path, baseDirectory, "/", 1)
		fmt.Printf("Adding endpoint at %s\n", relativePath)
		http.HandleFunc(relativePath, fileHandler(path))
	}
	return nil
}

// Add an endpoint for every file (not directories) in the 
//    base directory
func setUpFileEndpoints()  {
	err := filepath.Walk(baseDirectory, visitFile)
	if err != nil  {
		log.Fatal("setUpFileEndpoints: ", err)
	}
}

func main() {
	baseDirectory = os.Args[1]
	port := ":" + os.Args[2]
	fileCache = cache.NewCache()
	http.HandleFunc("/healthcheck", healthcheck)
	setUpFileEndpoints()
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Listen and serve: ", err)
	}
}