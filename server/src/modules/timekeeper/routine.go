package timekeeper

import (
	"fmt"
	"ghostrunner-server/modules/database"
	"ghostrunner-server/modules/utilities"
	"ghostrunner-server/modules/wrapper"
	"log"
	"strings"
)

func routine(venvName string, pyListArgs []string) {
	d := listDevices(venvName, pyListArgs) // Retrieve the Online devices.
	curTasks := database.RetrieveTasks()

	for index, task := range curTasks {
		relevantNodeid := task.Nodeid

		log.Printf("Processing Task %d", index)
		if isNodeOnline(relevantNodeid, d.OnlineDevices) {
			log.Printf("Node online: %s", relevantNodeid)
			result := forgeAndExec(venvName, relevantNodeid, task.Command) // Forge the Python command and execute it with the Python libraries.
			log.Println(result)

			//generateResult()

			log.Println("Removing Task from database...")
			database.RemoveTask(task.Name, task.Nodeid)
		} else {
			log.Printf("Node offline %s", relevantNodeid) // Just a debug line to tell the user that the node is offline.
		}
	}

}

func listDevices(venvName string, pyArgs []string) utilities.PyOnlineDevices {
	onDevices, err := wrapper.PyListOnline(venvName, pyArgs)
	if err != nil {
		log.Println(utilities.ErrTag, err)
	}

	return onDevices
}

func isNodeOnline(nodeid string, onlineDevices []utilities.Device) bool {
	for _, device := range onlineDevices {
		if device.NodeID == nodeid {
			return true
		}
	}
	return false
}

func forgeAndExec(venvName string, nodeid, command string) string {
	log.Printf("Triggered %s, on %s", command, nodeid)

	pyArgs := strings.Fields(fmt.Sprintf("--run --nodeid %s --command", nodeid))
	pyArgs = append(pyArgs, command)

	return wrapper.ExecTask(venvName, pyArgs)
}
