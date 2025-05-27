package timekeeper

import (
	"ghostrunner-server/modules/utilities"
	"log"
	"strings"
	"time"
)

var ( // Debugging
	pyListArgs = strings.Fields("-lo")
)

func KeepTime(interval int, venvName string) {
	transInterval := time.Duration(interval) * time.Second

	ticker := time.NewTicker(transInterval)
	defer ticker.Stop()

	for t := range ticker.C {
		log.Println(utilities.InfoTag, "Tick at:", t)
		log.Println(utilities.InfoTag, "Starting Routine.")
		routine(venvName, pyListArgs)
	}
}
