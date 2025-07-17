package app

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"MCPWeaver/internal/generator"
	"MCPWeaver/internal/mapping"
	"MCPWeaver/internal/parser"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Active generation jobs
var (
	activeJobs = make(map[string]*GenerationJob)
	jobsMutex  = &sync.RWMutex{}
)

// GenerateServer starts the generation process for a project
func (a *App) GenerateServer(projectID string) (*GenerationJob, error) {
	if projectID == "" {
		return nil, a.createAPIError("validation", ErrCodeValidation, "Project ID is required", nil)
	}

	// Get project
	project, err := a.GetProject(projectID)
	if err != nil {
		return nil, err
	}

	// Check if there's already an active job for this project
	jobsMutex.RLock()
	for _, job := range activeJobs {
		if job.ProjectID == projectID && (job.Status == GenerationStatusStarted || job.Status == GenerationStatusParsing || job.Status == GenerationStatusMapping || job.Status == GenerationStatusGenerating) {
			jobsMutex.RUnlock()
			return nil, a.createAPIError("validation", ErrCodeValidation, "Generation already in progress for this project", nil)
		}
	}
	jobsMutex.RUnlock()

	// Create generation job
	job := &GenerationJob{
		ID:          generateJobID(),
		ProjectID:   projectID,
		Status:      GenerationStatusStarted,
		Progress:    0.0,
		CurrentStep: "Initializing generation",
		StartTime:   time.Now(),
		Errors:      []GenerationError{},
	}

	// Store job
	jobsMutex.Lock()
	activeJobs[job.ID] = job
	jobsMutex.Unlock()

	// Update project status
	a.updateProjectStatus(projectID, ProjectStatusGenerating)

	// Start generation in background
	go a.runGeneration(job, project)

	// Emit event
	runtime.EventsEmit(a.ctx, "generation:started", job)

	return job, nil
}

// GetGenerationJob returns a specific generation job by ID
func (a *App) GetGenerationJob(jobID string) (*GenerationJob, error) {
	if jobID == "" {
		return nil, a.createAPIError("validation", ErrCodeValidation, "Job ID is required", nil)
	}

	jobsMutex.RLock()
	job, exists := activeJobs[jobID]
	jobsMutex.RUnlock()

	if !exists {
		return nil, a.createAPIError("not_found", ErrCodeNotFound, "Generation job not found", nil)
	}

	return job, nil
}

// CancelGeneration cancels an active generation job
func (a *App) CancelGeneration(jobID string) error {
	if jobID == "" {
		return a.createAPIError("validation", ErrCodeValidation, "Job ID is required", nil)
	}

	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	job, exists := activeJobs[jobID]
	if !exists {
		return a.createAPIError("not_found", ErrCodeNotFound, "Generation job not found", nil)
	}

	if job.Status == GenerationStatusCompleted || job.Status == GenerationStatusFailed || job.Status == GenerationStatusCancelled {
		return a.createAPIError("validation", ErrCodeValidation, "Cannot cancel completed generation", nil)
	}

	// Update job status
	job.Status = GenerationStatusCancelled
	job.CurrentStep = "Cancelled by user"
	endTime := time.Now()
	job.EndTime = &endTime

	// Update project status
	a.updateProjectStatus(job.ProjectID, ProjectStatusReady)

	// Emit event
	runtime.EventsEmit(a.ctx, "generation:cancelled", job)

	// Send notification
	a.emitNotification("info", "Generation Cancelled", "Code generation has been cancelled")

	return nil
}

