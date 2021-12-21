package cmd

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/Skarlso/providers-example/pkg/providers/storer"
)

var (
	removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove a registered plugin.",
		Run:   runRemoveCmd,
	}
	removeArgs struct {
		name string
	}
)

func init() {
	rootCmd.AddCommand(removeCmd)
	flag := removeCmd.Flags()
	flag.StringVar(&removeArgs.name, "name", "", "--name bare")
}

func runRemoveCmd(cmd *cobra.Command, args []string) {
	out := zerolog.ConsoleWriter{
		Out: os.Stderr,
	}
	log := zerolog.New(out).With().
		Timestamp().
		Logger()

	store, err := storer.NewLiteStorer(log, rootArgs.location)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialise storer")
		os.Exit(1)
	}
	if err := store.Delete(context.Background(), removeArgs.name); err != nil {
		log.Error().Err(err).Msg("Failed to remove plugin")
		os.Exit(1)
	}
}
