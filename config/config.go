package config

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// ServerWorker is the data structure representing a single server worker
type ServerWorker struct {
	Host string `json:"host"`
	Port int `json:"port"`
	CacheDuration string `json:"cacheDuration"`
}

// Config is the data structure representing the server configuration
type Config struct {
	Workers []ServerWorker `json:"workers"`
	BaseDirectory string `json:"baseDirectory"`
	MasterPort int `json:"masterPort"`
}

// LoadConfiguration loads configuration given the config file
func LoadConfiguration(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Fatal("Load Configuration: ", err)
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	for _, w := range config.Workers  {
		_, err := time.ParseDuration(w.CacheDuration)
		if err != nil {
			log.Fatal("Parse Duration: ", err)
		}
	}
	return config, err
}
