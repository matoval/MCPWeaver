package unit

import (
	"encoding/json"
	"testing"
	"time"

	"MCPWeaver/internal/app"
	"MCPWeaver/internal/mapping"
	"MCPWeaver/tests/helpers"

	"github.com/stretchr/testify/suite"
)

type TypesTestSuite struct {
	suite.Suite
	helper *helpers.TestHelper
}

func (s *TypesTestSuite) SetupTest() {
	s.helper = helpers.NewTestHelper(s.T())
}

func TestTypesTestSuite(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}

func (s *TypesTestSuite) TestProjectStatus_Constants() {
	// Test project status constants are properly defined
	statuses := []app.ProjectStatus{
		app.ProjectStatusCreated,
		app.ProjectStatusValidating,
		app.ProjectStatusReady,
		app.ProjectStatusGenerating,
		app.ProjectStatusError,
	}

	for _, status := range statuses {
		s.helper.AssertNotEqual("", string(status))
	}
}

func (s *TypesTestSuite) TestGenerationStatus_Constants() {
	// Test generation status constants
	statuses := []app.GenerationStatus{
		app.GenerationStatusStarted,
		app.GenerationStatusParsing,
		app.GenerationStatusMapping,
		app.GenerationStatusGenerating,
		app.GenerationStatusValidating,
		app.GenerationStatusCompleted,
		app.GenerationStatusFailed,
		app.GenerationStatusCancelled,
	}

	for _, status := range statuses {
		s.helper.AssertNotEqual("", string(status))
	}
}

func (s *TypesTestSuite) TestErrorSeverity_Constants() {
	// Test error severity constants
	severities := []app.ErrorSeverity{
		app.ErrorSeverityLow,
		app.ErrorSeverityMedium,
		app.ErrorSeverityHigh,
		app.ErrorSeverityCritical,
	}

	for _, severity := range severities {
		s.helper.AssertNotEqual("", string(severity))
	}
}

func (s *TypesTestSuite) TestCreateProjectRequest_Validation() {
	request := app.CreateProjectRequest{
		Name:       "Test Project",
		SpecPath:   "/path/to/spec.yaml",
		OutputPath: "/tmp/output",
		Settings: app.ProjectSettings{
			PackageName:   "test-server",
			ServerPort:    8080,
			EnableLogging: true,
			LogLevel:      "info",
		},
	}

	s.helper.AssertEqual("Test Project", request.Name)
	s.helper.AssertEqual("/path/to/spec.yaml", request.SpecPath)
	s.helper.AssertEqual("/tmp/output", request.OutputPath)
	s.helper.AssertEqual("test-server", request.Settings.PackageName)
	s.helper.AssertEqual(8080, request.Settings.ServerPort)
}

func (s *TypesTestSuite) TestProjectSettings_JSONSerialization() {
	settings := app.ProjectSettings{
		PackageName:     "test-server",
		ServerPort:      8080,
		EnableLogging:   true,
		LogLevel:        "debug",
		CustomTemplates: []string{"template1", "template2"},
	}

	// Test marshaling
	jsonData, err := json.Marshal(settings)
	s.helper.AssertNoError(err)
	s.helper.AssertContains(string(jsonData), "test-server")

	// Test unmarshaling
	var deserializedSettings app.ProjectSettings
	err = json.Unmarshal(jsonData, &deserializedSettings)
	s.helper.AssertNoError(err)
	s.helper.AssertEqual(settings.PackageName, deserializedSettings.PackageName)
	s.helper.AssertEqual(settings.ServerPort, deserializedSettings.ServerPort)
	s.helper.AssertEqual(len(settings.CustomTemplates), len(deserializedSettings.CustomTemplates))
}

