package network

import (
	"fmt"

	"github.com/spf13/cobra"
)

// command variable needs to be exported
var printName = &cobra.Command{
	Use:   "printName", // camelCase
	Short: "printName",
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, print help
		name, _ := cmd.Flags().GetString("name")
		fmt.Printf("Hi There %s !\n", name)
	},
}

func init() {
	printName.Flags().String("name", "", "")
	HelloCmd.AddCommand(printName)
}
