package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// TemplateVariable represents a template variable (database layer)
type TemplateVariable struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Type         string   `json:"type"`
	DefaultValue string   `json:"defaultValue"`
	Required     bool     `json:"required"`
	Options      []string `json:"options,omitempty"`
	Validation   string   `json:"validation,omitempty"`
}

// AppTemplate represents a template (database layer)
type AppTemplate struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Version     string             `json:"version"`
	Author      string             `json:"author"`
	Type        string             `json:"type"`
	Path        string             `json:"path"`
	IsBuiltIn   bool               `json:"isBuiltIn"`
	Variables   []TemplateVariable `json:"variables"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

// createTemplateError creates a database error for template operations
func createTemplateError(operation, message string, err error) error {
	dbErr := NewDatabaseError(operation, "templates", fmt.Errorf(message))
	if err != nil {
		dbErr.Err = err
	}
	return dbErr
}

// TemplateRepository handles template database operations
type TemplateRepository struct {
	db *DB
}

// NewTemplateRepository creates a new template repository
func NewTemplateRepository(db *DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

// Create inserts a new template into the database
func (r *TemplateRepository) Create(template *AppTemplate) error {
	if template == nil {
		return createTemplateError("CREATE_TEMPLATE", "template cannot be nil", nil)
	}

	// Serialize variables to JSON
	variablesJSON, err := json.Marshal(template.Variables)
	if err != nil {
		return createTemplateError("SERIALIZE_VARIABLES", "failed to serialize template variables", err)
	}

	// Set timestamps
	now := time.Now()
	template.CreatedAt = now
	template.UpdatedAt = now

	query := `
		INSERT INTO templates (id, name, description, version, author, type, path, is_built_in, variables, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.conn.Exec(query,
		template.ID,
		template.Name,
		template.Description,
		template.Version,
		template.Author,
		template.Type,
		template.Path,
		template.IsBuiltIn,
		string(variablesJSON),
		template.CreatedAt,
		template.UpdatedAt,
	)

	if err != nil {
		return createTemplateError("CREATE_TEMPLATE", "failed to create template", err)
	}

	return nil
}

// GetByID retrieves a template by its ID
func (r *TemplateRepository) GetByID(id string) (*AppTemplate, error) {
	if id == "" {
		return nil, createTemplateError("template ID cannot be empty", "GET_TEMPLATE", nil)
	}

	query := `
		SELECT id, name, description, version, author, type, path, is_built_in, variables, created_at, updated_at
		FROM templates
		WHERE id = ?
	`

	var dbTemplate Template
	err := r.db.conn.QueryRow(query, id).Scan(
		&dbTemplate.ID,
		&dbTemplate.Name,
		&dbTemplate.Description,
		&dbTemplate.Version,
		&dbTemplate.Author,
		&dbTemplate.Type,
		&dbTemplate.Path,
		&dbTemplate.IsBuiltIn,
		&dbTemplate.Variables,
		&dbTemplate.CreatedAt,
		&dbTemplate.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, createTemplateError("template not found", "TEMPLATE_NOT_FOUND", err)
		}
		return nil, createTemplateError("failed to get template", "GET_TEMPLATE", err)
	}

	return r.convertToAppTemplate(&dbTemplate)
}

