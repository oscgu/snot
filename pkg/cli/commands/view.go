package commands

import (
	"github.com/oscgu/snot/pkg/cli/note/ui/viewer"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:                   "view",
	Short:                 "Browse through all your notes",
	Args:                  cobra.NoArgs,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		viewer.Create()
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
