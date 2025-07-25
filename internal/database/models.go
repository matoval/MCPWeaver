package database

import (
	"time"
)

// Project represents a project in the database
type Project struct {
	ID              string     `json:"id" db:"id"`
	Name            string     `json:"name" db:"name"`
	SpecPath        string     `json:"specPath" db:"spec_path"`
	SpecURL         string     `json:"specUrl" db:"spec_url"`
	OutputPath      string     `json:"outputPath" db:"output_path"`
	Settings        string     `json:"settings" db:"settings"` // JSON serialized ProjectSettings
	Status          string     `json:"status" db:"status"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time  `json:"updatedAt" db:"updated_at"`
	LastGenerated   *time.Time `json:"lastGenerated" db:"last_generated"`
	GenerationCount int        `json:"generationCount" db:"generation_count"`
}

// Generation represents a generation job in the database
type Generation struct {
	ID               string     `json:"id" db:"id"`
	ProjectID        string     `json:"projectId" db:"project_id"`
	Status           string     `json:"status" db:"status"`
	Progress         float64    `json:"progress" db:"progress"`
	CurrentStep      string     `json:"currentStep" db:"current_step"`
	StartTime        time.Time  `json:"startTime" db:"start_time"`
	EndTime          *time.Time `json:"endTime" db:"end_time"`
	Results          string     `json:"results" db:"results"` // JSON serialized GenerationResults
	Errors           string     `json:"errors" db:"errors"`   // JSON serialized GenerationError array
	ProcessingTimeMs int64      `json:"processingTimeMs" db:"processing_time_ms"`
}

// Template represents a template in the database
type Template struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Version     string    `json:"version" db:"version"`
	Author      string    `json:"author" db:"author"`
	Type        string    `json:"type" db:"type"`
	Path        string    `json:"path" db:"path"`
	IsBuiltIn   bool      `json:"isBuiltIn" db:"is_built_in"`
	Variables   string    `json:"variables" db:"variables"` // JSON serialized TemplateVariable array
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// AppSetting represents an application setting in the database
type AppSetting struct {
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"`
	Type      string    `json:"type" db:"type"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// ValidationCache represents a validation cache entry in the database
type ValidationCache struct {
	SpecHash         string    `json:"specHash" db:"spec_hash"`
	SpecPath         string    `json:"specPath" db:"spec_path"`
	SpecURL          string    `json:"specUrl" db:"spec_url"`
	ValidationResult string    `json:"validationResult" db:"validation_result"` // JSON serialized ValidationResult
	CachedAt         time.Time `json:"cachedAt" db:"cached_at"`
	ExpiresAt        time.Time `json:"expiresAt" db:"expires_at"`
}
