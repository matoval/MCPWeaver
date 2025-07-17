package app

import (
	"encoding/json"
	"fmt"
	"time"

	"MCPWeaver/internal/database"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// CreateProject creates a new project with the given configuration
func (a *App) CreateProject(request CreateProjectRequest) (*Project, error) {
	// Validate input
	if request.Name == "" {
		return nil, a.createAPIError("validation", ErrCodeValidation, "Project name is required", nil)
	}
	
	if request.OutputPath == "" {
		return nil, a.createAPIError("validation", ErrCodeValidation, "Output path is required", nil)
	}

	// Check if project with same name already exists
	existing, err := a.projectRepo.GetByName(request.Name)
	if err != nil && err.Error() != "project not found" {
		return nil, a.createAPIError("internal", ErrCodeInternalError, "Failed to check existing project", nil)
	}
	if existing != nil {
		return nil, a.createAPIError("validation", ErrCodeValidation, "Project with this name already exists", nil)
	}

	// Serialize settings to JSON
	settingsJSON, err := json.Marshal(request.Settings)
	if err != nil {
		return nil, a.createAPIError("internal", ErrCodeInternalError, "Failed to serialize settings", map[string]string{
			"error": err.Error(),
		})
	}

	// Create database project
	dbProject := &database.Project{
		ID:              fmt.Sprintf("proj_%d", time.Now().UnixNano()),
		Name:            request.Name,
		SpecPath:        request.SpecPath,
		SpecURL:         request.SpecURL,
		OutputPath:      request.OutputPath,
		Settings:        string(settingsJSON),
		Status:          string(ProjectStatusCreated),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		GenerationCount: 0,
	}

	// Save to database
	err = a.projectRepo.Create(dbProject)
	if err != nil {
		return nil, a.createAPIError("internal", ErrCodeInternalError, "Failed to create project", map[string]string{
			"error": err.Error(),
		})
	}

	// Convert to API project
	project := a.dbProjectToAPI(dbProject)

	// Add to recent projects
	a.addToRecentProjects(project.ID)

	// Emit event
	runtime.EventsEmit(a.ctx, "project:created", project)

	// Send notification
	a.emitNotification("success", "Project Created", fmt.Sprintf("Project '%s' has been created successfully", project.Name))

	return project, nil
}

// GetProjects returns all projects
func (a *App) GetProjects() ([]*Project, error) {
	dbProjects, err := a.projectRepo.GetAll()
	if err != nil {
		return nil, a.createAPIError("internal", ErrCodeInternalError, "Failed to retrieve projects", map[string]string{
			"error": err.Error(),
		})
	}

	projects := make([]*Project, len(dbProjects))
	for i, dbProject := range dbProjects {
		projects[i] = a.dbProjectToAPI(dbProject)
	}

	return projects, nil
}

// GetProject returns a specific project by ID
func (a *App) GetProject(id string) (*Project, error) {
	if id == "" {
		return nil, a.createAPIError("validation", ErrCodeValidation, "Project ID is required", nil)
	}

	dbProject, err := a.projectRepo.GetByID(id)
	if err != nil {
		if err.Error() == "project not found" {
			return nil, a.createAPIError("not_found", ErrCodeNotFound, "Project not found", nil)
		}
		return nil, a.createAPIError("internal", ErrCodeInternalError, "Failed to retrieve project", map[string]string{
			"error": err.Error(),
		})
	}

	return a.dbProjectToAPI(dbProject), nil
}

// UpdateProject updates an existing project
func (a *App) UpdateProject(id string, updates ProjectUpdateRequest) (*Project, error) {
	if id == "" {
		return nil, a.createAPIError("validation", ErrCodeValidation, "Project ID is required", nil)
	}

	// Get existing project
	dbProject, err := a.projectRepo.GetByID(id)
	if err != nil {
		if err.Error() == "project not found" {
			return nil, a.createAPIError("not_found", ErrCodeNotFound, "Project not found", nil)
		}
		return nil, a.createAPIError("internal", ErrCodeInternalError, "Failed to retrieve project", map[string]string{
			"error": err.Error(),
		})
	}

	// Apply updates
	if updates.Name != nil {
		dbProject.Name = *updates.Name
	}
	if updates.SpecPath != nil {
		dbProject.SpecPath = *updates.SpecPath
	}
	if updates.SpecURL != nil {
		dbProject.SpecURL = *updates.SpecURL
	}
	if updates.OutputPath != nil {
		dbProject.OutputPath = *updates.OutputPath
	}
	if updates.Settings != nil {
		// Serialize updated settings to JSON
		settingsJSON, err := json.Marshal(updates.Settings)
		if err != nil {
			return nil, a.createAPIError("internal", ErrCodeInternalError, "Failed to serialize settings", map[string]string{
				"error": err.Error(),
			})
		}
		dbProject.Settings = string(settingsJSON)
	}

	dbProject.UpdatedAt = time.Now()

	// Save to database
	err = a.projectRepo.Update(dbProject)
	if err != nil {
		return nil, a.createAPIError("internal", ErrCodeInternalError, "Failed to update project", map[string]string{
			"error": err.Error(),
		})
	}

	// Convert to API project
	project := a.dbProjectToAPI(dbProject)

	// Emit event
	runtime.EventsEmit(a.ctx, "project:updated", project)

	// Send notification
	a.emitNotification("success", "Project Updated", fmt.Sprintf("Project '%s' has been updated successfully", project.Name))

	return project, nil
}

// DeleteProject deletes a project by ID
func (a *App) DeleteProject(id string) error {
	if id == "" {
		return a.createAPIError("validation", ErrCodeValidation, "Project ID is required", nil)
	}

	// Get project to check if it exists
	dbProject, err := a.projectRepo.GetByID(id)
	if err != nil {
		if err.Error() == "project not found" {
			return a.createAPIError("not_found", ErrCodeNotFound, "Project not found", nil)
		}
		return a.createAPIError("internal", ErrCodeInternalError, "Failed to retrieve project", map[string]string{
			"error": err.Error(),
		})
	}

	// Delete from database
	err = a.projectRepo.Delete(id)
	if err != nil {
		return a.createAPIError("internal", ErrCodeInternalError, "Failed to delete project", map[string]string{
			"error": err.Error(),
		})
	}

	// Remove from recent projects
	a.removeFromRecentProjects(id)

	// Emit event
	runtime.EventsEmit(a.ctx, "project:deleted", id)

	// Send notification
	a.emitNotification("success", "Project Deleted", fmt.Sprintf("Project '%s' has been deleted successfully", dbProject.Name))

	return nil
}

// GetRecentProjects returns recently accessed projects
func (a *App) GetRecentProjects() ([]*Project, error) {
	projects := make([]*Project, 0)
	
	for _, projectID := range a.settings.RecentProjects {
		project, err := a.GetProject(projectID)
		if err == nil {
			projects = append(projects, project)
		}
	}

	return projects, nil
}

// dbProjectToAPI converts a database project to API project
func (a *App) dbProjectToAPI(dbProject *database.Project) *Project {
	var lastGenerated *time.Time
	if dbProject.LastGenerated != nil {
		lastGenerated = dbProject.LastGenerated
	}

	// Deserialize settings from JSON
	var settings ProjectSettings
	if err := json.Unmarshal([]byte(dbProject.Settings), &settings); err != nil {
		// Use default settings if deserialization fails
		settings = ProjectSettings{
			PackageName:     "generated-server",
			ServerPort:      8080,
			EnableLogging:   true,
			LogLevel:        "info",
			CustomTemplates: []string{},
		}
	}

	return &Project{
		ID:              dbProject.ID,
		Name:            dbProject.Name,
		SpecPath:        dbProject.SpecPath,
		SpecURL:         dbProject.SpecURL,
		OutputPath:      dbProject.OutputPath,
		Settings:        settings,
		Status:          ProjectStatus(dbProject.Status),
		CreatedAt:       dbProject.CreatedAt,
		UpdatedAt:       dbProject.UpdatedAt,
		LastGenerated:   lastGenerated,
		GenerationCount: dbProject.GenerationCount,
	}
}

// addToRecentProjects adds a project to the recent projects list
func (a *App) addToRecentProjects(projectID string) {
	// Remove if already exists
	a.removeFromRecentProjects(projectID)
	
	// Add to beginning
	a.settings.RecentProjects = append([]string{projectID}, a.settings.RecentProjects...)
	
	// Keep only last 10
	if len(a.settings.RecentProjects) > 10 {
		a.settings.RecentProjects = a.settings.RecentProjects[:10]
	}
}

// removeFromRecentProjects removes a project from the recent projects list
func (a *App) removeFromRecentProjects(projectID string) {
	for i, id := range a.settings.RecentProjects {
		if id == projectID {
			a.settings.RecentProjects = append(a.settings.RecentProjects[:i], a.settings.RecentProjects[i+1:]...)
			break
		}
	}
}