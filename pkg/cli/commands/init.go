package commands

import (
	config "github.com/oscgu/snot/pkg/cli/config"
	db "github.com/oscgu/snot/pkg/cli/snotdb"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:                   "init",
	Short:                 "Create all the necessary files for using the cli",
	Args:                  cobra.NoArgs,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		db.Init()
		config.Init()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
