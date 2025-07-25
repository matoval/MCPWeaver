package database

import (
	"database/sql"
)

// Repository provides access to all database repositories
type Repository struct {
	Projects    *ProjectRepository
	Generations *GenerationRepository
	Settings    *SettingsRepository
	Templates   *TemplateRepository
	db          *DB
}

// NewRepository creates a new repository manager
func NewRepository(db *DB) *Repository {
	return &Repository{
		Projects:    NewProjectRepository(db),
		Generations: NewGenerationRepository(db),
		Settings:    NewSettingsRepository(db),
		Templates:   NewTemplateRepository(db),
		db:          db,
	}
}

// Close closes the database connection
func (r *Repository) Close() error {
	return r.db.Close()
}

// GetDB returns the underlying database connection
func (r *Repository) GetDB() *DB {
	return r.db
}

// Transaction executes a function within a database transaction
func (r *Repository) Transaction(fn func(*sql.Tx) error) error {
	tx, err := r.db.conn.Begin()
	if err != nil {
		return err
	}

	// Execute the function
	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
