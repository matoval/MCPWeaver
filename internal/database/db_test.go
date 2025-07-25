package database

import (
	"testing"
	"time"
)

func TestDatabaseConnection(t *testing.T) {
	// Create a temporary database file
	tempFile := ":memory:" // Use in-memory database for testing

	// Test opening the database
	db, err := Open(tempFile)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.conn.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

func TestMigrations(t *testing.T) {
	// Create a temporary database file
	tempFile := ":memory:"

	// Test migrations
	db, err := Open(tempFile)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Verify that tables exist
	tables := []string{"projects", "generations", "templates", "app_settings", "validation_cache", "schema_version"}

	for _, table := range tables {
		var count int
		query := `SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`
		err := db.conn.QueryRow(query, table).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to check table %s: %v", table, err)
		}
		if count != 1 {
			t.Errorf("Table %s does not exist", table)
		}
	}
}

func TestProjectRepository(t *testing.T) {
	// Create a temporary database
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := NewProjectRepository(db)

	// Test creating a project
	project := &Project{
		ID:         "test-project-1",
		Name:       "Test Project",
		SpecPath:   "/path/to/spec.yaml",
		OutputPath: "/path/to/output",
		Settings:   `{"packageName": "test"}`,
		Status:     "created",
	}

	err = repo.Create(project)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Test retrieving the project
	retrieved, err := repo.GetByID("test-project-1")
	if err != nil {
		t.Fatalf("Failed to get project: %v", err)
	}

	if retrieved.Name != "Test Project" {
		t.Errorf("Expected name 'Test Project', got '%s'", retrieved.Name)
	}

	// Test updating the project
	retrieved.Name = "Updated Test Project"
	err = repo.Update(retrieved)
	if err != nil {
		t.Fatalf("Failed to update project: %v", err)
	}

	// Verify the update
	updated, err := repo.GetByID("test-project-1")
	if err != nil {
		t.Fatalf("Failed to get updated project: %v", err)
	}

	if updated.Name != "Updated Test Project" {
		t.Errorf("Expected name 'Updated Test Project', got '%s'", updated.Name)
	}

	// Test getting all projects
	projects, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all projects: %v", err)
	}

	if len(projects) != 1 {
		t.Errorf("Expected 1 project, got %d", len(projects))
	}

	// Test deleting the project
	err = repo.Delete("test-project-1")
	if err != nil {
		t.Fatalf("Failed to delete project: %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID("test-project-1")
	if err == nil {
		t.Error("Expected error when getting deleted project")
	}
}

func TestGenerationRepository(t *testing.T) {
	// Create a temporary database
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// First create a project
	projectRepo := NewProjectRepository(db)
	project := &Project{
		ID:         "test-project-1",
		Name:       "Test Project",
		SpecPath:   "/path/to/spec.yaml",
		OutputPath: "/path/to/output",
		Settings:   `{"packageName": "test"}`,
		Status:     "created",
	}

	err = projectRepo.Create(project)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Test generation repository
	repo := NewGenerationRepository(db)

	generation := &Generation{
		ID:          "test-generation-1",
		ProjectID:   "test-project-1",
		Status:      "started",
		Progress:    0.5,
		CurrentStep: "Processing",
		StartTime:   time.Now(),
		Results:     `{"success": true}`,
		Errors:      `[]`,
	}

	err = repo.Create(generation)
	if err != nil {
		t.Fatalf("Failed to create generation: %v", err)
	}

	// Test retrieving the generation
	retrieved, err := repo.GetByID("test-generation-1")
	if err != nil {
		t.Fatalf("Failed to get generation: %v", err)
	}

	if retrieved.Status != "started" {
		t.Errorf("Expected status 'started', got '%s'", retrieved.Status)
	}

	// Test getting generations by project
	generations, err := repo.GetByProjectID("test-project-1")
	if err != nil {
		t.Fatalf("Failed to get generations by project: %v", err)
	}

	if len(generations) != 1 {
		t.Errorf("Expected 1 generation, got %d", len(generations))
	}
}

func TestSettingsRepository(t *testing.T) {
	// Create a temporary database
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := NewSettingsRepository(db)

	// Test setting a value
	err = repo.Set("test_key", "test_value", "string")
	if err != nil {
		t.Fatalf("Failed to set setting: %v", err)
	}

	// Test getting the value
	setting, err := repo.Get("test_key")
	if err != nil {
		t.Fatalf("Failed to get setting: %v", err)
	}

	if setting.Value != "test_value" {
		t.Errorf("Expected value 'test_value', got '%s'", setting.Value)
	}

	// Test updating the value
	err = repo.Set("test_key", "updated_value", "string")
	if err != nil {
		t.Fatalf("Failed to update setting: %v", err)
	}

	// Verify update
	updated, err := repo.Get("test_key")
	if err != nil {
		t.Fatalf("Failed to get updated setting: %v", err)
	}

	if updated.Value != "updated_value" {
		t.Errorf("Expected value 'updated_value', got '%s'", updated.Value)
	}

	// Test getting all settings
	settings, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all settings: %v", err)
	}

	if len(settings) != 1 {
		t.Errorf("Expected 1 setting, got %d", len(settings))
	}
}
