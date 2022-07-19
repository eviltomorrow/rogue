package cmd

import (
	"fmt"

	"github.com/eviltomorrow/rogue/lib/buildinfo"
	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print version about rogue-email",
	Long:  `The version about rogue-email`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println(buildinfo.GetVersion())
	},
}

func init() {
	rootCommand.AddCommand(versionCommand)
}