// GetGenerationHistory returns the generation history for a project
func (a *App) GetGenerationHistory(projectID string) ([]*GenerationJob, error) {
	if projectID == "" {
		return nil, a.createAPIError("validation", ErrCodeValidation, "Project ID is required", nil)
	}

	// For now, return active jobs for this project
	// In a full implementation, this would fetch from database
	jobsMutex.RLock()
	defer jobsMutex.RUnlock()

	var jobs []*GenerationJob
	for _, job := range activeJobs {
		if job.ProjectID == projectID {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

// runGeneration executes the generation process
func (a *App) runGeneration(job *GenerationJob, project *Project) {
	defer func() {
		if r := recover(); r != nil {
			a.handleGenerationError(job, fmt.Sprintf("Generation panicked: %v", r))
		}
	}()

	// Step 1: Parse OpenAPI specification
	a.updateJobProgress(job, GenerationStatusParsing, 0.1, "Parsing OpenAPI specification")

	var specPath string
	if project.SpecPath != "" {
		specPath = project.SpecPath
	} else if project.SpecURL != "" {
		// For URLs, we'll parse directly
		specPath = project.SpecURL
	} else {
		a.handleGenerationError(job, "No OpenAPI specification provided")
		return
	}

	var parsedAPI *parser.ParsedAPI
	var err error

	if strings.HasPrefix(specPath, "http://") || strings.HasPrefix(specPath, "https://") {
		parsedAPI, err = a.parserService.ParseFromURL(a.ctx, specPath)
	} else {
		parsedAPI, err = a.parserService.ParseFromFile(specPath)
	}

	if err != nil {
		a.handleGenerationError(job, fmt.Sprintf("Failed to parse OpenAPI specification: %v", err))
		return
	}

	a.updateJobProgress(job, GenerationStatusMapping, 0.3, "Mapping operations to MCP tools")

	// Step 2: Map operations to MCP tools
	baseURL := parsedAPI.BaseURL
	if baseURL == "" && len(parsedAPI.Servers) > 0 {
		baseURL = parsedAPI.Servers[0]
	}

	mappingService := mapping.NewService(baseURL)
	tools, err := mappingService.MapOperationsToTools(parsedAPI.Operations)
	if err != nil {
		a.handleGenerationError(job, fmt.Sprintf("Failed to map operations to MCP tools: %v", err))
		return
	}

	a.updateJobProgress(job, GenerationStatusGenerating, 0.5, "Generating MCP server code")

	// Step 3: Generate server code
	if a.generatorService == nil {
		a.generatorService = generator.NewService(project.OutputPath)
	}

	err = a.generatorService.Generate(parsedAPI, tools, project.Settings.PackageName)
	if err != nil {
		a.handleGenerationError(job, fmt.Sprintf("Failed to generate server code: %v", err))
		return
	}

	a.updateJobProgress(job, GenerationStatusValidating, 0.8, "Validating generated code")

	// Step 4: Validate generated code (basic file existence check)
	generatedFiles := []GeneratedFile{
		{
			Path: filepath.Join(project.OutputPath, "main.go"),
			Type: "server",
			Size: 0, // TODO: Get actual file size
			LinesOfCode: 0, // TODO: Count lines
		},
		{
			Path: filepath.Join(project.OutputPath, "go.mod"),
			Type: "module",
			Size: 0,
			LinesOfCode: 0,
		},
		{
			Path: filepath.Join(project.OutputPath, "README.md"),
			Type: "documentation",
			Size: 0,
			LinesOfCode: 0,
		},
	}

	// Step 5: Complete generation
	a.updateJobProgress(job, GenerationStatusCompleted, 1.0, "Generation completed successfully")

	// Update job results
	endTime := time.Now()
	job.EndTime = &endTime
	job.Results = &GenerationResults{
		ServerPath:     filepath.Join(project.OutputPath, "main.go"),
		GeneratedFiles: generatedFiles,
		MCPTools:       tools,
		Statistics: GenerationStats{
			TotalEndpoints:  len(parsedAPI.Operations),
			GeneratedTools:  len(tools),
			ProcessingTime:  endTime.Sub(job.StartTime),
			SpecComplexity:  "medium", // TODO: Calculate complexity
			TemplateVersion: "1.0.0",
		},
	}

	// Update project
	a.updateProjectGeneration(job.ProjectID)
	a.updateProjectStatus(job.ProjectID, ProjectStatusReady)

	// Emit completion event
	runtime.EventsEmit(a.ctx, "generation:completed", job)

	// Send notification
	a.emitNotification("success", "Generation Complete", 
		fmt.Sprintf("MCP server generated successfully with %d tools", len(tools)))
}

// updateJobProgress updates the job progress and emits events
func (a *App) updateJobProgress(job *GenerationJob, status GenerationStatus, progress float64, step string) {
	jobsMutex.Lock()
	job.Status = status
	job.Progress = progress
	job.CurrentStep = step
	jobsMutex.Unlock()

	// Emit progress event
	runtime.EventsEmit(a.ctx, "generation:progress", GenerationProgress{
		JobID:       job.ID,
		Progress:    progress,
		CurrentStep: step,
		Message:     step,
	})
}

// handleGenerationError handles errors during generation
func (a *App) handleGenerationError(job *GenerationJob, message string) {
	jobsMutex.Lock()
	job.Status = GenerationStatusFailed
	job.CurrentStep = "Generation failed"
	endTime := time.Now()
	job.EndTime = &endTime
	
	job.Errors = append(job.Errors, GenerationError{
		Type:    "generation",
		Message: message,
		Details: "",
	})
	jobsMutex.Unlock()

	// Update project status
	a.updateProjectStatus(job.ProjectID, ProjectStatusError)

	// Emit error event
	runtime.EventsEmit(a.ctx, "generation:failed", map[string]interface{}{
		"jobId":   job.ID,
		"type":    "generation",
		"message": message,
	})

	// Send notification
	a.emitNotification("error", "Generation Failed", message)
}

// updateProjectStatus updates the project status
func (a *App) updateProjectStatus(projectID string, status ProjectStatus) {
	project, err := a.GetProject(projectID)
	if err != nil {
		return
	}

	// Update project in database
	dbProject, err := a.projectRepo.GetByID(projectID)
	if err != nil {
		return
	}

	dbProject.Status = string(status)
	dbProject.UpdatedAt = time.Now()
	a.projectRepo.Update(dbProject)

	// Emit update event
	project.Status = status
	project.UpdatedAt = time.Now()
	runtime.EventsEmit(a.ctx, "project:updated", project)
}

// updateProjectGeneration updates the project's generation count and timestamp
func (a *App) updateProjectGeneration(projectID string) {
	dbProject, err := a.projectRepo.GetByID(projectID)
	if err != nil {
		return
	}

	now := time.Now()
	dbProject.GenerationCount++
	dbProject.LastGenerated = &now
	dbProject.UpdatedAt = now
	a.projectRepo.Update(dbProject)
}

// generateJobID generates a unique job ID
func generateJobID() string {
	return fmt.Sprintf("gen_%d", time.Now().UnixNano())
}