package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "executor",
		Short: "The main entry point to the plugin running system.",
	}
	rootArgs struct {
		location string
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&rootArgs.location, "location", "", "--location /~.config/providers")
	if rootArgs.location == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Failed to get home folder: ", err)
			os.Exit(1)
		}
		homeFolder := filepath.Join(home, ".config", "providers")
		if err := os.MkdirAll(homeFolder, 0766); err != nil {
			fmt.Println("Failed to create config folder: ", err)
			os.Exit(1)
		}
		rootArgs.location = homeFolder
	}
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
