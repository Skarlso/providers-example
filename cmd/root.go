package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "executor",
		Short: "The main entry point to the plugin running system.",
		Run:   runRootCmd,
	}
	rootArgs struct {
	}
)

func runRootCmd(cmd *cobra.Command, args []string) {
	fmt.Println("Running executor")
}

func Execute() error {
	return rootCmd.Execute()
}
