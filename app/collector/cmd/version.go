package cmd

import (
	"fmt"

	"github.com/eviltomorrow/rogue/lib/buildinfo"
	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print version about rogue-collector",
	Long:  `The version about rogue-collecotr`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(buildinfo.GetVersion())
	},
}

func init() {
	rootCommand.AddCommand(versionCommand)
}
