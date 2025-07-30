package project

import (
	"context"
	"time"

	"MCPWeaver/internal/app"
	"MCPWeaver/internal/database"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Service handles project management operations
type Service struct {
	projectRepo *database.ProjectRepository
	errorManager *app.ErrorManager
	ctx         context.Context
}

// Project represents a project
type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	SpecPath  string    `json:"specPath"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateProjectRequest represents a request to create a project
type CreateProjectRequest struct {
	Name     string `json:"name"`
	SpecPath string `json:"specPath"`
}

// New creates a new project service
func New(projectRepo *database.ProjectRepository, errorManager *app.ErrorManager, ctx context.Context) *Service {
	return &Service{
		projectRepo:  projectRepo,
		errorManager: errorManager,
		ctx:          ctx,
	}
}

// Create creates a new project
func (s *Service) Create(ctx context.Context, req CreateProjectRequest) (*Project, error) {
	// Validate input
	if req.Name == "" {
		return nil, s.errorManager.CreateValidationError(
			"Project name is required",
			map[string]string{"field": "name"},
			[]string{"Provide a valid project name"},
		)
	}

	if req.SpecPath == "" {
		return nil, s.errorManager.CreateValidationError(
			"Spec path is required",
			map[string]string{"field": "specPath"},
			[]string{"Provide a valid OpenAPI specification path"},
		)
	}

	// Generate unique project ID
	projectID := uuid.New().String()

	// Create database project model
	dbProject := &database.Project{
		ID:              projectID,
		Name:            req.Name,
		SpecPath:        req.SpecPath,
		SpecURL:         "", // Empty for file-based specs
		OutputPath:      "", // Will be set based on settings later
		Settings:        "{}", // Empty JSON object for now
		Status:          "created",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		LastGenerated:   nil,
		GenerationCount: 0,
	}

	// Save to database
	if err := s.projectRepo.Create(dbProject); err != nil {
		return nil, s.errorManager.CreateDatabaseError(
			"Failed to create project in database",
			"create",
			"projects",
		)
	}

	// Convert to service model
	project := &Project{
		ID:        dbProject.ID,
		Name:      dbProject.Name,
		SpecPath:  dbProject.SpecPath,
		CreatedAt: dbProject.CreatedAt,
		UpdatedAt: dbProject.UpdatedAt,
	}

	// Emit project created event
	if s.ctx != nil {
		runtime.EventsEmit(s.ctx, "project:created", map[string]any{
			"projectId": project.ID,
			"name":      project.Name,
			"timestamp": time.Now(),
		})
	}

	return project, nil
}

// GetAll returns all projects
func (s *Service) GetAll(ctx context.Context) ([]*Project, error) {
	// Retrieve projects from database
	dbProjects, err := s.projectRepo.GetAll()
	if err != nil {
		return nil, s.errorManager.CreateDatabaseError(
			"Failed to retrieve projects from database",
			"getAll",
			"projects",
		)
	}

	// Convert database models to service models
	projects := make([]*Project, 0, len(dbProjects))
	for _, dbProject := range dbProjects {
		project := &Project{
			ID:        dbProject.ID,
			Name:      dbProject.Name,
			SpecPath:  dbProject.SpecPath,
			CreatedAt: dbProject.CreatedAt,
			UpdatedAt: dbProject.UpdatedAt,
		}
		projects = append(projects, project)
	}

	return projects, nil
}

// GetByID returns a project by ID
func (s *Service) GetByID(ctx context.Context, id string) (*Project, error) {
	// Validate input
	if id == "" {
		return nil, s.errorManager.CreateValidationError(
			"Project ID is required",
			map[string]string{"field": "id"},
			[]string{"Provide a valid project ID"},
		)
	}

	// Retrieve project from database
	dbProject, err := s.projectRepo.GetByID(id)
	if err != nil {
		if err.Error() == "project not found" {
			return nil, s.errorManager.CreateValidationError(
				"Project not found",
				map[string]string{"projectId": id},
				[]string{"Check that the project ID is correct"},
			)
		}
		return nil, s.errorManager.CreateDatabaseError(
			"Failed to retrieve project from database",
			"getByID",
			"projects",
		)
	}

	// Convert database model to service model
	project := &Project{
		ID:        dbProject.ID,
		Name:      dbProject.Name,
		SpecPath:  dbProject.SpecPath,
		CreatedAt: dbProject.CreatedAt,
		UpdatedAt: dbProject.UpdatedAt,
	}

	return project, nil
}

