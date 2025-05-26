package wrapper

import (
	"encoding/json"
	"fmt"
	"ghostrunner-server/modules/utilities"
	"os"
	"os/exec"
)

const (
	pyFile = "./../runner/runner.py"
)

func PyListOnline(venvName string) (utilities.PyOnlineDevices, error) {
	pyBin := fmt.Sprintf("./../runner/%s/bin/python", venvName)
	cmd := exec.Command(pyBin, pyFile, "-lo")

	rawData, err := cmd.CombinedOutput()
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
