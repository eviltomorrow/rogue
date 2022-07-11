package procutil

import (
	"io/ioutil"
	"os"
	"strconv"
)

func FindWithPidFile(pidFile string) (*os.Process, error) {
	buf, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return nil, err
	}
	pid, err := strconv.Atoi(string(buf))
	if err != nil {
		return nil, err
	}
	return os.FindProcess(pid)
}

func FindWithPid(pid int) (*os.Process, error) {
	return os.FindProcess(pid)
}
