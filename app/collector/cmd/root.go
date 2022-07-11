package cmd

import (
	"fmt"
	"os"

	"github.com/eviltomorrow/rogue/lib/buildinfo"
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "rogue-collector",
	Short: "Collect stock trade data from net126 or sina",
	Long:  "Collect stock trade data from net126 or sina",
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			fmt.Println(buildinfo.GetVersion())
			os.Exit(0)
		}
		cmd.Usage()
	},
}

var (
	version bool
)

func init() {
	rootCommand.Flags().BoolVarP(&version, "version", "v", false, "version about rogue-collector")
}

func Execute() error {
	rootCommand.CompletionOptions.HiddenDefaultCmd = true

	return rootCommand.Execute()
}
