package database

import (
	"database/sql"
	"fmt"
	"time"
)

// GenerationRepository handles CRUD operations for generations
type GenerationRepository struct {
	db *DB
}

// NewGenerationRepository creates a new generation repository
func NewGenerationRepository(db *DB) *GenerationRepository {
	return &GenerationRepository{db: db}
}

// Create creates a new generation
func (r *GenerationRepository) Create(generation *Generation) error {
	query := `
		INSERT INTO generations (id, project_id, status, progress, current_step, start_time, results, errors, processing_time_ms)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	if generation.StartTime.IsZero() {
		generation.StartTime = time.Now()
	}

	_, err := r.db.conn.Exec(query,
		generation.ID, generation.ProjectID, generation.Status, generation.Progress,
		generation.CurrentStep, generation.StartTime, generation.Results,
		generation.Errors, generation.ProcessingTimeMs)

	if err != nil {
		return fmt.Errorf("failed to create generation: %w", err)
	}

	return nil
}

// GetByID retrieves a generation by ID
func (r *GenerationRepository) GetByID(id string) (*Generation, error) {
	query := `
		SELECT id, project_id, status, progress, current_step, start_time, 
			   end_time, results, errors, processing_time_ms
		FROM generations 
		WHERE id = ?
	`

	generation := &Generation{}
	err := r.db.conn.QueryRow(query, id).Scan(
		&generation.ID, &generation.ProjectID, &generation.Status, &generation.Progress,
		&generation.CurrentStep, &generation.StartTime, &generation.EndTime,
		&generation.Results, &generation.Errors, &generation.ProcessingTimeMs)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("generation not found")
		}
		return nil, fmt.Errorf("failed to get generation: %w", err)
	}

	return generation, nil
}

// GetByProjectID retrieves all generations for a project
func (r *GenerationRepository) GetByProjectID(projectID string) ([]*Generation, error) {
	query := `
		SELECT id, project_id, status, progress, current_step, start_time, 
			   end_time, results, errors, processing_time_ms
		FROM generations 
		WHERE project_id = ?
		ORDER BY start_time DESC
	`

	rows, err := r.db.conn.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query generations: %w", err)
	}
	defer rows.Close()

	var generations []*Generation
	for rows.Next() {
		generation := &Generation{}
		err := rows.Scan(
			&generation.ID, &generation.ProjectID, &generation.Status, &generation.Progress,
			&generation.CurrentStep, &generation.StartTime, &generation.EndTime,
			&generation.Results, &generation.Errors, &generation.ProcessingTimeMs)
		if err != nil {
			return nil, fmt.Errorf("failed to scan generation: %w", err)
		}
		generations = append(generations, generation)
	}

	return generations, nil
}

// Update updates a generation
func (r *GenerationRepository) Update(generation *Generation) error {
	query := `
		UPDATE generations 
		SET status = ?, progress = ?, current_step = ?, end_time = ?, 
			results = ?, errors = ?, processing_time_ms = ?
		WHERE id = ?
	`

	result, err := r.db.conn.Exec(query,
		generation.Status, generation.Progress, generation.CurrentStep,
		generation.EndTime, generation.Results, generation.Errors,
		generation.ProcessingTimeMs, generation.ID)

	if err != nil {
		return fmt.Errorf("failed to update generation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("generation not found")
	}

	return nil
}

// Delete deletes a generation
func (r *GenerationRepository) Delete(id string) error {
	query := `DELETE FROM generations WHERE id = ?`

	result, err := r.db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete generation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("generation not found")
	}

	return nil
}

// GetRecent retrieves recent generations across all projects
func (r *GenerationRepository) GetRecent(limit int) ([]*Generation, error) {
	query := `
		SELECT id, project_id, status, progress, current_step, start_time, 
			   end_time, results, errors, processing_time_ms
		FROM generations 
		ORDER BY start_time DESC
		LIMIT ?
	`

	rows, err := r.db.conn.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent generations: %w", err)
	}
	defer rows.Close()

	var generations []*Generation
	for rows.Next() {
		generation := &Generation{}
		err := rows.Scan(
			&generation.ID, &generation.ProjectID, &generation.Status, &generation.Progress,
			&generation.CurrentStep, &generation.StartTime, &generation.EndTime,
			&generation.Results, &generation.Errors, &generation.ProcessingTimeMs)
		if err != nil {
			return nil, fmt.Errorf("failed to scan generation: %w", err)
		}
		generations = append(generations, generation)
	}

	return generations, nil
}

// GetByStatus retrieves generations by status
func (r *GenerationRepository) GetByStatus(status string) ([]*Generation, error) {
	query := `
		SELECT id, project_id, status, progress, current_step, start_time, 
			   end_time, results, errors, processing_time_ms
		FROM generations 
		WHERE status = ?
		ORDER BY start_time DESC
	`

	rows, err := r.db.conn.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query generations by status: %w", err)
	}
	defer rows.Close()

	var generations []*Generation
	for rows.Next() {
		generation := &Generation{}
		err := rows.Scan(
			&generation.ID, &generation.ProjectID, &generation.Status, &generation.Progress,
			&generation.CurrentStep, &generation.StartTime, &generation.EndTime,
			&generation.Results, &generation.Errors, &generation.ProcessingTimeMs)
		if err != nil {
			return nil, fmt.Errorf("failed to scan generation: %w", err)
		}
		generations = append(generations, generation)
	}

	return generations, nil
}
