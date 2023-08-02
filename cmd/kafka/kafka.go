package kafka

import (
	"github.com/spf13/cobra"
)

// command variable needs to be exported
var Topic = &cobra.Command{
	Use: "topic", // camelCase
	// Short: "hello",
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, print help
		cmd.Help()
	},
}
