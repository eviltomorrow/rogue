package buildinfo

import (
	"fmt"
	"runtime"
)

var (
	MainVersion string
	GoVersion   = runtime.Version()
	GoOSArch    = runtime.GOOS + "/" + runtime.GOARCH
	GitSha      string
	GitTag      string
	GitBranch   string
	BuildTime   string
)

func GetVersion() string {
	return fmt.Sprintf(`Current Version: %s
Git Sha: %s
Git Tag: %s
Git Branch: %s
Go Version: %s
GO OS/Arch: %s
Build Time: %s`, MainVersion, GitSha, GitTag, GitBranch, GoVersion, GoOSArch, BuildTime)
}
