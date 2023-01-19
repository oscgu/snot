package commands

import (
	"fmt"
	"strings"
	"time"

	config "github.com/oscgu/snot/pkg/cli/config"
	d "github.com/oscgu/snot/pkg/cli/dataproviders/snotdb"
	note "github.com/oscgu/snot/pkg/cli/note"
	editor "github.com/oscgu/snot/pkg/cli/note/ui/editor"
	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:                   "note [topic] (title)",
	Short:                 "Create a new note",
	Args:                  cobra.MinimumNArgs(2),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		topic := args[0]
		title := strings.Join(args[1:], " ")

		text, cancelled, created := editor.Create(topic, title, time.Now())

		if cancelled {
			return
		}

		n := note.Note{
			Topic:   topic,
			Title:   title,
			Author:  config.Conf.User.Name,
			Content: text,
			Created: created,
		}

		if d.Snotdb.Db.Model(n).Where("title = ?", n.Title).Updates(note.Note{Content: n.Content, LastChanged: time.Now()}).RowsAffected == 0 {
			d.Snotdb.Db.Create(n)
			fmt.Println("Note created.")
		} else {
			fmt.Println("Note updated.")
		}
	},
}

func init() {
	rootCmd.AddCommand(noteCmd)
}
