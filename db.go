package vcd

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/noonat/vcd/migrate"
	"github.com/pkg/errors"
)

var migrations = []migrate.Migration{
	{
		Comment: "Create vessels, vessel_clicks, and vessel_pilot_clicks tables",
		Up: migrate.ExecQueries([]migrate.MigrationQuery{
			{
				Comment: "creating vessels table",
				Query: `
					CREATE TABLE vessels (
						id INT UNSIGNED NOT NULL AUTO_INCREMENT,
						created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
						ip VARBINARY(16) NOT NULL,
						cfe TEXT NOT NULL,
						data TEXT NOT NULL,
						PRIMARY KEY (id)
					) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
				`,
			},
			{
				Comment: "creating vessel_clicks table",
				Query: `
					CREATE TABLE vessel_clicks (
						id INT UNSIGNED NOT NULL AUTO_INCREMENT,
						created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
						ip VARBINARY(16) NOT NULL,
						referrer VARCHAR(1024) NOT NULL,
						vessel_id INT UNSIGNED NOT NULL,
						PRIMARY KEY (id)
					) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
				`,
			},
			{
				Comment: "creating vessel_pilot_clicks table",
				Query: `
					CREATE TABLE vessel_pilot_clicks (
						id INT UNSIGNED NOT NULL AUTO_INCREMENT,
						created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
						ip VARBINARY(16) NOT NULL,
						referrer VARCHAR(1024) NOT NULL,
						vessel_id INT UNSIGNED NOT NULL,
						PRIMARY KEY (id)
					) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
				`,
			},
		}),
		Down: migrate.ExecQueries([]migrate.MigrationQuery{
			{
				Comment: "dropping vessel_pilot_clicks table",
				Query:   `DROP TABLE vessel_pilot_clicks`,
			},
			{
				Comment: "dropping vessel_clicks table",
				Query:   `DROP TABLE vessel_clicks`,
			},
			{
				Comment: "dropping vessels table",
				Query:   `DROP TABLE vessels`,
			},
		}),
	},
}

var (
	retryAttempts = 5
	retryDuration = 1 * time.Second
)

// openDB opens a database connection and runs migrations on it.
func openDB(ctx context.Context, dsn string) (*sql.DB, error) {
	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing dsn")
	}
	if config.Params == nil {
		config.Params = map[string]string{}
	}
	config.Params["collation"] = "utf8mb4_unicode_ci"
	config.Params["sql_mode"] = "'" + strings.Join([]string{
		"ERROR_FOR_DIVISION_BY_ZERO",
		"NO_ENGINE_SUBSTITUTION",
		"NO_ZERO_DATE",
		"NO_ZERO_IN_DATE",
		"ONLY_FULL_GROUP_BY",
		"STRICT_TRANS_TABLES",
	}, ",") + "'"
	config.ParseTime = true
	dsn = config.FormatDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to db")
	}
	for i := 0; i < retryAttempts; i++ {
		err = db.PingContext(ctx)
		if err == nil {
			break
		}
		log.Printf("Error pinging database, retrying after %f seconds...", retryDuration.Seconds())
		time.Sleep(retryDuration)
	}
	if err != nil {
		return nil, errors.Wrap(err, "error pinging db")
	}

	if err := migrate.Migrate(ctx, db, migrations); err != nil {
		return nil, errors.Wrap(err, "error migrating db")
	}

	return db, nil
}
