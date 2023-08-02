package cmd

import (
	"custom-cli/cmd/kafka"
	"custom-cli/cmd/network"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "custom-cli",
	Short: "A sample CLI application using Cobra",
	Long:  `A sample CLI application demonstrating how to use Cobra for command-line applications.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, print help
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(network.NetworkCmd)
	rootCmd.AddCommand(network.HelloCmd)
	rootCmd.AddCommand(kafka.Topic)
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
