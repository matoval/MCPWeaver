package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// DB represents the database connection
type DB struct {
	conn *sql.DB
}

// Open opens a connection to the SQLite database
func Open(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{conn: conn}
	
	// Optimize database connection for performance
	if err := db.optimize(); err != nil {
		return nil, fmt.Errorf("failed to optimize database: %w", err)
	}
	
	// Run migrations
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// migrate is implemented in migrations.go

// optimize configures SQLite for better performance
func (db *DB) optimize() error {
	optimizations := []string{
		"PRAGMA journal_mode=WAL",        // Write-Ahead Logging for better concurrency
		"PRAGMA synchronous=NORMAL",      // Faster writes with reasonable safety
		"PRAGMA cache_size=10000",        // 10MB cache (default is 2MB)
		"PRAGMA temp_store=MEMORY",       // Keep temp tables in memory
		"PRAGMA mmap_size=268435456",     // 256MB memory map size
		"PRAGMA optimize",               // Optimize query planner
	}
	
	for _, pragma := range optimizations {
		if _, err := db.conn.Exec(pragma); err != nil {
			return fmt.Errorf("failed to execute optimization '%s': %w", pragma, err)
		}
	}
	
	return nil
}

// GetConn returns the underlying database connection
func (db *DB) GetConn() *sql.DB {
	return db.conn
}