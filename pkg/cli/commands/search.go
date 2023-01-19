package commands

import (
	"fmt"

	d "github.com/oscgu/snot/pkg/cli/dataproviders/snotdb"
	"github.com/oscgu/snot/pkg/cli/note"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:                   "search [title|content] (string)",
	Short:                 "Search for a string thats either in the title or content",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var notes []note.Note
		d.Snotdb.Db.Where(args[0]+" LIKE ?", "%"+args[1]+"%").Find(&notes)
		fmt.Println(notes)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
