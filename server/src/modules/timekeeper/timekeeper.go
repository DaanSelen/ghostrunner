package timekeeper

import (
	"ghostrunner-server/modules/utilities"
	"log"
	"time"
)

func KeepTime(interval int) {
	transInterval := time.Duration(interval) * time.Second

	ticker := time.NewTicker(transInterval)
	defer ticker.Stop()

	for t := range ticker.C {
		log.Println(utilities.InfoTag, "Tick at:", t)
	}
}
