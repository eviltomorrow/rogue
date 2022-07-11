package main

import (
	"log"

	"github.com/eviltomorrow/rogue/app/collector/cmd"
	"github.com/eviltomorrow/rogue/lib/buildinfo"
)

var (
	MainVersion = "unknown"
	GitSha      = "unknown"
	GitTag      = "unknown"
	GitBranch   = "unknown"
	BuildTime   = "unknown"
)

func init() {
	buildinfo.MainVersion = MainVersion
	buildinfo.GitSha = GitSha
	buildinfo.GitTag = GitTag
	buildinfo.GitBranch = GitBranch
	buildinfo.BuildTime = BuildTime
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
