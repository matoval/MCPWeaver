package database

import (
	"database/sql"
	"fmt"
	"time"
)

// ProjectRepository handles CRUD operations for projects
type ProjectRepository struct {
	db *DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create creates a new project
func (r *ProjectRepository) Create(project *Project) error {
	query := `
		INSERT INTO projects (id, name, spec_path, spec_url, output_path, settings, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	project.CreatedAt = now
	project.UpdatedAt = now
	
	_, err := r.db.conn.Exec(query, 
		project.ID, project.Name, project.SpecPath, project.SpecURL, 
		project.OutputPath, project.Settings, project.Status, 
		project.CreatedAt, project.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	
	return nil
}

// GetByID retrieves a project by ID
func (r *ProjectRepository) GetByID(id string) (*Project, error) {
	query := `
		SELECT id, name, spec_path, spec_url, output_path, settings, status, 
			   created_at, updated_at, last_generated, generation_count
		FROM projects 
		WHERE id = ?
	`
	
	project := &Project{}
	err := r.db.conn.QueryRow(query, id).Scan(
		&project.ID, &project.Name, &project.SpecPath, &project.SpecURL,
		&project.OutputPath, &project.Settings, &project.Status,
		&project.CreatedAt, &project.UpdatedAt, &project.LastGenerated,
		&project.GenerationCount)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found")
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	
	return project, nil
}

// GetAll retrieves all projects
func (r *ProjectRepository) GetAll() ([]*Project, error) {
	query := `
		SELECT id, name, spec_path, spec_url, output_path, settings, status, 
			   created_at, updated_at, last_generated, generation_count
		FROM projects 
		ORDER BY updated_at DESC
	`
	
	rows, err := r.db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()
	
	var projects []*Project
	for rows.Next() {
		project := &Project{}
		err := rows.Scan(
			&project.ID, &project.Name, &project.SpecPath, &project.SpecURL,
			&project.OutputPath, &project.Settings, &project.Status,
			&project.CreatedAt, &project.UpdatedAt, &project.LastGenerated,
			&project.GenerationCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, project)
	}
	
	return projects, nil
}

// Update updates a project
func (r *ProjectRepository) Update(project *Project) error {
	query := `
		UPDATE projects 
		SET name = ?, spec_path = ?, spec_url = ?, output_path = ?, 
			settings = ?, status = ?, updated_at = ?, last_generated = ?, 
			generation_count = ?
		WHERE id = ?
	`
	
	project.UpdatedAt = time.Now()
	
	result, err := r.db.conn.Exec(query, 
		project.Name, project.SpecPath, project.SpecURL, project.OutputPath,
		project.Settings, project.Status, project.UpdatedAt, project.LastGenerated,
		project.GenerationCount, project.ID)
	
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("project not found")
	}
	
	return nil
}

// Delete deletes a project
func (r *ProjectRepository) Delete(id string) error {
	query := `DELETE FROM projects WHERE id = ?`
	
	result, err := r.db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("project not found")
	}
	
	return nil
}

// GetByName retrieves a single project by exact name match
func (r *ProjectRepository) GetByName(name string) (*Project, error) {
	query := `
		SELECT id, name, spec_path, spec_url, output_path, settings, status, 
			   created_at, updated_at, last_generated, generation_count
		FROM projects 
		WHERE name = ?
	`
	
	project := &Project{}
	err := r.db.conn.QueryRow(query, name).Scan(
		&project.ID, &project.Name, &project.SpecPath, &project.SpecURL,
		&project.OutputPath, &project.Settings, &project.Status,
		&project.CreatedAt, &project.UpdatedAt, &project.LastGenerated,
		&project.GenerationCount)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found")
		}
		return nil, fmt.Errorf("failed to get project by name: %w", err)
	}
	
	return project, nil
}

// SearchByName retrieves projects by name (useful for search)
func (r *ProjectRepository) SearchByName(name string) ([]*Project, error) {
	query := `
		SELECT id, name, spec_path, spec_url, output_path, settings, status, 
			   created_at, updated_at, last_generated, generation_count
		FROM projects 
		WHERE name LIKE ? 
		ORDER BY updated_at DESC
	`
	
	rows, err := r.db.conn.Query(query, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to query projects by name: %w", err)
	}
	defer rows.Close()
	
	var projects []*Project
	for rows.Next() {
		project := &Project{}
		err := rows.Scan(
			&project.ID, &project.Name, &project.SpecPath, &project.SpecURL,
			&project.OutputPath, &project.Settings, &project.Status,
			&project.CreatedAt, &project.UpdatedAt, &project.LastGenerated,
			&project.GenerationCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, project)
	}
	
	return projects, nil
}

// GetRecent retrieves recently updated projects
func (r *ProjectRepository) GetRecent(limit int) ([]*Project, error) {
	query := `
		SELECT id, name, spec_path, spec_url, output_path, settings, status, 
			   created_at, updated_at, last_generated, generation_count
		FROM projects 
		ORDER BY updated_at DESC
		LIMIT ?
	`
	
	rows, err := r.db.conn.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent projects: %w", err)
	}
	defer rows.Close()
	
	var projects []*Project
	for rows.Next() {
		project := &Project{}
		err := rows.Scan(
			&project.ID, &project.Name, &project.SpecPath, &project.SpecURL,
			&project.OutputPath, &project.Settings, &project.Status,
			&project.CreatedAt, &project.UpdatedAt, &project.LastGenerated,
			&project.GenerationCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, project)
	}
	
	return projects, nil
}