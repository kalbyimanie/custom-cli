package network

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// command variable needs to be exported
var NetworkCmd = &cobra.Command{
	Use: "network", // camelCase
	Run: func(cmd *cobra.Command, args []string) {
		//NOTE check if --url flag exists
		if flag := cmd.Flags().Lookup("url"); flag == nil {
			fmt.Fprintf(os.Stderr, "[ERROR] need to set --url flag\n")
			cmd.Help()
			os.Exit(1)
		}
		//NOTE check if given flag is other than the expected --url flag
		if cmd.Flags().Changed("") {
			fmt.Fprintf(os.Stderr, "[ERROR] need to set --url flag\n")
			cmd.Help()
			os.Exit(1)
		}
	},
}

// command variable needs to be exported
var HelloCmd = &cobra.Command{
	Use:   "hello", // camelCase
	Short: "hello",
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, print help
		cmd.Help()
	},
}
