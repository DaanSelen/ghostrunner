package timekeeper

import (
	"ghostrunner-server/modules/utilities"
	"ghostrunner-server/modules/wrapper"
	"log"
	"strings"
	"time"
)

var ( // Debugging
	pyListArgs = strings.Fields("-lo")
)

func KeepTime(interval int, cfg utilities.ConfigStruct) {
	transInterval := time.Duration(interval) * time.Second

	ticker := time.NewTicker(transInterval)
	defer ticker.Stop()

	for t := range ticker.C {
		log.Println(utilities.InfoTag, "Tick at:", t)
		log.Println(utilities.InfoTag, "Starting Routine.")
		routine(cfg, pyListArgs)
	}
}

func listDevices(cfg utilities.ConfigStruct, pyArgs []string) utilities.PyOnlineDevices {
	onDevices, err := wrapper.PyListOnline(cfg.PyVenvName, pyArgs)
	if err != nil {
		log.Println(utilities.ErrTag, err)
	}

	return onDevices
}
