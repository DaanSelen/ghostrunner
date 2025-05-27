package main

import (
	"flag"
	"fmt"
	"ghostrunner-server/modules/database"
	"ghostrunner-server/modules/restapi"
	"ghostrunner-server/modules/timekeeper"
	"ghostrunner-server/modules/utilities"
	"log"
)

func main() {
	// Begin by trying to find the configuration file.
	cfgPtr := flag.String("conf", "./conf/ghostserver.conf", "Specify a config file location yourself. Relative to the program.")
	cfg := utilities.ReadConf(*cfgPtr)

	hmacKey, err := utilities.LoadHMACKey(cfg.TokenKeyFile)
	if err != nil {
		log.Println(utilities.ErrTag, err)
	}

	log.Println(utilities.InfoTag, "Starting the Sqlite3 database connection.")
	database.InitSqlite(cfg.AdminToken, hmacKey)

	log.Println(utilities.InfoTag, "Starting the API-Server backend.")
	restapi.InitApiServer(cfg, hmacKey)

	log.Println(utilities.InfoTag, "Components should have started.")
	log.Println(utilities.InfoTag, "Letting TimeKeeper take over...")
	log.Println(utilities.InfoTag, fmt.Sprintf("Interval set at: %d seconds.", cfg.Interval))

	timekeeper.KeepTime(cfg.Interval, cfg.PyVenvName)
}
