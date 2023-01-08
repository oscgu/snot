package cli

import (
	"strings"
	"time"

	config "github.com/oscgu/snot/pkg/cli/config"
	note "github.com/oscgu/snot/pkg/cli/note"
	ui "github.com/oscgu/snot/pkg/cli/note/ui/textarea"
	snotdb "github.com/oscgu/snot/pkg/cli/snotdb"
	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note [topic] (title)",
	Short: "Create a new note",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		title := strings.Join(args[1:], " ")
		text, cancelled, created := ui.TextArea(title)

		if cancelled {
			return
		}

		n := note.Note{
			Topic:   args[0],
			Title:   title,
			Author:  config.Conf.User.Name,
			Content: text,
			Created: created,
		}

		if snotdb.Db.Model(n).Where("title = ?", n.Title).Updates(note.Note{Content: n.Content, LastChanged: time.Now()}).RowsAffected == 0 {
			snotdb.Db.Create(n)
		}
	},
}

func init() {
	rootCmd.AddCommand(noteCmd)
}
