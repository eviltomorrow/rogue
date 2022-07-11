package cmd

import (
	"fmt"

	"github.com/eviltomorrow/rogue/lib/buildinfo"
	"github.com/spf13/cobra"
)

var stopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stop rogue-email app",
	Long:  `Stop rogue-email app`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(buildinfo.GetVersion())
	},
}

func init() {
	rootCommand.AddCommand(stopCommand)
}
