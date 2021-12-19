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
	err = l.Create(ctx, &models.Plugin{
		Name:     "test",
		Type:     models.Bare,
		Location: "/tmp/plugins",
	})
	assert.NoError(t, err)

	// get
	p1, err := l.Get(ctx, "test")
	assert.NoError(t, err)
	assert.True(t, p1.ID > 0)
	assert.Equal(t, "test", p1.Name)
	assert.Equal(t, "/tmp/plugins", p1.Location)
	assert.Equal(t, models.Bare, p1.Type)

	// list
	err = l.Create(ctx, &models.Plugin{
		Name:     "test2",
		Type:     models.Container,
		Location: "/tmp/plugins",
	})
	assert.NoError(t, err)
	plugins, err := l.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, plugins, 2)
	assert.Equal(t, "test", plugins[0].Name)
	assert.Equal(t, "test2", plugins[1].Name)

	// delete
	err = l.Delete(ctx, "test2")
	assert.NoError(t, err)
	_, err = l.Get(ctx, "test2")
	assert.EqualError(t, err, "failed to run get: sql: no rows in result set")
}