// GetByName retrieves a template by its name
func (r *TemplateRepository) GetByName(name string) (*AppTemplate, error) {
	if name == "" {
		return nil, createTemplateError("template name cannot be empty", "GET_TEMPLATE", nil)
	}

	query := `
		SELECT id, name, description, version, author, type, path, is_built_in, variables, created_at, updated_at
		FROM templates
		WHERE name = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	var dbTemplate Template
	err := r.db.conn.QueryRow(query, name).Scan(
		&dbTemplate.ID,
		&dbTemplate.Name,
		&dbTemplate.Description,
		&dbTemplate.Version,
		&dbTemplate.Author,
		&dbTemplate.Type,
		&dbTemplate.Path,
		&dbTemplate.IsBuiltIn,
		&dbTemplate.Variables,
		&dbTemplate.CreatedAt,
		&dbTemplate.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, createTemplateError("template not found", "TEMPLATE_NOT_FOUND", err)
		}
		return nil, createTemplateError("failed to get template", "GET_TEMPLATE", err)
	}

	return r.convertToAppTemplate(&dbTemplate)
}

// GetAll retrieves all templates
func (r *TemplateRepository) GetAll() ([]*AppTemplate, error) {
	query := `
		SELECT id, name, description, version, author, type, path, is_built_in, variables, created_at, updated_at
		FROM templates
		ORDER BY is_built_in DESC, name ASC
	`

	rows, err := r.db.conn.Query(query)
	if err != nil {
		return nil, createTemplateError("failed to get templates", "GET_TEMPLATES", err)
	}
	defer rows.Close()

	var templates []*AppTemplate
	for rows.Next() {
		var dbTemplate Template
		err := rows.Scan(
			&dbTemplate.ID,
			&dbTemplate.Name,
			&dbTemplate.Description,
			&dbTemplate.Version,
			&dbTemplate.Author,
			&dbTemplate.Type,
			&dbTemplate.Path,
			&dbTemplate.IsBuiltIn,
			&dbTemplate.Variables,
			&dbTemplate.CreatedAt,
			&dbTemplate.UpdatedAt,
		)

		if err != nil {
			return nil, createTemplateError("failed to scan template", "SCAN_TEMPLATE", err)
		}

		template, err := r.convertToAppTemplate(&dbTemplate)
		if err != nil {
			return nil, err
		}

		templates = append(templates, template)
	}

	if err = rows.Err(); err != nil {
		return nil, createTemplateError("error iterating templates", "ITERATE_TEMPLATES", err)
	}

	return templates, nil
}

// GetByType retrieves templates by their type
func (r *TemplateRepository) GetByType(templateType string) ([]*AppTemplate, error) {
	query := `
		SELECT id, name, description, version, author, type, path, is_built_in, variables, created_at, updated_at
		FROM templates
		WHERE type = ?
		ORDER BY name ASC
	`

	rows, err := r.db.conn.Query(query, templateType)
	if err != nil {
		return nil, createTemplateError("failed to get templates by type", "GET_TEMPLATES_BY_TYPE", err)
	}
	defer rows.Close()

	var templates []*AppTemplate
	for rows.Next() {
		var dbTemplate Template
		err := rows.Scan(
			&dbTemplate.ID,
			&dbTemplate.Name,
			&dbTemplate.Description,
			&dbTemplate.Version,
			&dbTemplate.Author,
			&dbTemplate.Type,
			&dbTemplate.Path,
			&dbTemplate.IsBuiltIn,
			&dbTemplate.Variables,
			&dbTemplate.CreatedAt,
			&dbTemplate.UpdatedAt,
		)

		if err != nil {
			return nil, createTemplateError("failed to scan template", "SCAN_TEMPLATE", err)
		}

		template, err := r.convertToAppTemplate(&dbTemplate)
		if err != nil {
			return nil, err
		}

		templates = append(templates, template)
	}

	if err = rows.Err(); err != nil {
		return nil, createTemplateError("error iterating templates", "ITERATE_TEMPLATES", err)
	}

	return templates, nil
}

// Update updates an existing template
func (r *TemplateRepository) Update(template *AppTemplate) error {
	if template == nil {
		return createTemplateError("template cannot be nil", "UPDATE_TEMPLATE", nil)
	}

	if template.ID == "" {
		return createTemplateError("template ID cannot be empty", "UPDATE_TEMPLATE", nil)
	}

	// Serialize variables to JSON
	variablesJSON, err := json.Marshal(template.Variables)
	if err != nil {
		return createTemplateError("SERIALIZE_VARIABLES", "failed to serialize template variables", err)
	}

	// Update timestamp
	template.UpdatedAt = time.Now()

	query := `
		UPDATE templates
		SET name = ?, description = ?, version = ?, author = ?, type = ?, path = ?, 
		    is_built_in = ?, variables = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.conn.Exec(query,
		template.Name,
		template.Description,
		template.Version,
		template.Author,
		template.Type,
		template.Path,
		template.IsBuiltIn,
		string(variablesJSON),
		template.UpdatedAt,
		template.ID,
	)

	if err != nil {
		return createTemplateError("failed to update template", "UPDATE_TEMPLATE", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return createTemplateError("failed to check update result", "UPDATE_TEMPLATE", err)
	}

	if rowsAffected == 0 {
		return createTemplateError("template not found", "TEMPLATE_NOT_FOUND", nil)
	}

	return nil
}

// Delete removes a template from the database
func (r *TemplateRepository) Delete(id string) error {
	if id == "" {
		return createTemplateError("template ID cannot be empty", "DELETE_TEMPLATE", nil)
	}

	query := `DELETE FROM templates WHERE id = ?`

	result, err := r.db.conn.Exec(query, id)
	if err != nil {
		return createTemplateError("failed to delete template", "DELETE_TEMPLATE", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return createTemplateError("failed to check delete result", "DELETE_TEMPLATE", err)
	}

	if rowsAffected == 0 {
		return createTemplateError("template not found", "TEMPLATE_NOT_FOUND", nil)
	}

	return nil
}

// Search searches templates by name or description
func (r *TemplateRepository) Search(query string) ([]*AppTemplate, error) {
	if query == "" {
		return r.GetAll()
	}

	searchQuery := `
		SELECT id, name, description, version, author, type, path, is_built_in, variables, created_at, updated_at
		FROM templates
		WHERE name LIKE ? OR description LIKE ?
		ORDER BY is_built_in DESC, name ASC
	`

	searchPattern := "%" + query + "%"
	rows, err := r.db.conn.Query(searchQuery, searchPattern, searchPattern)
	if err != nil {
		return nil, createTemplateError("failed to search templates", "SEARCH_TEMPLATES", err)
	}
	defer rows.Close()

	var templates []*AppTemplate
	for rows.Next() {
		var dbTemplate Template
		err := rows.Scan(
			&dbTemplate.ID,
			&dbTemplate.Name,
			&dbTemplate.Description,
			&dbTemplate.Version,
			&dbTemplate.Author,
			&dbTemplate.Type,
			&dbTemplate.Path,
			&dbTemplate.IsBuiltIn,
			&dbTemplate.Variables,
			&dbTemplate.CreatedAt,
			&dbTemplate.UpdatedAt,
		)

		if err != nil {
			return nil, createTemplateError("failed to scan template", "SCAN_TEMPLATE", err)
		}

		template, err := r.convertToAppTemplate(&dbTemplate)
		if err != nil {
			return nil, err
		}

		templates = append(templates, template)
	}

	if err = rows.Err(); err != nil {
		return nil, createTemplateError("error iterating templates", "ITERATE_TEMPLATES", err)
	}

	return templates, nil
}

// GetBuiltInTemplates retrieves all built-in templates
func (r *TemplateRepository) GetBuiltInTemplates() ([]*AppTemplate, error) {
	query := `
		SELECT id, name, description, version, author, type, path, is_built_in, variables, created_at, updated_at
		FROM templates
		WHERE is_built_in = true
		ORDER BY name ASC
	`

	rows, err := r.db.conn.Query(query)
	if err != nil {
		return nil, createTemplateError("failed to get built-in templates", "GET_BUILTIN_TEMPLATES", err)
	}
	defer rows.Close()

	var templates []*AppTemplate
	for rows.Next() {
		var dbTemplate Template
		err := rows.Scan(
			&dbTemplate.ID,
			&dbTemplate.Name,
			&dbTemplate.Description,
			&dbTemplate.Version,
			&dbTemplate.Author,
			&dbTemplate.Type,
			&dbTemplate.Path,
			&dbTemplate.IsBuiltIn,
			&dbTemplate.Variables,
			&dbTemplate.CreatedAt,
			&dbTemplate.UpdatedAt,
		)

		if err != nil {
			return nil, createTemplateError("failed to scan template", "SCAN_TEMPLATE", err)
		}

		template, err := r.convertToAppTemplate(&dbTemplate)
		if err != nil {
			return nil, err
		}

		templates = append(templates, template)
	}

	if err = rows.Err(); err != nil {
		return nil, createTemplateError("error iterating templates", "ITERATE_TEMPLATES", err)
	}

	return templates, nil
}

// convertToAppTemplate converts a database template to an app template
func (r *TemplateRepository) convertToAppTemplate(dbTemplate *Template) (*AppTemplate, error) {
	// Parse variables JSON
	var variables []TemplateVariable
	if dbTemplate.Variables != "" {
		if err := json.Unmarshal([]byte(dbTemplate.Variables), &variables); err != nil {
			return nil, createTemplateError("failed to parse template variables", "PARSE_VARIABLES", err)
		}
	}

	return &AppTemplate{
		ID:          dbTemplate.ID,
		Name:        dbTemplate.Name,
		Description: dbTemplate.Description,
		Version:     dbTemplate.Version,
		Author:      dbTemplate.Author,
		Type:        dbTemplate.Type,
		Path:        dbTemplate.Path,
		IsBuiltIn:   dbTemplate.IsBuiltIn,
		Variables:   variables,
		CreatedAt:   dbTemplate.CreatedAt,
		UpdatedAt:   dbTemplate.UpdatedAt,
	}, nil
}

// Exists checks if a template exists by ID
func (r *TemplateRepository) Exists(id string) (bool, error) {
	if id == "" {
		return false, createTemplateError("template ID cannot be empty", "CHECK_TEMPLATE_EXISTS", nil)
	}

	query := `SELECT COUNT(*) FROM templates WHERE id = ?`
	var count int
	err := r.db.conn.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, createTemplateError("failed to check template existence", "CHECK_TEMPLATE_EXISTS", err)
	}

	return count > 0, nil
}

// Count returns the total number of templates
func (r *TemplateRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM templates`
	var count int
	err := r.db.conn.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, createTemplateError("failed to count templates", "COUNT_TEMPLATES", err)
	}

	return count, nil
}

// CountByType returns the number of templates by type
func (r *TemplateRepository) CountByType(templateType string) (int, error) {
	query := `SELECT COUNT(*) FROM templates WHERE type = ?`
	var count int
	err := r.db.conn.QueryRow(query, templateType).Scan(&count)
	if err != nil {
		return 0, createTemplateError("failed to count templates by type", "COUNT_TEMPLATES_BY_TYPE", err)
	}

	return count, nil
}
