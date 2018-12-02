package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Workers []struct {
		Host     string `json:"host"`
		Port     int `json:"port"`
	} `json:"workers"`
}

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

func main() {
	f := "/Users/donaldhamnett/GolandProjects/gows/conf/config.json"
	c := LoadConfiguration(f)
	fmt.Print(c)
}
