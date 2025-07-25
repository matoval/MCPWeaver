package database

import (
	"errors"
	"fmt"
)

// Common database errors
var (
	ErrNotFound      = errors.New("record not found")
	ErrAlreadyExists = errors.New("record already exists")
	ErrInvalidData   = errors.New("invalid data")
	ErrTransaction   = errors.New("transaction failed")
	ErrConnection    = errors.New("database connection failed")
	ErrMigration     = errors.New("migration failed")
)

// DatabaseError represents a database error with additional context
type DatabaseError struct {
	Op     string // Operation that failed
	Err    error  // Underlying error
	Table  string // Table involved
	ID     string // Record ID if applicable
	Detail string // Additional detail
}

// Error implements the error interface
func (e *DatabaseError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("database error in %s on table %s: %v - %s", e.Op, e.Table, e.Err, e.Detail)
	}
	return fmt.Sprintf("database error in %s on table %s: %v", e.Op, e.Table, e.Err)
}

// Unwrap returns the underlying error
func (e *DatabaseError) Unwrap() error {
	return e.Err
}

// NewDatabaseError creates a new database error
func NewDatabaseError(op, table string, err error) *DatabaseError {
	return &DatabaseError{
		Op:    op,
		Table: table,
		Err:   err,
	}
}

// WithID adds an ID to the database error
func (e *DatabaseError) WithID(id string) *DatabaseError {
	e.ID = id
	return e
}

// WithDetail adds additional detail to the database error
func (e *DatabaseError) WithDetail(detail string) *DatabaseError {
	e.Detail = detail
	return e
}

// IsNotFound checks if the error is a "not found" error
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsAlreadyExists checks if the error is an "already exists" error
func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}

// IsInvalidData checks if the error is an "invalid data" error
func IsInvalidData(err error) bool {
	return errors.Is(err, ErrInvalidData)
}

// IsTransaction checks if the error is a transaction error
func IsTransaction(err error) bool {
	return errors.Is(err, ErrTransaction)
}

// IsConnection checks if the error is a connection error
func IsConnection(err error) bool {
	return errors.Is(err, ErrConnection)
}

// IsMigration checks if the error is a migration error
func IsMigration(err error) bool {
	return errors.Is(err, ErrMigration)
}
