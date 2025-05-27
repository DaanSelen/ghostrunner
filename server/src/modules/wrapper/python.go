package wrapper

import (
	"encoding/json"
	"fmt"
	"ghostrunner-server/modules/utilities"
	"log"
	"os"
	"os/exec"
)

const (
	pyFile = "./../runner/runner.py"
)

func PyListOnline(venvName string, pyArgs []string) (utilities.PyOnlineDevices, error) {
	pyBin := fmt.Sprintf("./../runner/%s/bin/python", venvName)
	runtimeArgs := append([]string{pyFile}, pyArgs...)

	cmd := exec.Command(pyBin, runtimeArgs...)

	rawData, err := cmd.CombinedOutput()
	log.Println(string(rawData))
	if err != nil {
		cwd, _ := os.Getwd()
		return utilities.PyOnlineDevices{}, fmt.Errorf("python execution failed, working directory: %s", cwd)
	}

	var data utilities.PyOnlineDevices
	if err := json.Unmarshal(rawData, &data); err != nil {
		return utilities.PyOnlineDevices{}, fmt.Errorf("error unmarshaling: %v", err)
	}

	return data, nil
}
