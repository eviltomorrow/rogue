package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/eviltomorrow/rogue/lib/procutil"
	"github.com/eviltomorrow/rogue/lib/runutil"
	"github.com/spf13/cobra"
	"gopkg.in/natefinch/lumberjack.v2"
)

var startCommand = &cobra.Command{
	Use:   "start",
	Short: "Startup rogue-email app",
	Long:  `Startup rogue-email app`,
	Run: func(cmd *cobra.Command, args []string) {
		var writer = &lumberjack.Logger{
			Filename:   filepath.Join(runutil.ExecutableDir, "../log/error.log"),
			MaxSize:    20,
			MaxBackups: 10,
			MaxAge:     28,
			Compress:   true,
		}
		if err := procutil.RunInBackground(runutil.ExecutableName, []string{"--pid", pidFile}, nil, writer); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	},
}

func init() {
	rootCommand.AddCommand(startCommand)
}
