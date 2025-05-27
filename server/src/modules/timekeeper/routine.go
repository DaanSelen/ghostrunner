package timekeeper

import (
	"ghostrunner-server/modules/database"
	"ghostrunner-server/modules/utilities"
	"log"
)

func routine(cfg utilities.ConfigStruct, pyListArgs []string) {
	d := listDevices(cfg, pyListArgs) // Retrieve the Online devices.
	curTasks := database.RetrieveTasks()

	for index, task := range curTasks {
		relevantNodeids := task.Nodeids

		log.Printf("Processing Task %d", index)
		for _, nodeid := range relevantNodeids {
			if isNodeOnline(nodeid, d.OnlineDevices) {
				//result := wrapper.ExecCommand(nodeid, task.Command)
				log.Printf("Node online: %s", nodeid)
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
