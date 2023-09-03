package danger

import (
	"os/exec"
	"runtime"
)

func Cmdstart(msg string) string {
	if isWindows() {
		newmsg := "chcp 65001 & " + msg
		output, err := execCmd([]string{"cmd", "/c", newmsg})
		if err != nil {
			return ""
		}
		return output
	}
	output, err := execCmd([]string{"sh", "-c", msg})
	if err != nil {
		return ""
	}
	return output
}

func execCmd(cmd []string) (string, error) {
	command := exec.Command(cmd[0], cmd[1:]...)
	output, err := command.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}
