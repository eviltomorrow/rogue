package runutil

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/eviltomorrow/rogue/lib/util"
)

var (
	ExecutableName string
	ExecutableDir  string
	Pid            = os.Getpid()
	LaunchTime     = time.Now()
	HostName       string
	OS             = runtime.GOOS
	Arch           = runtime.GOARCH
	RunningTime    = func() string {
		return util.FormatDuration(time.Since(LaunchTime))
	}
	IP string
)

func init() {
	path, err := os.Executable()
	if err != nil {
		panic(fmt.Errorf("panic: get Executable path failure, nest error: %v", err))
	}
	path, err = filepath.Abs(path)
	if err != nil {
		panic(fmt.Errorf("panic: abs RootDir failure, nest error: %v", err))
	}
	ExecutableDir = filepath.Dir(path)
	ExecutableName = filepath.Base(path)

	name, err := os.Hostname()
	if err == nil {
		HostName = name
	}

	localIP, err := util.GetLocalIP2()
	if err == nil {
		IP = localIP
	}
}
