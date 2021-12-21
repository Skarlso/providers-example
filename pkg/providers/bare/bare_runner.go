package bare

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/Skarlso/providers-example/pkg/providers"
)

// Config contains the configuration for this runner.
type Config struct {
}

// Dependencies any providers which this provider needs.
type Dependencies struct {
	Logger zerolog.Logger
	Store  providers.Storer
}

// Runner is a bare runner
type Runner struct {
	Config
	Dependencies
}

var _ providers.Runner = &Runner{}

// NewBareRunner creates a new Bare runner.
func NewBareRunner(cfg Config, deps Dependencies) *Runner {
	return &Runner{
		Dependencies: deps,
		Config:       cfg,
	}
}

// Run executes a bare metal plugin.
func (r *Runner) Run(ctx context.Context, name string, args []string) error {
	r.Logger.Info().Str("name", name).Strs("args", args).Msg("running bare metal plugin...")
	return nil
}
