// Package migrate provides helpers for running SQL database migrations. It's
// designed for migrations that are specified in code and distributed as part
// of the application binary, and applied as part of the application startup
// (rather than via external files and an external tool).
package migrate

import (
	"context"
	"database/sql"
	"log"

	"github.com/pkg/errors"
)

// MigrationFunc is type of function used for the up and down migrations.
type MigrationFunc func(ctx context.Context, db *sql.DB) error

// Migration represents an individual migration step. The Up function is run
// to migrate from the previous version to this version, and the Down function
// can be run to go back the other way. The Comment is inserted into the
// schema_versions table after migrating to this version.
type Migration struct {
	Comment string
	Up      MigrationFunc
	Down    MigrationFunc
}

// MigrationQuery specifies a comment and a SQL query for a migration. This
// is used with the ExecQueries helper to generate a migration step from raw
// SQL queries.
type MigrationQuery struct {
	Comment string
	Query   string
}

// ExecQueries generates a migration function from a list of SQL queries.
// Running the returned function will execute each of the SQL queries as its
// migration step.
func ExecQueries(queries []MigrationQuery) MigrationFunc {
	return func(ctx context.Context, db *sql.DB) error {
		for _, q := range queries {
			_, err := db.ExecContext(ctx, q.Query)
			if err != nil {
				return errors.Errorf("error %s: %s", q.Comment, err)
			}
		}
		return nil
	}
}

// Migrate upgrades the given database to the latest migration in the list
// of passed migrations. It will use the schema_versions table to track the
// migrations that have been run (and it will create that table if it doesn't
// already exist).
func Migrate(ctx context.Context, db *sql.DB, migrations []Migration) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_versions (
			version INT UNSIGNED NOT NULL PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			comment TEXT NOT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`)
	if err != nil {
		return errors.Wrap(err, "error creating schema_versions table")
	}

	var currentVersion int
	row := db.QueryRowContext(ctx, `SELECT COALESCE(MAX(version), 0) FROM schema_versions`)
	if err := row.Scan(&currentVersion); err != nil {
		return errors.Wrap(err, "error querying current schema version")
	}

	for i, m := range migrations {
		version := i + 1
		if version <= currentVersion {
			continue
		}
		log.Printf("Upgrading database to version %d", version)
		if err := m.Up(ctx, db); err != nil {
			return errors.Wrapf(err, "error upgrading database to version %d", version)
		}
		_, err := db.ExecContext(ctx, `
			INSERT INTO schema_versions (version, comment) VALUES (?, ?)
		`, version, m.Comment)
		if err != nil {
			return errors.Wrapf(err, "error inserting schema_versions row for version %d", version)
		}
	}

	return nil
}
