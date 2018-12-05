package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/robcarney/gows/cache"
)

// Cache to be used
var fileCache *cache.Cache
// Cache duration
var cacheDuration time.Duration
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
			fileCache.Set(path, content, cacheDuration)
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
		relativePath := strings.Replace(path, baseDirectory, "", 1)
		if relativePath[:1] != "/"  {
			relativePath = "/" + relativePath
		}
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
	var port int
	var cacheDurationArg string
	var err error
	flag.IntVar(&port, "port", 9000, "The port for this server to listen on")
	flag.StringVar(&baseDirectory, "baseDirectory", "", "The base directory to serve files from")
	flag.StringVar(&baseDirectory, "baseDirectory", "", "The base directory to serve files from")
	flag.StringVar(&cacheDurationArg, "cacheDuration", "10s", "Duration for how long to keep files in cache")
	cacheDuration, err = time.ParseDuration(cacheDurationArg)
	if err != nil  {
		log.Fatal("Parse Duration ", err)
	}
	flag.Parse()
	fileCache = cache.NewCache()
	http.HandleFunc("/healthcheck", healthcheck)
	setUpFileEndpoints()
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("Listen and serve: ", err)
	}
}
