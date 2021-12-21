package cmd

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/Skarlso/providers-example/pkg/models"
	"github.com/Skarlso/providers-example/pkg/providers/storer"
)

var (
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds a new plugin to the list of plugins.",
		Run:   runAddCmd,
	}
	addArgs struct {
		_type    string
		name     string
		location string
		image    string
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
	flag := addCmd.Flags()
	flag.StringVar(&addArgs._type, "type", models.Bare, "--type bare")
	flag.StringVar(&addArgs.name, "name", "", "--name bare")
	flag.StringVar(&addArgs.location, "file-location", "", "--file-location ~/.config/providers/")
	flag.StringVar(&addArgs.image, "image", "", "--image skarlso/providers:echo-v1")
}

func runAddCmd(cmd *cobra.Command, args []string) {
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
	plugin := &models.Plugin{
		Name: addArgs.name,
		Type: addArgs._type,
	}
	if addArgs._type == models.Container {
		plugin.Container = &models.ContainerPlugin{
			Image: addArgs.image,
		}
	} else if addArgs._type == models.Bare {
		plugin.Bare = &models.BareMetalPlugin{
			Location: addArgs.location,
		}
	} else {
		log.Error().Str("type", addArgs._type).Msg("Invalid type.")
		os.Exit(1)
	}
	if err := store.Create(context.Background(), plugin); err != nil {
		log.Error().Err(err).Msg("Failed to add plugin")
		os.Exit(1)
	}
}
