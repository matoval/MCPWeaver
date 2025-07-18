package database

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestValidationCacheRepository_StoreAndGet(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := NewValidationCacheRepository(db)

	// Create test cache entry
	cache := &ValidationCache{
		SpecHash:         "test-hash-123",
		SpecPath:         "/test/path/spec.yaml",
		SpecURL:          "",
		ValidationResult: `{"valid": true, "errors": [], "warnings": []}`,
		CachedAt:         time.Now(),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	}

	// Store cache entry
	err = repo.Store(cache)
	if err != nil {
		t.Fatalf("Failed to store cache: %v", err)
	}

	// Retrieve cache entry
	retrieved, err := repo.GetByHash("test-hash-123")
	if err != nil {
		t.Fatalf("Failed to get cache: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Cache entry not found")
	}

	if retrieved.SpecHash != cache.SpecHash {
		t.Errorf("SpecHash = %s, want %s", retrieved.SpecHash, cache.SpecHash)
	}

	if retrieved.SpecPath != cache.SpecPath {
		t.Errorf("SpecPath = %s, want %s", retrieved.SpecPath, cache.SpecPath)
	}

	if retrieved.ValidationResult != cache.ValidationResult {
		t.Errorf("ValidationResult = %s, want %s", retrieved.ValidationResult, cache.ValidationResult)
	}
}

func TestValidationCacheRepository_GetByHash_NotFound(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := NewValidationCacheRepository(db)

	// Try to get non-existent cache entry
	retrieved, err := repo.GetByHash("non-existent-hash")
	if err != nil {
		t.Fatalf("Failed to get cache: %v", err)
	}

	if retrieved != nil {
		t.Error("Should return nil for non-existent cache entry")
	}
}

func TestValidationCacheRepository_ExpiredEntry(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := NewValidationCacheRepository(db)

	// Create expired cache entry
	cache := &ValidationCache{
		SpecHash:         "expired-hash",
		SpecPath:         "/test/path/spec.yaml",
		SpecURL:          "",
		ValidationResult: `{"valid": true, "errors": [], "warnings": []}`,
		CachedAt:         time.Now().Add(-25 * time.Hour),
		ExpiresAt:        time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
	}

	// Store expired cache entry
	err = repo.Store(cache)
	if err != nil {
		t.Fatalf("Failed to store cache: %v", err)
	}

	// Try to retrieve expired cache entry
	retrieved, err := repo.GetByHash("expired-hash")
	if err != nil {
		t.Fatalf("Failed to get cache: %v", err)
	}

	if retrieved != nil {
		t.Error("Should return nil for expired cache entry")
	}
}

func TestValidationCacheRepository_CleanExpired(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := NewValidationCacheRepository(db)

	// Create valid cache entry
	validCache := &ValidationCache{
		SpecHash:         "valid-hash",
		SpecPath:         "/test/path/spec.yaml",
		SpecURL:          "",
		ValidationResult: `{"valid": true, "errors": [], "warnings": []}`,
		CachedAt:         time.Now(),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	}

	// Create expired cache entry
	expiredCache := &ValidationCache{
		SpecHash:         "expired-hash",
		SpecPath:         "/test/path/spec2.yaml",
		SpecURL:          "",
		ValidationResult: `{"valid": true, "errors": [], "warnings": []}`,
		CachedAt:         time.Now().Add(-25 * time.Hour),
		ExpiresAt:        time.Now().Add(-1 * time.Hour),
	}

	// Store both entries
	err = repo.Store(validCache)
	if err != nil {
		t.Fatalf("Failed to store valid cache: %v", err)
	}

	err = repo.Store(expiredCache)
	if err != nil {
		t.Fatalf("Failed to store expired cache: %v", err)
	}

	// Clean expired entries
	err = repo.CleanExpired()
	if err != nil {
		t.Fatalf("Failed to clean expired entries: %v", err)
	}

	// Check that valid entry still exists
	validRetrieved, err := repo.GetByHash("valid-hash")
	if err != nil {
		t.Fatalf("Failed to get valid cache: %v", err)
	}

	if validRetrieved == nil {
		t.Error("Valid cache entry should still exist after cleaning")
	}

	// Check that expired entry is gone
	expiredRetrieved, err := repo.GetByHash("expired-hash")
	if err != nil {
		t.Fatalf("Failed to get expired cache: %v", err)
	}

	if expiredRetrieved != nil {
		t.Error("Expired cache entry should be removed after cleaning")
	}
}

