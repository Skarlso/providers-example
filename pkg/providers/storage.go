package providers

import "github.com/Skarlso/providers-example/pkg/models"

// Storer can store information about the plugins that were created.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_storer_client.go . Storer
type Storer interface {
	Create(plugin *models.Plugin) error
	Get(name string) (*models.Plugin, error)
	Delete(name string) error
	List() ([]*models.Plugin, error)
}
