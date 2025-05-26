package wrapper

import (
	"encoding/json"
	"fmt"
	"ghostrunner-server/modules/utilities"
	"log"
	"os"
	"os/exec"
)

func PyListOnline(venvName string) {
	pythonBin := fmt.Sprintf("./../runner/%s/bin/python", venvName)
	cmd := exec.Command(pythonBin, "./../runner/runner.py", "-lo")

	data, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(utilities.ErrTag, err, data)
		cwd, _ := os.Getwd()
		log.Println("Working directory:", cwd)
		return
	}

	var status utilities.PyOnlineDevices
	if err := json.Unmarshal(data, &status); err != nil {
		fmt.Println("Error unmarshaling:", err)
		return
	}

	fmt.Printf("Parsed Struct: %+v\n", status)
}
