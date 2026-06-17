package sqllite

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations runs all pending migrations
func RunMigrations(db *sql.DB) error {
	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		return err
	}

	// Read all migration files
	migrationFiles, err := readMigrationFiles()
	if err != nil {
		return err
	}

	// Get applied migrations
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	// Run pending migrations
	for _, migration := range migrationFiles {
		if applied[migration.name] {
			log.Printf("Migration %s already applied, skipping...", migration.name)
			continue
		}

		log.Printf("Running migration: %s", migration.name)

		if err := runMigration(db, migration); err != nil {
			return fmt.Errorf("migration %s failed: %w", migration.name, err)
		}

		// Record migration
		if err := recordMigration(db, migration.name); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration.name, err)
		}
	}

	return nil
}

type migration struct {
	name string
	sql  string
}

func createMigrationsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		migration_name TEXT PRIMARY KEY,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(query)
	return err
}

func readMigrationFiles() ([]migration, error) {
	var migrations []migration

	// Get the migrations directory path
	// This assumes the binary is run from the backend directory
	migrationDir := "pkg/db/migration/sqlite"
	
	// Try multiple possible paths
	possiblePaths := []string{
		migrationDir,
		"./pkg/db/migration/sqlite",
		"../pkg/db/migration/sqlite",
	}
	
	var entries []fs.DirEntry
	var err error
	
	for _, path := range possiblePaths {
		entries, err = os.ReadDir(path)
		if err == nil {
			migrationDir = path
			break
		}
	}
	
	if err != nil {
		return nil, fmt.Errorf("migrations directory not found: %w", err)
	}

	// Map to store migrations by name
	migMap := make(map[string]*migration)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Only process .up.sql files
		if len(name) < 8 || !strings.HasSuffix(name, ".up.sql") {
			continue
		}

		// Extract base name (e.g., "000001_create_users_table")
		baseName := name[:len(name)-7] // Remove ".up.sql"

		content, err := os.ReadFile(filepath.Join(migrationDir, name))
		if err != nil {
			return nil, err
		}

		if _, exists := migMap[baseName]; !exists {
			migMap[baseName] = &migration{name: baseName}
		}

		migMap[baseName].sql = string(content)
	}

	// Convert map to sorted slice
	for _, mig := range migMap {
		if mig.sql != "" { // Only add if we have SQL
			migrations = append(migrations, *mig)
		}
	}

	// Sort by name to ensure order
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].name < migrations[j].name
	})

	return migrations, nil
}

func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	applied := make(map[string]bool)

	rows, err := db.Query("SELECT migration_name FROM schema_migrations")
	if err != nil {
		return applied, nil // Table might not exist yet
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		applied[name] = true
	}

	return applied, nil
}

func runMigration(db *sql.DB, mig migration) error {
	_, err := db.Exec(mig.sql)
	return err
}

func recordMigration(db *sql.DB, name string) error {
	_, err := db.Exec("INSERT INTO schema_migrations (migration_name) VALUES (?)", name)
	return err
}
