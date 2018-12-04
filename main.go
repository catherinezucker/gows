package main


import (
	"os"
	"log"

	"github.com/robcarney/gows/master"
	"github.com/robcarney/gows/config"
)

func main()  {
	log.SetOutput(os.Stderr)
	serverConfig, err := config.LoadConfiguration("conf/config.json")
	if err != nil {
		log.Fatal("Load Configuration: ", err)
	}
	master.Run(serverConfig)
}