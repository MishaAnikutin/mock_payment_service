package infrastructure

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func UpgradeHead(db *sql.DB) error {
	path := filepath.Join("database", "migrations.sql")
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %v", err)
	}

	queries := strings.Split(string(content), ";\n")

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	for i, q := range queries {
		q = strings.TrimSpace(q)
		if q == "" {
			continue
		}

		if _, err := tx.Exec(q); err != nil {
			return fmt.Errorf("migration %d failed: %v\nQuery: %s", i+1, err, q)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
