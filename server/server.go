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

var fileCache cache.Cache
var baseDirectory string

func healthcheck(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w, "Healthcheck succeeded")
}

func fileHandler(path string) func(w http.ResponseWriter, r *http.Request)  {
	return func(w http.ResponseWriter, r *http.Request)  {
		content := fileCache.Get(path)
		if content != nil {
			w.Write(content)
		} else {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				log.Printf("Could not read file at path %s\n", path)
				return
			}
			fileCache.Set(path, content, time.Second * 50)
			w.Write(content)
		}
	}
}

func visitFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.IsDir() {
		relativePath := strings.Replace(path, baseDirectory, "", 1)
		http.HandleFunc(relativePath, fileHandler(path))
	}
	return nil
}

func setUpFileEndpoints()  {
	err := filepath.Walk(baseDirectory, visitFile)
	if err != nil  {
		log.Fatal("setUpFileEndpoints: ", err)
	}
}

func main() {
	baseDirectory = os.Args[1]
	port := ":" + os.Args[2]
	http.HandleFunc("/healthcheck", healthcheck)
	setUpFileEndpoints()
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Listen and serve: ", err)
	}
}