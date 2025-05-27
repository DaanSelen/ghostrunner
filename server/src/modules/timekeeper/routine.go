package timekeeper

import (
	"ghostrunner-server/modules/database"
	"ghostrunner-server/modules/utilities"
	"log"
)

func routine(cfg utilities.ConfigStruct, pyListArgs []string) {
	d := listDevices(cfg, pyListArgs) // Retrieve the Online devices.
	curTasks := database.RetrieveTasks()

	for index, device := range d.OnlineDevices {
		log.Println(index, device)
	}

	for index, task := range curTasks {
		log.Println(index, task)
		for _, nodeid := range task.Nodeids {
			if isNodeOnline(nodeid, d.OnlineDevices) {
				log.Printf("NodeID %s is online\n", nodeid)
			}
		}
	}

}

func isNodeOnline(nodeid string, onlineDevices []utilities.Device) bool {
	for _, device := range onlineDevices {
		if device.NodeID == nodeid {
			return true
		}
	}
	return false
}
