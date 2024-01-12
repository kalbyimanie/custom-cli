package jenkins

import (
	"github.com/spf13/cobra"
)

// command variable needs to be exported
var JenkinsCmd = &cobra.Command{
	Use: "jenkins", // camelCase
	// Short: "hello",
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, print help
		cmd.Help()
	},
}