func (s *TypesTestSuite) TestGenerationJob_Structure() {
	job := app.GenerationJob{
		ID:          "gen_123",
		ProjectID:   "project_456",
		Status:      app.GenerationStatusStarted,
		Progress:    0.5,
		CurrentStep: "Processing",
		StartTime:   time.Now(),
		Errors:      []app.GenerationError{},
		Warnings:    []string{},
	}

	s.helper.AssertEqual("gen_123", job.ID)
	s.helper.AssertEqual("project_456", job.ProjectID)
	s.helper.AssertEqual(app.GenerationStatusStarted, job.Status)
	s.helper.AssertEqual(0.5, job.Progress)
	s.helper.AssertEqual("Processing", job.CurrentStep)
	s.helper.AssertEqual(0, len(job.Errors))
	s.helper.AssertEqual(0, len(job.Warnings))
}

func (s *TypesTestSuite) TestAPIError_Structure() {
	apiError := app.APIError{
		Type:        "validation",
		Code:        "INVALID_INPUT",
		Message:     "Test error",
		Severity:    app.ErrorSeverityHigh,
		Recoverable: true,
		Timestamp:   time.Now(),
	}

	s.helper.AssertEqual("validation", apiError.Type)
	s.helper.AssertEqual("INVALID_INPUT", apiError.Code)
	s.helper.AssertEqual("Test error", apiError.Message)
	s.helper.AssertEqual(app.ErrorSeverityHigh, apiError.Severity)
	s.helper.AssertEqual(true, apiError.Recoverable)
}

func (s *TypesTestSuite) TestAPIError_ErrorInterface() {
	apiError := app.APIError{
		Message: "Test error message",
	}

	// Test Error() method
	s.helper.AssertEqual("Test error message", apiError.Error())
}

func (s *TypesTestSuite) TestAPIError_IsRetryable() {
	retryDelay := 5 * time.Second

	// Retryable error
	retryableError := app.APIError{
		Recoverable: true,
		RetryAfter:  &retryDelay,
	}
	s.helper.AssertEqual(true, retryableError.IsRetryable())
	s.helper.AssertEqual(retryDelay, retryableError.GetRetryDelay())

	// Non-retryable error
	nonRetryableError := app.APIError{
		Recoverable: false,
	}
	s.helper.AssertEqual(false, nonRetryableError.IsRetryable())
	s.helper.AssertEqual(time.Duration(0), nonRetryableError.GetRetryDelay())
}

func (s *TypesTestSuite) TestErrorCollection_Methods() {
	collection := app.ErrorCollection{
		Errors: []app.APIError{
			{Message: "Error 1"},
			{Message: "Error 2"},
		},
		Warnings: []app.APIError{
			{Message: "Warning 1"},
		},
		Operation: "test_operation",
	}

	s.helper.AssertEqual(true, collection.HasErrors())
	s.helper.AssertEqual(true, collection.HasWarnings())
	s.helper.AssertContains(collection.Error(), "2 errors")
}

func (s *TypesTestSuite) TestDefaultRetryPolicy() {
	policy := app.DefaultRetryPolicy()

	s.helper.AssertEqual(3, policy.MaxRetries)
	s.helper.AssertEqual(time.Second, policy.InitialDelay)
	s.helper.AssertEqual(30*time.Second, policy.MaxDelay)
	s.helper.AssertEqual(2.0, policy.BackoffMultiplier)
	s.helper.AssertEqual(true, policy.JitterEnabled)
	s.helper.AssertEqual(true, len(policy.RetryableErrors) > 0)
}

func (s *TypesTestSuite) TestValidationResult_Structure() {
	result := app.ValidationResult{
		Valid:          true,
		Errors:         []app.ValidationError{},
		Warnings:       []app.ValidationWarning{},
		Suggestions:    []string{"suggestion1"},
		ValidationTime: time.Millisecond * 100,
		CacheHit:       false,
		ValidatedAt:    time.Now(),
	}

	s.helper.AssertEqual(true, result.Valid)
	s.helper.AssertEqual(0, len(result.Errors))
	s.helper.AssertEqual(0, len(result.Warnings))
	s.helper.AssertEqual(1, len(result.Suggestions))
	s.helper.AssertEqual(time.Millisecond*100, result.ValidationTime)
	s.helper.AssertEqual(false, result.CacheHit)
}

