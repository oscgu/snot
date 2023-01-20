package commands

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/oscgu/snot/internal/log"
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

		var text string
		dateTimeNow := time.Now()

		if config.Conf.Editor != "default" && config.Conf.Editor != "" {
			cmd := exec.Command(config.Conf.Editor, "/tmp/temp.snot")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}

			fileData, err := ioutil.ReadFile("/tmp/temp.snot")
			if err != nil {
				log.Fatal(err)
			}

			text = string(fileData)
			if err := os.Remove("/tmp/temp.snot"); err != nil {
				log.Fatal(err)
			}
		} else {
			var cancelled bool

			text, cancelled = editor.Create(topic, title, dateTimeNow)

			if cancelled {
				return
			}

		}

		n := note.Note{
			Topic:   topic,
			Title:   title,
			Author:  config.Conf.User.Name,
			Content: text,
			Created: dateTimeNow,
		}

		if d.Snotdb.Db.Model(n).Where("title = ?", n.Title).Where("topic = ?", n.Topic).Updates(note.Note{Content: n.Content, LastChanged: time.Now()}).RowsAffected == 0 {
			d.Snotdb.Db.Create(n)
			log.Info("Note created")
		} else {
			log.Info("Note upserted")
		}
	},
}

func init() {
	rootCmd.AddCommand(noteCmd)
}
