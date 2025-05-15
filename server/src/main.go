package main

import (
	"flag"
	"fmt"
	"ghostrunner-server/modules/confread"
	"ghostrunner-server/modules/restapi"
	"log"
)

func main() {
	// Begin by trying to find the configuration file.
	confPtr := flag.String("conf", "./conf/ghostserver.conf", "Specify a config file location yourself. Relative to the program.")
	config := confread.ReadConf(*confPtr)

	log.Println("Starting the API-Server backend.")
	restapi.InitApiServer(config)

	fmt.Scanln()
}
