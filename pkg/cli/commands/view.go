package commands

import (
	"github.com/oscgu/snot/pkg/cli/note/ui/viewer"
	"github.com/oscgu/snot/pkg/cli/snotdb"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:                   "view",
	Short:                 "Browse through all your notes",
	Args:                  cobra.NoArgs,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		var topics []string
		snotdb.Db.Table("notes").Distinct("topic").Scan(&topics)

		viewer.Create(topics)
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