func (s *TypesTestSuite) TestImportResult_Structure() {
	result := app.ImportResult{
		Content:      "openapi: 3.0.0",
		Valid:        true,
		ImportedFrom: "file",
		FilePath:     "/path/to/spec.yaml",
		FileSize:     1024,
		ImportedAt:   time.Now(),
	}

	s.helper.AssertEqual("openapi: 3.0.0", result.Content)
	s.helper.AssertEqual(true, result.Valid)
	s.helper.AssertEqual("file", result.ImportedFrom)
	s.helper.AssertEqual("/path/to/spec.yaml", result.FilePath)
	s.helper.AssertEqual(int64(1024), result.FileSize)
}

func (s *TypesTestSuite) TestGenerationResults_Structure() {
	files := []app.GeneratedFile{
		{Path: "/path/main.go", Type: "server", Size: 1024, LinesOfCode: 50},
	}

	tools := []mapping.MCPTool{
		{Name: "get_users", Description: "Get users"},
	}

	stats := app.GenerationStats{
		TotalEndpoints:  5,
		GeneratedTools:  1,
		ProcessingTime:  time.Second,
		SpecComplexity:  "low",
		TemplateVersion: "1.0.0",
	}

	results := app.GenerationResults{
		ServerPath:     "/path/main.go",
		GeneratedFiles: files,
		MCPTools:       tools,
		Statistics:     stats,
	}

	s.helper.AssertEqual("/path/main.go", results.ServerPath)
	s.helper.AssertEqual(1, len(results.GeneratedFiles))
	s.helper.AssertEqual(1, len(results.MCPTools))
	s.helper.AssertEqual(5, results.Statistics.TotalEndpoints)
}

func (s *TypesTestSuite) TestTemplateTypes() {
	// Test template type constants
	types := []app.TemplateType{
		app.TemplateTypeDefault,
		app.TemplateTypeCustom,
		app.TemplateTypePlugin,
	}

	for _, templateType := range types {
		s.helper.AssertNotEqual("", string(templateType))
	}
}

func (s *TypesTestSuite) TestAppSettings_Structure() {
	settings := app.AppSettings{
		Theme:             "dark",
		Language:          "en",
		AutoSave:          true,
		DefaultOutputPath: "/tmp",
		RecentProjects:    []string{"proj1", "proj2"},
		RecentFiles:       []string{"file1", "file2"},
		WindowSettings: app.WindowSettings{
			Width:     1280,
			Height:    720,
			Maximized: false,
		},
		EditorSettings: app.EditorSettings{
			FontSize:        14,
			FontFamily:      "Monaco",
			TabSize:         4,
			WordWrap:        true,
			LineNumbers:     true,
			SyntaxHighlight: true,
		},
	}

	s.helper.AssertEqual("dark", settings.Theme)
	s.helper.AssertEqual("en", settings.Language)
	s.helper.AssertEqual(true, settings.AutoSave)
	s.helper.AssertEqual(2, len(settings.RecentProjects))
	s.helper.AssertEqual(1280, settings.WindowSettings.Width)
	s.helper.AssertEqual(14, settings.EditorSettings.FontSize)
}

// Performance Tests
func (s *TypesTestSuite) TestTypeSerialization_Performance() {
	project := app.Project{
		ID:         "test-project",
		Name:       "Test Project",
		SpecPath:   "/path/to/spec.yaml",
		OutputPath: "/tmp/output",
		Settings: app.ProjectSettings{
			PackageName:     "test-server",
			ServerPort:      8080,
			EnableLogging:   true,
			LogLevel:        "info",
			CustomTemplates: make([]string, 100),
		},
		Status:          app.ProjectStatusReady,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		GenerationCount: 5,
	}

	s.helper.AssertPerformance(func() {
		for i := 0; i < 1000; i++ {
			jsonData, err := json.Marshal(project)
			s.helper.AssertNoError(err)

			var deserializedProject app.Project
			err = json.Unmarshal(jsonData, &deserializedProject)
			s.helper.AssertNoError(err)
		}
	}, 50*1000000) // 50ms
}
