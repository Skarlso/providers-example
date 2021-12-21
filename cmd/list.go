package cmd

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/Skarlso/providers-example/pkg/models"
	"github.com/Skarlso/providers-example/pkg/providers"
	"github.com/Skarlso/providers-example/pkg/providers/storer"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists all registered plugins.",
		Run:   runListCmd,
	}
	listArgs struct {
		_type string
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
	flag := listCmd.Flags()
	flag.StringVar(&listArgs._type, "type", "", "--type bare")
}

func runListCmd(cmd *cobra.Command, args []string) {
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
	results, err := store.List(context.Background(), providers.ListOpts{
		TypeFilter: listArgs._type,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to list plugins")
		os.Exit(1)
	}
	data := make([][]string, 0)
	for _, result := range results {
		d := []string{
			result.Name,
			result.Type,
		}
		if result.Type == models.Container {
			d = append(d, result.Container.Image)
		} else {
			d = append(d, result.Bare.Location)
		}
		data = append(data, d)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Type", "Image/Location"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}
