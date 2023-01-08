package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "snot",
	Short: "snot - a simple cli tool to take notes",
	Long:  "snot - is a simple cli tool to take notes",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "There was an error executing the command: %s", err)
		os.Exit(1)
	}
}
