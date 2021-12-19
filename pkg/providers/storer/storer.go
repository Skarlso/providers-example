package storer

import (
	"context"
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
func NewLiteStorer(logger zerolog.Logger, location string) (*LiteStorer, error) {
	l := &LiteStorer{Logger: logger, DBLocation: location}
	if err := l.Init(); err != nil {
		return nil, err
	}
	return l, nil
}

var _ providers.Storer = &LiteStorer{}

// LiteStorer stores information in a SQLite backed storage medium.
type LiteStorer struct {
	Logger     zerolog.Logger
	DBLocation string
}

// Create will create a new entry in our storage.
func (l *LiteStorer) Create(ctx context.Context, plugin *models.Plugin) error {
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
	if _, err = db.Exec("insert into plugins(name, type, location) values($1, $2, $3);", plugin.Name, plugin.Type, plugin.Location); err != nil {
		return fmt.Errorf("failed to run insert into: %w", err)
	}
	l.Logger.Info().Str("name", plugin.Name).Msg("done")
	return nil
}

// Get returns plugin details.
func (l *LiteStorer) Get(ctx context.Context, name string) (*models.Plugin, error) {
	l.Logger.Info().Str("name", name).Msg("Getting plugin...")
	db, err := l.createConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			l.Logger.Error().Err(err).Msg("failed to close db connection")
		}
	}()
	// we could use a transaction here and all the jazz, but this is a blog post project. :)

	var (
		storedID       int
		storedName     string
		storedType     string
		storedLocation string
	)
	if err := db.QueryRow("select id, name, type, location from plugins where name = $1;", name).Scan(&storedID, &storedName, &storedType, &storedLocation); err != nil {
		return nil, fmt.Errorf("failed to run get: %w", err)
	}
	result := &models.Plugin{
		ID:       storedID,
		Name:     storedName,
		Type:     models.Type(storedType),
		Location: storedLocation,
	}
	return result, nil
}

func (l *LiteStorer) Delete(ctx context.Context, name string) error {
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
	l.Logger.Info().Str("name", name).Msg("done")
	return nil
}

func (l *LiteStorer) List(ctx context.Context) ([]*models.Plugin, error) {
	db, err := l.createConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			l.Logger.Error().Err(err).Msg("failed to close db connection")
		}
	}()
	// we could use a transaction here and all the jazz, but this is a blog post project. :)
	row, err := db.Query("select id, name, type, location from plugins")
	if err != nil {
		return nil, fmt.Errorf("failed to run query: %w", err)
	}
	var result []*models.Plugin
	for row.Next() {
		var (
			storedID       int
			storedName     string
			storedType     string
			storedLocation string
		)
		if err := row.Scan(&storedID, &storedName, &storedType, &storedLocation); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		result = append(result, &models.Plugin{
			ID:       storedID,
			Name:     storedName,
			Type:     models.Type(storedType),
			Location: storedLocation,
		})
	}
	return result, nil
}

func (l *LiteStorer) createConnection() (*sql.DB, error) {
	// check if db exist. If not, bootstrap it.
	db, err := sql.Open("sqlite3", filepath.Join(l.DBLocation, "provider.db"))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return db, nil
}

func (l *LiteStorer) Init() error {
	l.Logger.Debug().Str("location", l.DBLocation).Msg("Creating new database...®")
	if _, err := os.Stat(filepath.Join(l.DBLocation, "provider.db")); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat db file: %w", err)
	}
	db, err := sql.Open("sqlite3", filepath.Join(l.DBLocation, "provider.db"))
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	sqlStmt := `create table plugins (id integer primary key, name text, type text, location text);`
	if _, err := db.Exec(sqlStmt); err != nil {
		return fmt.Errorf("failed to execute bootstrap statement: %w", err)
	}
	return nil
}
