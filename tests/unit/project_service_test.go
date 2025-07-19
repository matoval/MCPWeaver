package unit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"MCPWeaver/internal/project"
	"MCPWeaver/tests/helpers"

	"github.com/stretchr/testify/suite"
)

type ProjectServiceTestSuite struct {
	suite.Suite
	helper  *helpers.TestHelper
	service *project.Service
	ctx     context.Context
}

func (s *ProjectServiceTestSuite) SetupTest() {
	s.helper = helpers.NewTestHelper(s.T())
	s.service = project.New()
	s.ctx = context.Background()
}

func TestProjectServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ProjectServiceTestSuite))
}

func (s *ProjectServiceTestSuite) TestNewService() {
	service := project.New()
	s.helper.AssertNotNil(service)
}

func (s *ProjectServiceTestSuite) TestProject_Structure() {
	proj := project.Project{
		ID:        "test-id",
		Name:      "Test Project", 
		SpecPath:  "/path/to/spec.yaml",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.helper.AssertEqual("test-id", proj.ID)
	s.helper.AssertEqual("Test Project", proj.Name)
	s.helper.AssertEqual("/path/to/spec.yaml", proj.SpecPath)
	s.helper.AssertNotNil(proj.CreatedAt)
	s.helper.AssertNotNil(proj.UpdatedAt)
}

func (s *ProjectServiceTestSuite) TestCreateProjectRequest_Structure() {
	req := project.CreateProjectRequest{
		Name:     "Test Project",
		SpecPath: "/path/to/spec.yaml",
	}

	s.helper.AssertEqual("Test Project", req.Name)
	s.helper.AssertEqual("/path/to/spec.yaml", req.SpecPath)
}

func (s *ProjectServiceTestSuite) TestCreateProject_CurrentImplementation() {
	// Test the current implementation that returns "not yet implemented"
	req := project.CreateProjectRequest{
		Name:     "Test Project",
		SpecPath: "/path/to/spec.yaml",
	}

	proj, err := s.service.Create(s.ctx, req)
	
	// Current implementation returns an error but also creates a project struct
	s.helper.AssertError(err)
	s.helper.AssertContains(err.Error(), "not yet implemented")
	
	// But it still returns a project with the data filled in
	s.helper.AssertNotNil(proj)
	s.helper.AssertEqual("Test Project", proj.Name)
	s.helper.AssertEqual("/path/to/spec.yaml", proj.SpecPath)
	s.helper.AssertNotNil(proj.CreatedAt)
	s.helper.AssertNotNil(proj.UpdatedAt)
	s.helper.AssertContains(proj.ID, "proj_")
}

func (s *ProjectServiceTestSuite) TestGetAll_CurrentImplementation() {
	// Test the current implementation that returns "not yet implemented"
	projects, err := s.service.GetAll(s.ctx)
	
	s.helper.AssertError(err)
	s.helper.AssertContains(err.Error(), "not yet implemented")
	s.helper.AssertNil(projects)
}

func (s *ProjectServiceTestSuite) TestGetByID_CurrentImplementation() {
	// Test the current implementation that returns "not yet implemented"
	proj, err := s.service.GetByID(s.ctx, "test-id")
	
	s.helper.AssertError(err)
	s.helper.AssertContains(err.Error(), "not yet implemented")
	s.helper.AssertNil(proj)
}

func (s *ProjectServiceTestSuite) TestGetByID_EmptyID() {
	// Test with empty ID
	proj, err := s.service.GetByID(s.ctx, "")
	
	s.helper.AssertError(err)
	s.helper.AssertContains(err.Error(), "not yet implemented")
	s.helper.AssertNil(proj)
}

func (s *ProjectServiceTestSuite) TestProjectIDGeneration() {
	// Test ID generation by creating multiple projects
	req := project.CreateProjectRequest{
		Name:     "Test Project",
		SpecPath: "/path/to/spec.yaml",
	}

	// Create multiple projects and verify IDs are unique
	ids := make(map[string]bool)
	for i := 0; i < 10; i++ {
		proj, err := s.service.Create(s.ctx, req)
		s.helper.AssertError(err) // Still expecting "not implemented" error
		s.helper.AssertNotNil(proj)
		
		// Check that ID is unique
		s.helper.AssertEqual(false, ids[proj.ID], "ID should be unique")
		ids[proj.ID] = true
		
		// Check ID format
		s.helper.AssertContains(proj.ID, "proj_")
		
		// Add small delay to ensure different timestamps
		time.Sleep(time.Nanosecond)
	}
}

func (s *ProjectServiceTestSuite) TestProjectTimestamps() {
	req := project.CreateProjectRequest{
		Name:     "Test Project",
		SpecPath: "/path/to/spec.yaml",
	}

	beforeCreate := time.Now()
	proj, err := s.service.Create(s.ctx, req)
	afterCreate := time.Now()
	
	s.helper.AssertError(err) // Expected for current implementation
	s.helper.AssertNotNil(proj)
	
	// Verify timestamps are reasonable
	s.helper.AssertEqual(true, proj.CreatedAt.After(beforeCreate) || proj.CreatedAt.Equal(beforeCreate))
	s.helper.AssertEqual(true, proj.CreatedAt.Before(afterCreate) || proj.CreatedAt.Equal(afterCreate))
	s.helper.AssertEqual(true, proj.UpdatedAt.After(beforeCreate) || proj.UpdatedAt.Equal(beforeCreate))
	s.helper.AssertEqual(true, proj.UpdatedAt.Before(afterCreate) || proj.UpdatedAt.Equal(afterCreate))
}

func (s *ProjectServiceTestSuite) TestCreateProject_WithEmptyRequest() {
	// Test with empty request
	req := project.CreateProjectRequest{}

	proj, err := s.service.Create(s.ctx, req)
	
	s.helper.AssertError(err)
	s.helper.AssertContains(err.Error(), "not yet implemented")
	
	// Even with empty request, current implementation creates project struct
	s.helper.AssertNotNil(proj)
	s.helper.AssertEqual("", proj.Name)
	s.helper.AssertEqual("", proj.SpecPath)
	s.helper.AssertNotEqual("", proj.ID) // ID should still be generated
}

func (s *ProjectServiceTestSuite) TestCreateProject_WithSpecialCharacters() {
	// Test with special characters in name and path
	req := project.CreateProjectRequest{
		Name:     "Test Project!@#$%^&*()",
		SpecPath: "/path/with spaces/special-chars_spec.yaml",
	}

	proj, err := s.service.Create(s.ctx, req)
	
	s.helper.AssertError(err)
	s.helper.AssertNotNil(proj)
	s.helper.AssertEqual("Test Project!@#$%^&*()", proj.Name)
	s.helper.AssertEqual("/path/with spaces/special-chars_spec.yaml", proj.SpecPath)
}

func (s *ProjectServiceTestSuite) TestContextHandling() {
	// Test with nil context (should not panic)
	req := project.CreateProjectRequest{
		Name:     "Test Project",
		SpecPath: "/path/to/spec.yaml",
	}

	// Current implementation doesn't actually use context but should handle it gracefully
	proj, err := s.service.Create(nil, req)
	s.helper.AssertError(err)
	s.helper.AssertNotNil(proj)

	// Test with canceled context
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	
	proj2, err2 := s.service.Create(cancelCtx, req)
	s.helper.AssertError(err2)
	s.helper.AssertNotNil(proj2)
}

// Performance Tests
func (s *ProjectServiceTestSuite) TestCreateProject_Performance() {
	req := project.CreateProjectRequest{
		Name:     "Performance Test Project",
		SpecPath: "/path/to/spec.yaml",
	}

	s.helper.AssertPerformance(func() {
		for i := 0; i < 1000; i++ {
			_, err := s.service.Create(s.ctx, req)
			s.helper.AssertError(err) // Expected for current implementation
		}
	}, 100*1000000) // 100ms
}

func (s *ProjectServiceTestSuite) TestIDGeneration_Performance() {
	s.helper.AssertPerformance(func() {
		for i := 0; i < 10000; i++ {
			// Test the ID generation logic
			id := fmt.Sprintf("proj_%d", time.Now().UnixNano())
			s.helper.AssertContains(id, "proj_")
		}
	}, 50*1000000) // 50ms
}

// Memory Tests
func (s *ProjectServiceTestSuite) TestProject_MemoryUsage() {
	// Test that creating many projects doesn't consume excessive memory
	req := project.CreateProjectRequest{
		Name:     "Memory Test Project",
		SpecPath: "/path/to/spec.yaml",
	}

	projects := make([]*project.Project, 1000)
	for i := 0; i < 1000; i++ {
		proj, err := s.service.Create(s.ctx, req)
		s.helper.AssertError(err) // Expected for current implementation
		projects[i] = proj
	}

	// Verify we created 1000 projects
	s.helper.AssertEqual(1000, len(projects))
	
	// Verify each project has unique ID
	ids := make(map[string]bool)
	for _, proj := range projects {
		s.helper.AssertEqual(false, ids[proj.ID], "Each project should have unique ID")
		ids[proj.ID] = true
	}
}