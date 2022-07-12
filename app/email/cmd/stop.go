package cmd

import (
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/eviltomorrow/rogue/lib/procutil"
	"github.com/eviltomorrow/rogue/lib/runutil"
	"github.com/spf13/cobra"
)

var stopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stop rogue-email app",
	Long:  `Stop rogue-email app`,
	Run: func(cmd *cobra.Command, args []string) {
		var pidFile = filepath.Join(runutil.ExecutableDir, "../var/run/rogue-email.pid")
		process, err := procutil.FindWithPidFile(pidFile)
		if err != nil {
			log.Printf("[E] Find process with pidFile failure, nest error: %v", err)
			os.Exit(1)
		} else {
			if err := process.Signal(syscall.SIGQUIT); err != nil {
				log.Printf("[E] Signal process failure, nest error: %v", err)
			}
		}
	},
}

func init() {
	rootCommand.AddCommand(stopCommand)
}
