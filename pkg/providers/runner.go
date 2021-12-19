package providers

import "context"

// Runner runs a plugin.
type Runner interface {
	Run(ctx context.Context, name string, args []string) error
}
