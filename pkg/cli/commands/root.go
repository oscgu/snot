package commands

import (
	"fmt"
	"os"

	config "github.com/oscgu/snot/pkg/cli/config"
	db "github.com/oscgu/snot/pkg/cli/snotdb"
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
		fmt.Fprintf(os.Stderr, "There was an error executing the command: %s", err)
		os.Exit(1)
	}
}
