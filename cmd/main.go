package main

import (
	"os"

	"github.com/spf13/cobra"
)

var swapCmd = &cobra.Command{
	Use:   "swap",
	Short: "swap is a commond line client for bytom contract",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			os.Exit(0)
		}
	},
}

func main() {
	swapCmd.AddCommand(deployCmd)
	swapCmd.AddCommand(callCmd)
	swapCmd.AddCommand(deployHTLCCmd)
	swapCmd.AddCommand(callHTLCCmd)
	if err := swapCmd.Execute(); err != nil {
		os.Exit(0)
	}
}
