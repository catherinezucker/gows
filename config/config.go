package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// ServerWorker is the data structure representing a single server worker
type ServerWorker struct {
	Host string `json:"host"`
	Port int `json:"port"`
}

// Config is the data structure representing the server configuration
type Config struct {
	Workers []ServerWorker `json:"workers"`
	BaseDirectory string `json:"baseDirectory"`
	MasterPort int `json:"masterPort"`
}

// LoadConfiguration loads configuration given the config file
func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
