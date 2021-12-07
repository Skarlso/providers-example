package storer

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"

	"github.com/Skarlso/providers-example/pkg/models"
	"github.com/Skarlso/providers-example/pkg/providers"
)

// NewLiteStorer creates a storer provider.
func NewLiteStorer(logger zerolog.Logger) *LiteStorer {
	return &LiteStorer{
		Logger: logger,
	}
}

var _ providers.Storer = &LiteStorer{}

// LiteStorer stores information in a SQLite backed storage medium.
type LiteStorer struct {
	Logger     zerolog.Logger
	DBLocation string
}

// Create will create a new entry in our storage.
func (l *LiteStorer) Create(plugin *models.Plugin) error {
	l.Logger.Info().Str("name", plugin.Name).Msg("Creating new plugin...")
	db, err := l.createConnection()
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			l.Logger.Error().Err(err).Msg("failed to close db connection")
		}
	}()
	// we could use a transaction here and all the jazz, but this is a blog post project. :)
	if _, err = db.Exec("insert into plugins(name, type) values($1, $2);", plugin.Name, plugin.Type); err != nil {
		return fmt.Errorf("failed to run insert into: %w", err)
	}
	l.Logger.Info().Str("name", plugin.Name).Msg("done")
	return nil
}

// Get returns plugin details.
func (l *LiteStorer) Get(name string) (*models.Plugin, error) {
	return nil, nil
}

func (l *LiteStorer) Delete(name string) error {
	l.Logger.Info().Str("name", name).Msg("Deleting plugin...")
	db, err := l.createConnection()
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			l.Logger.Error().Err(err).Msg("failed to close db connection")
		}
	}()
	// we could use a transaction here and all the jazz, but this is a blog post project. :)
	if _, err = db.Exec("delete from plugins where name = $1;", name); err != nil {
		return fmt.Errorf("failed to run insert into: %w", err)
	}
	l.Logger.Info().Str("name", plugin.Name).Msg("done")
	return nil
}

func (l *LiteStorer) List() ([]*models.Plugin, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LiteStorer) createConnection() (*sql.DB, error) {
	// check if db exist. If not, bootstrap it.
	if _, err := os.Stat(filepath.Join(l.DBLocation, "provider.db")); os.IsNotExist(err) {
		if err := l.bootstrapStore(); err != nil {
			return nil, fmt.Errorf("failed to bootstrap database: %w", err)
		}
	}
	db, err := sql.Open("sqlite3", filepath.Join(l.DBLocation, "provider.db"))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return db, nil
}

func (l *LiteStorer) bootstrapStore() error {
	db, err := sql.Open("sqlite3", filepath.Join(l.DBLocation, "provider.db"))
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	sqlStmt := `create table plugins (id integer primary key, name text, type text);`
	if _, err := db.Exec(sqlStmt); err != nil {
		return fmt.Errorf("failed to execute bootstrap statement: %w", err)
	}
	return nil
}
