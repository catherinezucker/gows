package main


import (
	"os"
	"log"

	"github.com/robcarney/gows/master"
	"github.com/robcarney/gows/config"
)

func main()  {
	log.SetOutput(os.Stderr)
	serverConfig := config.LoadConfiguration("conf/config.json")
	master.Run(serverConfig)
}