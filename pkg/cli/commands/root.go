package commands

import (
	"github.com/oscgu/snot/internal/log"
	config "github.com/oscgu/snot/pkg/cli/config"
	db "github.com/oscgu/snot/pkg/cli/dataproviders/snotdb"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "snot",
	Short: "A simple cli tool to take notes and view them in the terminal",
}

func Execute() {
	db.Init()
	config.Init()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
