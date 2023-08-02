package network

import (
	"fmt"

	"github.com/spf13/cobra"
)

// command variable needs to be exported
var helloCmd = &cobra.Command{
	Use:   "hello", // camelCase
	Short: "hello",
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, print help
		fmt.Println("hello cmd")
	},
}

func init() {
	HelloCmd.AddCommand(helloCmd)
}
