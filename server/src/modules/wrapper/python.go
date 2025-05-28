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

func pyExec(venvName string, pyArgs []string) ([]byte, error) {
	pyBin := fmt.Sprintf("./../runner/%s/bin/python", venvName)
	runtimeArgs := append([]string{pyFile}, pyArgs...)

	cmd := exec.Command(pyBin, runtimeArgs...)

	return cmd.CombinedOutput()
}

func PyListOnline(venvName string, pyArgs []string) (utilities.PyOnlineDevices, error) {
	rawData, err := pyExec(venvName, pyArgs)
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

func ExecTask(venvName string, pyArgs []string) string {
	rawData, err := pyExec(venvName, pyArgs)
	if err != nil {
		cwd, _ := os.Getwd()
		log.Println("FAILED,", err, "CWD:", cwd)
	}

	return string(rawData)
}
