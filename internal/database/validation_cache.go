package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"
)

// ValidationCacheRepository handles CRUD operations for validation cache
type ValidationCacheRepository struct {
	db *DB
}

// NewValidationCacheRepository creates a new validation cache repository
func NewValidationCacheRepository(db *DB) *ValidationCacheRepository {
	return &ValidationCacheRepository{db: db}
}

// GetByHash retrieves a validation result by spec hash
func (r *ValidationCacheRepository) GetByHash(specHash string) (*ValidationCache, error) {
	query := `
		SELECT spec_hash, spec_path, spec_url, validation_result, cached_at, expires_at
		FROM validation_cache 
		WHERE spec_hash = ? AND expires_at > ?
	`
	
	cache := &ValidationCache{}
	err := r.db.conn.QueryRow(query, specHash, time.Now()).Scan(
		&cache.SpecHash, &cache.SpecPath, &cache.SpecURL, 
		&cache.ValidationResult, &cache.CachedAt, &cache.ExpiresAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get validation cache: %w", err)
	}
	
	return cache, nil
}

// Store stores a validation result in the cache
func (r *ValidationCacheRepository) Store(cache *ValidationCache) error {
	query := `
		INSERT OR REPLACE INTO validation_cache 
		(spec_hash, spec_path, spec_url, validation_result, cached_at, expires_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.conn.Exec(query, 
		cache.SpecHash, cache.SpecPath, cache.SpecURL, 
		cache.ValidationResult, cache.CachedAt, cache.ExpiresAt)
	
	if err != nil {
		return fmt.Errorf("failed to store validation cache: %w", err)
	}
	
	return nil
}

// CleanExpired removes expired validation cache entries
func (r *ValidationCacheRepository) CleanExpired() error {
	query := `DELETE FROM validation_cache WHERE expires_at <= ?`
	
	result, err := r.db.conn.Exec(query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to clean expired validation cache: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	fmt.Printf("Cleaned %d expired validation cache entries\n", rowsAffected)
	return nil
}

// GetStats returns cache statistics
func (r *ValidationCacheRepository) GetStats() (*ValidationCacheStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_entries,
			COUNT(CASE WHEN expires_at > ? THEN 1 END) as active_entries,
			COUNT(CASE WHEN expires_at <= ? THEN 1 END) as expired_entries
		FROM validation_cache
	`
	
	stats := &ValidationCacheStats{}
	now := time.Now()
	
	err := r.db.conn.QueryRow(query, now, now).Scan(
		&stats.TotalEntries, &stats.ActiveEntries, &stats.ExpiredEntries)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get cache stats: %w", err)
	}
	
	return stats, nil
}

// GenerateSpecHash generates a hash for a specification file
func (r *ValidationCacheRepository) GenerateSpecHash(specPath string) (string, error) {
	file, err := os.Open(specPath)
	if err != nil {
		return "", fmt.Errorf("failed to open spec file: %w", err)
	}
	defer file.Close()
	
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash spec file: %w", err)
	}
	
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// GenerateURLHash generates a hash for a specification URL
func (r *ValidationCacheRepository) GenerateURLHash(specURL string, lastModified time.Time) string {
	hash := sha256.New()
	hash.Write([]byte(specURL))
	hash.Write([]byte(lastModified.Format(time.RFC3339)))
	return hex.EncodeToString(hash.Sum(nil))
}

// ValidationCacheStats represents cache statistics
type ValidationCacheStats struct {
	TotalEntries   int `json:"totalEntries"`
	ActiveEntries  int `json:"activeEntries"`
	ExpiredEntries int `json:"expiredEntries"`
}