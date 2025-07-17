package database

import (
	"database/sql"
	"fmt"
	"log"
)

// Migration represents a database migration
type Migration struct {
	Version int
	Name    string
	Up      string
	Down    string
}

// migrations contains all database migrations in order
var migrations = []Migration{
	{
		Version: 1,
		Name:    "initial_schema",
		Up: `
			-- Projects table
			CREATE TABLE projects (
				id TEXT PRIMARY KEY,
				name TEXT NOT NULL,
				spec_path TEXT,
				spec_url TEXT,
				output_path TEXT NOT NULL,
				settings TEXT NOT NULL, -- JSON serialized ProjectSettings
				status TEXT NOT NULL DEFAULT 'created',
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				last_generated DATETIME,
				generation_count INTEGER DEFAULT 0
			);

			-- Indexes for projects table
			CREATE INDEX idx_projects_name ON projects(name);
			CREATE INDEX idx_projects_status ON projects(status);
			CREATE INDEX idx_projects_created_at ON projects(created_at);
			CREATE INDEX idx_projects_last_generated ON projects(last_generated);

			-- Generations table
			CREATE TABLE generations (
				id TEXT PRIMARY KEY,
				project_id TEXT NOT NULL,
				status TEXT NOT NULL,
				progress REAL DEFAULT 0.0,
				current_step TEXT,
				start_time DATETIME DEFAULT CURRENT_TIMESTAMP,
				end_time DATETIME,
				results TEXT, -- JSON serialized GenerationResults
				errors TEXT, -- JSON serialized GenerationError array
				processing_time_ms INTEGER,
				FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
			);

			-- Indexes for generations table
			CREATE INDEX idx_generations_project_id ON generations(project_id);
			CREATE INDEX idx_generations_status ON generations(status);
			CREATE INDEX idx_generations_start_time ON generations(start_time);

			-- Templates table
			CREATE TABLE templates (
				id TEXT PRIMARY KEY,
				name TEXT NOT NULL,
				description TEXT,
				version TEXT NOT NULL,
				author TEXT,
				type TEXT NOT NULL, -- 'default', 'custom', 'plugin'
				path TEXT NOT NULL,
				is_built_in BOOLEAN DEFAULT FALSE,
				variables TEXT, -- JSON serialized TemplateVariable array
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
			);

			-- Indexes for templates table
			CREATE INDEX idx_templates_name ON templates(name);
			CREATE INDEX idx_templates_type ON templates(type);
			CREATE INDEX idx_templates_is_built_in ON templates(is_built_in);

			-- Application settings table
			CREATE TABLE app_settings (
				key TEXT PRIMARY KEY,
				value TEXT NOT NULL,
				type TEXT NOT NULL, -- 'string', 'number', 'boolean', 'json'
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
			);

			-- Validation cache table
			CREATE TABLE validation_cache (
				spec_hash TEXT PRIMARY KEY,
				spec_path TEXT,
				spec_url TEXT,
				validation_result TEXT NOT NULL, -- JSON serialized ValidationResult
				cached_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				expires_at DATETIME NOT NULL
			);

			-- Index for validation cache cleanup
			CREATE INDEX idx_validation_cache_expires_at ON validation_cache(expires_at);

			-- Insert initial schema version
			INSERT INTO schema_version (version) VALUES (1);
		`,
		Down: `
			DROP TABLE IF EXISTS validation_cache;
			DROP TABLE IF EXISTS app_settings;
			DROP TABLE IF EXISTS templates;
			DROP TABLE IF EXISTS generations;
			DROP TABLE IF EXISTS projects;
			DROP TABLE IF EXISTS schema_version;
		`,
	},
}

// getCurrentVersion returns the current database schema version
func (db *DB) getCurrentVersion() (int, error) {
	var version int
	err := db.conn.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_version").Scan(&version)
	if err != nil {
		// If table doesn't exist, we're at version 0
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return version, nil
}

// migrate runs all pending migrations
func (db *DB) migrate() error {
	// Create schema_version table if it doesn't exist
	_, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS schema_version (
			version INTEGER PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_version table: %w", err)
	}

	currentVersion, err := db.getCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	log.Printf("Current database schema version: %d", currentVersion)

	// Apply pending migrations
	for _, migration := range migrations {
		if migration.Version > currentVersion {
			log.Printf("Applying migration %d: %s", migration.Version, migration.Name)
			
			// Begin transaction
			tx, err := db.conn.Begin()
			if err != nil {
				return fmt.Errorf("failed to begin transaction for migration %d: %w", migration.Version, err)
			}

			// Execute migration
			_, err = tx.Exec(migration.Up)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to execute migration %d: %w", migration.Version, err)
			}

			// Update schema version (only if not already done in the migration)
			if migration.Version > 1 {
				_, err = tx.Exec("INSERT INTO schema_version (version) VALUES (?)", migration.Version)
				if err != nil {
					tx.Rollback()
					return fmt.Errorf("failed to update schema version for migration %d: %w", migration.Version, err)
				}
			}

			// Commit transaction
			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
			}

			log.Printf("Successfully applied migration %d", migration.Version)
		}
	}

	log.Printf("Database migrations completed")
	return nil
}

// rollback rolls back the database to a specific version
func (db *DB) rollback(targetVersion int) error {
	currentVersion, err := db.getCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	if targetVersion >= currentVersion {
		return fmt.Errorf("target version %d is not less than current version %d", targetVersion, currentVersion)
	}

	// Apply rollback migrations in reverse order
	for i := len(migrations) - 1; i >= 0; i-- {
		migration := migrations[i]
		if migration.Version > targetVersion && migration.Version <= currentVersion {
			log.Printf("Rolling back migration %d: %s", migration.Version, migration.Name)
			
			// Begin transaction
			tx, err := db.conn.Begin()
			if err != nil {
				return fmt.Errorf("failed to begin transaction for rollback %d: %w", migration.Version, err)
			}

			// Execute rollback
			_, err = tx.Exec(migration.Down)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to execute rollback %d: %w", migration.Version, err)
			}

			// Remove schema version entry
			_, err = tx.Exec("DELETE FROM schema_version WHERE version = ?", migration.Version)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to remove schema version for rollback %d: %w", migration.Version, err)
			}

			// Commit transaction
			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("failed to commit rollback %d: %w", migration.Version, err)
			}

			log.Printf("Successfully rolled back migration %d", migration.Version)
		}
	}

	log.Printf("Database rollback completed")
	return nil
}