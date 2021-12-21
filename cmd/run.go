package cmd

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/Skarlso/providers-example/pkg/providers/bare"
	"github.com/Skarlso/providers-example/pkg/providers/container"
	"github.com/Skarlso/providers-example/pkg/providers/storer"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run a plugin.",
		Run:   runRunCmd,
	}
	runArgs struct {
		name string
		args []string
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	flag := runCmd.Flags()
	flag.StringVar(&runArgs.name, "name", "", "--name")
	flag.StringSliceVar(&runArgs.args, "args", nil, "--args")
}

func runRunCmd(cmd *cobra.Command, args []string) {
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
	barePlugin := bare.NewBareRunner(bare.Config{}, bare.Dependencies{
		Logger: log,
		Storer: store,
	})
	containerPlugin, err := container.NewRunner(container.Config{
		DefaultMaximumCommandRuntime: 15,
	}, container.Dependencies{
		Storer: store,
		Next:   barePlugin,
		Logger: log,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create container runner")
		os.Exit(1)
	}
	if err := containerPlugin.Run(context.Background(), runArgs.name, runArgs.args); err != nil {
		log.Error().Err(err).Msg("Failed to run plugin")
		os.Exit(1)
	}
	log.Info().Msg("All done.")
}
