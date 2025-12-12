package salmonstack

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func openSQLite() (*sql.DB, error) {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("openSQLite: DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := applyPragmas(db); err != nil {
		return nil, fmt.Errorf("apply pragmas: %w", err)
	}

	return db, nil
}

// Sets all required SQLite PRAGMAS
func applyPragmas(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA foreign_keys = ON;",
		"PRAGMA strict = ON;",
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA journal_size_limit = 67108864;", // 64MB
		"PRAGMA mmap_size = 134217728;",         // 128MB
		"PRAGMA cache_size = 2000;",
		"PRAGMA busy_timeout = 5000;",
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return fmt.Errorf("apply %q: %w", p, err)
		}
	}

	return nil
}
