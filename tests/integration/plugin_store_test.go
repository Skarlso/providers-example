package livestore

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/Skarlso/providers-example/pkg/models"
	"github.com/Skarlso/providers-example/pkg/providers/storer"
)

func TestPluginStore_MainFlow(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	l, err := storer.NewLiteStorer(logger, testDbLocation)
	assert.NoError(t, err)
	ctx := context.Background()
	// create a container plugin
	err = l.Create(ctx, &models.Plugin{
		Name: "test-bare-1",
		Type: models.Bare,
		Bare: &models.BareMetalPlugin{
			Location: "/tmp/plugins",
		},
	})
	assert.NoError(t, err)
	// create a bare metal plugin
	err = l.Create(ctx, &models.Plugin{
		Name: "test-container-1",
		Type: models.Container,
		Container: &models.ContainerPlugin{
			Image: "skarlso/container:v0.0.1",
		},
	})
	assert.NoError(t, err)

	// get
	p1, err := l.Get(ctx, "test-bare-1")
	assert.NoError(t, err)
	assert.True(t, p1.ID > 0)
	assert.Equal(t, "test-bare-1", p1.Name)
	assert.Equal(t, "/tmp/plugins", p1.Bare.Location)
	assert.Equal(t, models.Bare, p1.Type)
	p2, err := l.Get(ctx, "test-container-1")
	assert.NoError(t, err)
	assert.True(t, p2.ID > 0)
	assert.Equal(t, "test-container-1", p2.Name)
	assert.Equal(t, "skarlso/container:v0.0.1", p2.Container.Image)
	assert.Equal(t, models.Container, p2.Type)

	// list
	plugins, err := l.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, plugins, 2)
	assert.Equal(t, "test-bare-1", plugins[0].Name)
	assert.Equal(t, "test-container-1", plugins[1].Name)

	// delete
	err = l.Delete(ctx, "test-bare-1")
	assert.NoError(t, err)
	_, err = l.Get(ctx, "test-bare-1")
	assert.EqualError(t, err, "failed to run get: sql: no rows in result set")
}
