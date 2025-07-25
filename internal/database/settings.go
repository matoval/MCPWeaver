package database

import (
	"database/sql"
	"fmt"
	"time"
)

// SettingsRepository handles CRUD operations for application settings
type SettingsRepository struct {
	db *DB
}

// NewSettingsRepository creates a new settings repository
func NewSettingsRepository(db *DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

// Set sets a setting value
func (r *SettingsRepository) Set(key, value, settingType string) error {
	query := `
		INSERT INTO app_settings (key, value, type, updated_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			type = excluded.type,
			updated_at = excluded.updated_at
	`

	_, err := r.db.conn.Exec(query, key, value, settingType, time.Now())
	if err != nil {
		return fmt.Errorf("failed to set setting: %w", err)
	}

	return nil
}

// Get retrieves a setting value
func (r *SettingsRepository) Get(key string) (*AppSetting, error) {
	query := `
		SELECT key, value, type, updated_at
		FROM app_settings 
		WHERE key = ?
	`

	setting := &AppSetting{}
	err := r.db.conn.QueryRow(query, key).Scan(
		&setting.Key, &setting.Value, &setting.Type, &setting.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("setting not found")
		}
		return nil, fmt.Errorf("failed to get setting: %w", err)
	}

	return setting, nil
}

// GetAll retrieves all settings
func (r *SettingsRepository) GetAll() ([]*AppSetting, error) {
	query := `
		SELECT key, value, type, updated_at
		FROM app_settings 
		ORDER BY key
	`

	rows, err := r.db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query settings: %w", err)
	}
	defer rows.Close()

	var settings []*AppSetting
	for rows.Next() {
		setting := &AppSetting{}
		err := rows.Scan(&setting.Key, &setting.Value, &setting.Type, &setting.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan setting: %w", err)
		}
		settings = append(settings, setting)
	}

	return settings, nil
}

// Delete deletes a setting
func (r *SettingsRepository) Delete(key string) error {
	query := `DELETE FROM app_settings WHERE key = ?`

	result, err := r.db.conn.Exec(query, key)
	if err != nil {
		return fmt.Errorf("failed to delete setting: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("setting not found")
	}

	return nil
}

// GetByType retrieves settings by type
func (r *SettingsRepository) GetByType(settingType string) ([]*AppSetting, error) {
	query := `
		SELECT key, value, type, updated_at
		FROM app_settings 
		WHERE type = ?
		ORDER BY key
	`

	rows, err := r.db.conn.Query(query, settingType)
	if err != nil {
		return nil, fmt.Errorf("failed to query settings by type: %w", err)
	}
	defer rows.Close()

	var settings []*AppSetting
	for rows.Next() {
		setting := &AppSetting{}
		err := rows.Scan(&setting.Key, &setting.Value, &setting.Type, &setting.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan setting: %w", err)
		}
		settings = append(settings, setting)
	}

	return settings, nil
}

// SetMultiple sets multiple settings in a transaction
func (r *SettingsRepository) SetMultiple(settings map[string]AppSetting) error {
	tx, err := r.db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO app_settings (key, value, type, updated_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			type = excluded.type,
			updated_at = excluded.updated_at
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now()
	for key, setting := range settings {
		_, err := stmt.Exec(key, setting.Value, setting.Type, now)
		if err != nil {
			return fmt.Errorf("failed to execute statement for key %s: %w", key, err)
		}
	}

	return tx.Commit()
}
