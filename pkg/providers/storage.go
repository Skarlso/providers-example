package providers

import (
	"context"

	"github.com/Skarlso/providers-example/pkg/models"
)

// ListOpts defines options for listing plugins.
type ListOpts struct {
	TypeFilter string
}

// Storer can store information about the plugins that were created.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_storer_client.go . Storer
type Storer interface {
	Init() error
	Create(ctx context.Context, plugin *models.Plugin) error
	Get(ctx context.Context, name string) (*models.Plugin, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context, opts ListOpts) ([]*models.Plugin, error)
}
