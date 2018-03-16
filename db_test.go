package vcd

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func requireTestDB(t *testing.T, fixtures []string) *sql.DB {
	db, err := openDB(context.Background(), "root:@tcp(127.0.0.1:3306)/test_vcd")
	require.NoError(t, err)

	// truncate the tables
	tables := []string{
		"vessels",
		"vessel_clicks",
		"vessel_pilot_clicks",
	}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table))
		require.NoError(t, err)
	}

	// load data from fixtures
	for _, fn := range fixtures {
		b, err := ioutil.ReadFile(filepath.Join("testdata", fn))
		require.NoError(t, err)
		_, err = db.Exec(string(b))
		require.NoError(t, err)
	}

	return db
}