func TestValidationCacheRepository_GetStats(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := NewValidationCacheRepository(db)

	// Initially should have no entries
	stats, err := repo.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats.TotalEntries != 0 {
		t.Errorf("TotalEntries = %d, want 0", stats.TotalEntries)
	}

	if stats.ActiveEntries != 0 {
		t.Errorf("ActiveEntries = %d, want 0", stats.ActiveEntries)
	}

	if stats.ExpiredEntries != 0 {
		t.Errorf("ExpiredEntries = %d, want 0", stats.ExpiredEntries)
	}

	// Add some entries
	validCache := &ValidationCache{
		SpecHash:         "valid-hash",
		SpecPath:         "/test/path/spec.yaml",
		SpecURL:          "",
		ValidationResult: `{"valid": true, "errors": [], "warnings": []}`,
		CachedAt:         time.Now(),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	}

	expiredCache := &ValidationCache{
		SpecHash:         "expired-hash",
		SpecPath:         "/test/path/spec2.yaml",
		SpecURL:          "",
		ValidationResult: `{"valid": true, "errors": [], "warnings": []}`,
		CachedAt:         time.Now().Add(-25 * time.Hour),
		ExpiresAt:        time.Now().Add(-1 * time.Hour),
	}

	err = repo.Store(validCache)
	if err != nil {
		t.Fatalf("Failed to store valid cache: %v", err)
	}

	err = repo.Store(expiredCache)
	if err != nil {
		t.Fatalf("Failed to store expired cache: %v", err)
	}

	// Get stats
	stats, err = repo.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats.TotalEntries != 2 {
		t.Errorf("TotalEntries = %d, want 2", stats.TotalEntries)
	}

	if stats.ActiveEntries != 1 {
		t.Errorf("ActiveEntries = %d, want 1", stats.ActiveEntries)
	}

	if stats.ExpiredEntries != 1 {
		t.Errorf("ExpiredEntries = %d, want 1", stats.ExpiredEntries)
	}
}

func TestValidationCacheRepository_GenerateSpecHash(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := NewValidationCacheRepository(db)

	// Create test file
	testContent := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      responses:
        '200':
          description: Success
`

	testFilePath := filepath.Join(tempDir, "test-spec.yaml")
	err = os.WriteFile(testFilePath, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Generate hash
	hash1, err := repo.GenerateSpecHash(testFilePath)
	if err != nil {
		t.Fatalf("Failed to generate hash: %v", err)
	}

	if hash1 == "" {
		t.Error("Generated hash should not be empty")
	}

	// Generate hash again for same file
	hash2, err := repo.GenerateSpecHash(testFilePath)
	if err != nil {
		t.Fatalf("Failed to generate hash: %v", err)
	}

	if hash1 != hash2 {
		t.Error("Hash should be consistent for same file")
	}

	// Modify file and generate hash again
	modifiedContent := testContent + "\n  /test2:\n    get:\n      responses:\n        '200':\n          description: Success"
	err = os.WriteFile(testFilePath, []byte(modifiedContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write modified test file: %v", err)
	}

	hash3, err := repo.GenerateSpecHash(testFilePath)
	if err != nil {
		t.Fatalf("Failed to generate hash for modified file: %v", err)
	}

	if hash1 == hash3 {
		t.Error("Hash should be different for modified file")
	}
}

func TestValidationCacheRepository_GenerateURLHash(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := NewValidationCacheRepository(db)

	url := "https://example.com/api/spec.yaml"
	time1 := time.Now()
	time2 := time1.Add(1 * time.Hour)

	// Generate hash for same URL at different times
	hash1 := repo.GenerateURLHash(url, time1)
	hash2 := repo.GenerateURLHash(url, time2)

	if hash1 == "" {
		t.Error("Generated hash should not be empty")
	}

	if hash1 == hash2 {
		t.Error("Hash should be different for different times")
	}

	// Generate hash for same URL and time
	hash3 := repo.GenerateURLHash(url, time1)
	if hash1 != hash3 {
		t.Error("Hash should be consistent for same URL and time")
	}
}

func TestValidationCacheRepository_GenerateSpecHash_NonExistentFile(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := NewValidationCacheRepository(db)

	// Try to generate hash for non-existent file
	_, err = repo.GenerateSpecHash("/non/existent/file.yaml")
	if err == nil {
		t.Error("Should return error for non-existent file")
	}
}