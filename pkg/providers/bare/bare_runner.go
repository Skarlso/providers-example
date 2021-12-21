package bare

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog"

	"github.com/Skarlso/providers-example/pkg/providers"
)

// Config contains the configuration for this runner.
type Config struct {
}

// Dependencies any providers which this provider needs.
type Dependencies struct {
	Logger zerolog.Logger
	Storer providers.Storer
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
	plugin, err := r.Storer.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("plugin not found: %w", err)
	}
	cmd := exec.Command(filepath.Join(plugin.Bare.Location, name), args...)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run plugin: %w", err)
	}
	fmt.Println(string(output))
	return nil
}
