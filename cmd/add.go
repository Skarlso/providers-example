package cmd

import "github.com/spf13/cobra"

var (
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds a new plugin to the list of plugins.",
		Run:   runAddCmd,
	}
	addArgs struct {
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAddCmd(cmd *cobra.Command, args []string) {}
