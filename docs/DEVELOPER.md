# MCPWeaver Developer Guide

This guide provides comprehensive information for developers who want to contribute to MCPWeaver, understand its architecture, or extend its functionality.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Architecture Overview](#architecture-overview)
- [Code Organization](#code-organization)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Building and Packaging](#building-and-packaging)
- [Contributing Guidelines](#contributing-guidelines)
- [Code Style and Standards](#code-style-and-standards)
- [Debugging](#debugging)
- [Performance Optimization](#performance-optimization)
- [Security Considerations](#security-considerations)
- [Release Process](#release-process)

## Getting Started

### Prerequisites

Before you begin development, ensure you have the following tools installed:

**Required:**
- **Go 1.23 or later** - [Download Go](https://golang.org/dl/)
- **Node.js 18 or later** - [Download Node.js](https://nodejs.org/)
- **Wails CLI v2.10.1** - `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- **Git** - [Download Git](https://git-scm.com/)

**Platform-Specific:**

**Windows:**
- Visual Studio 2019 Community or later
- Windows 10 SDK
- WebView2 Evergreen Standalone Installer

**macOS:**
- Xcode Command Line Tools: `xcode-select --install`
- macOS 10.15 SDK or later

**Linux:**
- GCC or Clang compiler
- GTK3 development libraries
- WebKit2GTK development libraries

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install build-essential libgtk-3-dev libwebkit2gtk-4.0-dev

# Fedora/CentOS/RHEL
sudo dnf install gcc-c++ gtk3-devel webkit2gtk3-devel

# Arch Linux
sudo pacman -S base-devel gtk3 webkit2gtk
```

### Quick Setup

1. **Clone the Repository**
   ```bash
   git clone https://github.com/matoval/MCPWeaver.git
   cd MCPWeaver
   ```

2. **Install Dependencies**
   ```bash
   # Install Go dependencies
   go mod tidy
   
   # Install Wails CLI (if not already installed)
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   
   # Install frontend dependencies
   cd frontend
   npm install
   cd ..
   ```

3. **Verify Setup**
   ```bash
   # Check Wails installation
   wails doctor
   
   # Run development server
   wails dev
   ```

## Development Setup

### IDE Configuration

**Visual Studio Code (Recommended):**

Install recommended extensions:
- Go (Google)
- Wails Snippets
- TypeScript and JavaScript Language Features
- React snippets
- GitLens

Workspace settings (`.vscode/settings.json`):
```json
{
  "go.useLanguageServer": true,
  "go.toolsManagement.autoUpdate": true,
  "typescript.preferences.quoteStyle": "double",
  "editor.formatOnSave": true,
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--fast"]
}
```

**GoLand/IntelliJ IDEA:**
- Enable Go plugin
- Configure Wails run configurations
- Set up code style according to gofmt

### Environment Variables

Create a `.env.development` file in the project root:

```bash
# Development settings
WAILS_DEV_MODE=true
WAILS_LOG_LEVEL=debug
MCPWEAVER_DB_PATH=./dev.db
MCPWEAVER_LOG_LEVEL=debug
MCPWEAVER_DEBUG_MODE=true

# Frontend development
REACT_APP_ENV=development
REACT_APP_API_URL=http://localhost:34115
```

### Database Setup

MCPWeaver uses SQLite for local storage. During development:

```bash
# The database will be created automatically
# For testing, you can reset it by deleting:
rm -f mcpweaver.db dev.db

# Database migrations are handled automatically on startup
```

## Architecture Overview

MCPWeaver follows a clean architecture pattern with clear separation of concerns:

```
┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │    Backend      │
│   (React/TS)    │◄──►│    (Go)         │
└─────────────────┘    └─────────────────┘
        │                       │
        │              ┌─────────────────┐
        │              │   Application   │
        │              │     Layer       │
        │              └─────────────────┘
        │                       │
        │              ┌─────────────────┐
        │              │   Service       │
        │              │     Layer       │
        │              └─────────────────┘
        │                       │
        │              ┌─────────────────┐
        │              │ Repository      │
        │              │     Layer       │
        │              └─────────────────┘
        │                       │
        │              ┌─────────────────┐
        └──────────────┤   Database      │
                       │   (SQLite)      │
                       └─────────────────┘
```

### Layer Responsibilities

**Frontend Layer (React/TypeScript):**
- User interface components
- State management
- API communication via Wails bindings
- Event handling and real-time updates

**Application Layer (internal/app):**
- Wails context binding
- Request/response handling
- Error management
- Event emission

**Service Layer:**
- **Parser Service** (`internal/parser`): OpenAPI specification parsing
- **Validator Service** (`internal/validator`): Spec validation logic
- **Generator Service** (`internal/generator`): MCP server generation
- **Mapping Service** (`internal/mapping`): OpenAPI to MCP mapping

**Repository Layer (internal/database):**
- Data access objects
- Database schema management
- Query optimization

### Key Components

**App Struct** (`internal/app/app.go`):
```go
type App struct {
    ctx                 context.Context
    db                  *sql.DB
    projectRepo         *database.ProjectRepository
    validationCacheRepo *database.ValidationCacheRepository
    parserService       *parser.Service
    mappingService      *mapping.Service
    generatorService    *generator.Service
    validatorService    *validator.Service
    settings            *AppSettings
    errorManager        *ErrorManager
    performanceMonitor  *PerformanceMonitor
}
```

## Code Organization

### Project Structure

```
MCPWeaver/
├── README.md                    # Project overview
├── LICENSE                     # AGPL v3 license
├── wails.json                  # Wails configuration
├── go.mod & go.sum            # Go module files
├── main.go                    # Application entry point
├── build/                     # Build outputs
├── internal/                  # Go backend code
│   ├── app/                   # Application layer
│   │   ├── app.go            # Main app struct and lifecycle
│   │   ├── files.go          # File operations
│   │   ├── types.go          # Type definitions
│   │   ├── errors.go         # Error handling
│   │   ├── settings.go       # Settings management
│   │   └── monitoring.go     # Performance monitoring
│   ├── database/             # Data layer
│   │   ├── connection.go     # Database connection
│   │   ├── migrations.go     # Schema migrations
│   │   ├── project_repo.go   # Project repository
│   │   └── validation_repo.go # Validation cache
│   ├── generator/            # Code generation
│   │   ├── service.go        # Generation service
│   │   ├── templates.go      # Template management
│   │   └── mcp_generator.go  # MCP server generator
│   ├── mapping/              # OpenAPI to MCP mapping
│   │   ├── service.go        # Mapping service
│   │   ├── endpoints.go      # Endpoint mapping
│   │   └── types.go          # Type mapping
│   ├── parser/               # OpenAPI parsing
│   │   ├── service.go        # Parser service
│   │   ├── openapi.go        # OpenAPI parser
│   │   └── validator.go      # Basic validation
│   └── validator/            # Advanced validation
│       ├── service.go        # Validation service
│       ├── rules.go          # Validation rules
│       └── suggestions.go    # Improvement suggestions
├── frontend/                 # React frontend
│   ├── package.json          # Node.js dependencies
│   ├── tsconfig.json         # TypeScript configuration
│   ├── public/               # Static assets
│   ├── src/                  # Source code
│   │   ├── components/       # React components
│   │   ├── services/         # API services
│   │   ├── types/            # TypeScript types
│   │   ├── hooks/            # Custom React hooks
│   │   ├── utils/            # Utility functions
│   │   └── styles/           # CSS/styling
│   └── tests/                # Frontend tests
├── docs/                     # Documentation
├── specs/                    # Technical specifications
└── scripts/                  # Build and utility scripts
```

### Naming Conventions

**Go Code:**
- **Packages**: Lowercase, single word (e.g., `parser`, `validator`)
- **Files**: Lowercase with underscores (e.g., `project_repo.go`)
- **Types**: PascalCase (e.g., `ProjectRepository`)
- **Functions/Methods**: PascalCase for exported, camelCase for private
- **Variables**: camelCase (e.g., `projectID`, `validationResult`)
- **Constants**: PascalCase or SCREAMING_SNAKE_CASE for package-level

**TypeScript/React:**
- **Components**: PascalCase (e.g., `ProjectManager.tsx`)
- **Files**: PascalCase for components, camelCase for utilities
- **Interfaces**: PascalCase with descriptive names (e.g., `ProjectData`)
- **Functions**: camelCase (e.g., `createProject`)
- **Variables**: camelCase (e.g., `projectList`)
- **Constants**: SCREAMING_SNAKE_CASE (e.g., `API_BASE_URL`)

## Development Workflow

### Branch Strategy

We use **Git Flow** with these branches:

- **main**: Production-ready code
- **develop**: Integration branch for features
- **feature/***: New features (e.g., `feature/project-templates`)
- **hotfix/***: Critical fixes (e.g., `hotfix/validation-crash`)
- **release/***: Release preparation (e.g., `release/v1.2.0`)

### Typical Development Flow

1. **Create Feature Branch**
   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b feature/your-feature-name
   ```

2. **Development**
   ```bash
   # Make changes
   # Run tests frequently
   wails dev  # For live development
   
   # Test your changes
   go test ./...
   npm test --prefix frontend
   ```

3. **Pre-commit Checks**
   ```bash
   # Format code
   go fmt ./...
   npm run format --prefix frontend
   
   # Lint code
   golangci-lint run
   npm run lint --prefix frontend
   
   # Run tests
   go test ./...
   npm test --prefix frontend
   
   # Build to ensure no issues
   wails build
   ```

4. **Commit and Push**
   ```bash
   git add .
   git commit -m "feat: add project template support"
   git push origin feature/your-feature-name
   ```

5. **Create Pull Request**
   - Use the GitHub PR template
   - Include description of changes
   - Link to related issues
   - Add appropriate labels
   - Request review from maintainers

### Commit Message Format

We follow **Conventional Commits** specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or modifying tests
- `chore`: Maintenance tasks

**Examples:**
```bash
feat(generator): add custom template support
fix(validation): handle empty OpenAPI specs
docs(api): update method documentation
test(parser): add edge case tests for malformed specs
```

## Testing

### Testing Strategy

MCPWeaver uses a multi-layered testing approach:

1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test component interactions
3. **End-to-End Tests**: Test complete workflows
4. **Performance Tests**: Verify performance requirements

### Backend Testing (Go)

**Test Organization:**
```
internal/
├── app/
│   ├── app_test.go
│   ├── files_test.go
│   └── testdata/
├── database/
│   ├── project_repo_test.go
│   └── migrations_test.go
└── parser/
    ├── service_test.go
    └── testdata/
        ├── valid_spec.json
        └── invalid_spec.yaml
```

**Writing Tests:**
```go
// internal/parser/service_test.go
package parser

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestParseOpenAPISpec(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected *OpenAPISpec
        wantErr  bool
    }{
        {
            name:  "valid OpenAPI 3.0 spec",
            input: `{"openapi":"3.0.0","info":{"title":"Test","version":"1.0.0"},"paths":{}}`,
            expected: &OpenAPISpec{
                OpenAPI: "3.0.0",
                Info: Info{
                    Title:   "Test",
                    Version: "1.0.0",
                },
                Paths: map[string]PathItem{},
            },
            wantErr: false,
        },
        {
            name:    "invalid JSON",
            input:   `{"invalid": json}`,
            wantErr: true,
        },
    }

    service := NewService()

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := service.ParseOpenAPISpec(tt.input)
            
            if tt.wantErr {
                require.Error(t, err)
                return
            }
            
            require.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func TestParseOpenAPISpecFromFile(t *testing.T) {
    service := NewService()
    
    // Test with valid file
    result, err := service.ParseOpenAPISpecFromFile("testdata/valid_spec.json")
    require.NoError(t, err)
    assert.NotNil(t, result)
    
    // Test with non-existent file
    _, err = service.ParseOpenAPISpecFromFile("testdata/non_existent.json")
    require.Error(t, err)
}

// Benchmark tests
func BenchmarkParseOpenAPISpec(b *testing.B) {
    service := NewService()
    spec := `{"openapi":"3.0.0","info":{"title":"Test","version":"1.0.0"},"paths":{}}`
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = service.ParseOpenAPISpec(spec)
    }
}
```

**Running Tests:**
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with detailed coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -run TestParseOpenAPISpec ./internal/parser

# Run benchmarks
go test -bench=. ./internal/parser

# Run tests with race detection
go test -race ./...

# Verbose output
go test -v ./...
```

### Frontend Testing (React/TypeScript)

**Test Organization:**
```
frontend/src/
├── components/
│   ├── ProjectManager.tsx
│   ├── ProjectManager.test.tsx
│   └── __tests__/
├── hooks/
│   ├── useAPI.ts
│   └── useAPI.test.ts
├── services/
│   ├── apiService.ts
│   └── apiService.test.ts
└── utils/
    ├── validation.ts
    └── validation.test.ts
```

**Testing Tools:**
- **Jest**: Test runner and assertion library
- **React Testing Library**: Component testing utilities
- **MSW**: API mocking
- **Cypress**: End-to-end testing

**Component Testing:**
```typescript
// components/ProjectManager.test.tsx
import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { ProjectManager } from './ProjectManager';

// Mock the Wails API
const mockAPI = {
  GetProjects: jest.fn(),
  CreateProject: jest.fn(),
  DeleteProject: jest.fn(),
};

(global as any).window = {
  go: {
    app: {
      App: mockAPI
    }
  }
};

describe('ProjectManager', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders project list', async () => {
    const mockProjects = [
      {
        id: '1',
        name: 'Test Project',
        description: 'Test Description',
        status: 'active',
        createdAt: '2023-01-01T00:00:00Z',
        updatedAt: '2023-01-01T00:00:00Z'
      }
    ];

    mockAPI.GetProjects.mockResolvedValue(mockProjects);

    render(<ProjectManager />);

    await waitFor(() => {
      expect(screen.getByText('Test Project')).toBeInTheDocument();
    });
  });

  it('creates new project', async () => {
    const newProject = {
      id: '2',
      name: 'New Project',
      description: 'New Description',
      status: 'active',
      createdAt: '2023-01-01T00:00:00Z',
      updatedAt: '2023-01-01T00:00:00Z'
    };

    mockAPI.CreateProject.mockResolvedValue(newProject);
    mockAPI.GetProjects.mockResolvedValue([]);

    render(<ProjectManager />);

    const createButton = screen.getByText('Create Project');
    fireEvent.click(createButton);

    const nameInput = screen.getByPlaceholderText('Project name');
    fireEvent.change(nameInput, { target: { value: 'New Project' } });

    const submitButton = screen.getByText('Create');
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(mockAPI.CreateProject).toHaveBeenCalledWith({
        name: 'New Project',
        description: '',
        openAPISpec: '',
        template: 'default'
      });
    });
  });

  it('handles API errors gracefully', async () => {
    mockAPI.GetProjects.mockRejectedValue(new Error('Network error'));

    render(<ProjectManager />);

    await waitFor(() => {
      expect(screen.getByText(/error/i)).toBeInTheDocument();
    });
  });
});
```

**Running Frontend Tests:**
```bash
# Run all tests
npm test

# Run tests in watch mode
npm run test:watch

# Run tests with coverage
npm run test:coverage

# Run specific test file
npm test ProjectManager.test.tsx

# Run tests matching pattern
npm test -- --testNamePattern="creates new project"
```

### Integration Testing

**Database Integration Tests:**
```go
// internal/database/integration_test.go
func TestProjectRepositoryIntegration(t *testing.T) {
    // Create temporary database
    db := setupTestDB(t)
    defer teardownTestDB(t, db)
    
    repo := database.NewProjectRepository(db)
    
    // Test create project
    project := &database.Project{
        Name:        "Integration Test Project",
        Description: "Test Description",
        OpenAPISpec: `{"openapi":"3.0.0"}`,
        Template:    "standard",
    }
    
    err := repo.Create(project)
    require.NoError(t, err)
    assert.NotEmpty(t, project.ID)
    
    // Test get project
    retrieved, err := repo.GetByID(project.ID)
    require.NoError(t, err)
    assert.Equal(t, project.Name, retrieved.Name)
    
    // Test update project
    retrieved.Description = "Updated Description"
    err = repo.Update(retrieved)
    require.NoError(t, err)
    
    // Verify update
    updated, err := repo.GetByID(project.ID)
    require.NoError(t, err)
    assert.Equal(t, "Updated Description", updated.Description)
}

func setupTestDB(t *testing.T) *database.DB {
    db, err := database.Open(":memory:")
    require.NoError(t, err)
    return db
}

func teardownTestDB(t *testing.T, db *database.DB) {
    err := db.Close()
    require.NoError(t, err)
}
```

### End-to-End Testing

**Cypress E2E Tests:**
```typescript
// frontend/cypress/integration/project_workflow.spec.ts
describe('Project Workflow', () => {
  beforeEach(() => {
    cy.visit('/');
    cy.wait(1000); // Wait for app to load
  });

  it('completes full project workflow', () => {
    // Create new project
    cy.get('[data-testid=create-project-button]').click();
    cy.get('[data-testid=project-name-input]').type('E2E Test Project');
    cy.get('[data-testid=project-description-input]').type('End-to-end test project');
    
    // Import OpenAPI spec
    cy.get('[data-testid=import-spec-button]').click();
    cy.get('[data-testid=spec-url-input]').type('https://petstore3.swagger.io/api/v3/openapi.json');
    cy.get('[data-testid=import-url-button]').click();
    
    // Wait for import to complete
    cy.get('[data-testid=import-success]', { timeout: 10000 }).should('be.visible');
    
    // Validate spec
    cy.get('[data-testid=validate-spec-button]').click();
    cy.get('[data-testid=validation-success]', { timeout: 5000 }).should('be.visible');
    
    // Generate server
    cy.get('[data-testid=generate-server-button]').click();
    
    // Monitor generation progress
    cy.get('[data-testid=generation-progress]', { timeout: 30000 }).should('contain', '100%');
    cy.get('[data-testid=generation-success]').should('be.visible');
    
    // Verify generated files
    cy.get('[data-testid=output-files-list]').should('contain', 'server.go');
    cy.get('[data-testid=output-files-list]').should('contain', 'handlers.go');
  });

  it('handles validation errors gracefully', () => {
    cy.get('[data-testid=create-project-button]').click();
    cy.get('[data-testid=project-name-input]').type('Invalid Spec Project');
    
    // Import invalid spec
    cy.get('[data-testid=import-spec-button]').click();
    cy.get('[data-testid=spec-content-textarea]').type('{"invalid": "spec"}');
    cy.get('[data-testid=import-content-button]').click();
    
    // Validate spec
    cy.get('[data-testid=validate-spec-button]').click();
    
    // Expect validation errors
    cy.get('[data-testid=validation-errors]').should('be.visible');
    cy.get('[data-testid=validation-suggestions]').should('be.visible');
  });
});
```

**Running E2E Tests:**
```bash
# Open Cypress UI
npm run cypress:open

# Run headless
npm run cypress:run

# Run specific test
npm run cypress:run -- --spec "cypress/integration/project_workflow.spec.ts"
```

## Building and Packaging

### Development Builds

```bash
# Development build (fast, includes debug info)
wails build -devMode

# Development with hot reload
wails dev

# Build for specific platform
wails build -platform windows/amd64
wails build -platform darwin/universal
wails build -platform linux/amd64
```

### Production Builds

```bash
# Production build (optimized)
wails build -clean -ldflags "-s -w"

# Cross-platform build
wails build -platform windows/amd64,darwin/universal,linux/amd64

# Build with version info
wails build -ldflags "-X main.version=1.0.0 -X main.commit=$(git rev-parse HEAD)"
```

### Build Configuration

**wails.json:**
```json
{
  "name": "MCPWeaver",
  "outputfilename": "MCPWeaver",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto",
  "author": {
    "name": "MCPWeaver Contributors",
    "email": "hello@mcpweaver.dev"
  },
  "info": {
    "companyName": "MCPWeaver",
    "productName": "MCPWeaver",
    "productVersion": "1.0.0",
    "copyright": "Copyright © 2024, MCPWeaver Contributors",
    "comments": "Transform OpenAPI specifications into MCP servers"
  },
  "nsisType": "multiple",
  "obfuscated": false,
  "garbleargs": "",
  "buildType": "platform"
}
```

### Advanced Build Features

**Code Signing (macOS):**
```bash
# Sign the application
codesign --deep --force --verify --verbose --sign "Developer ID Application: Your Name" MCPWeaver.app

# Verify signature
codesign --verify --verbose MCPWeaver.app
spctl --assess --verbose MCPWeaver.app
```

**Code Signing (Windows):**
```bash
# Sign with certificate
signtool sign /f certificate.p12 /p password /t http://timestamp.digicert.com MCPWeaver.exe

# Verify signature
signtool verify /pa MCPWeaver.exe
```

**Creating Installers:**
```bash
# Windows NSIS installer
makensis install.nsi

# macOS DMG
create-dmg --volname "MCPWeaver" --window-size 800 400 MCPWeaver.dmg MCPWeaver.app

# Linux AppImage
linuxdeploy --appdir AppDir --executable MCPWeaver --desktop-file mcpweaver.desktop --output appimage
```

## Contributing Guidelines

### Code Review Process

1. **Self Review**: Review your own code before submitting
2. **Automated Checks**: Ensure CI passes
3. **Peer Review**: At least one maintainer approval required
4. **Testing**: Verify tests pass and coverage meets requirements
5. **Documentation**: Update docs if needed

### Review Criteria

**Code Quality:**
- Follows established patterns and conventions
- Proper error handling
- Clear variable and function names
- Appropriate comments and documentation

**Testing:**
- Unit tests for new functionality
- Integration tests for component interactions
- E2E tests for user workflows
- Performance tests for critical paths

**Security:**
- Input validation
- SQL injection prevention
- XSS protection
- Secure file handling

**Performance:**
- Efficient algorithms
- Memory management
- Database query optimization
- UI responsiveness

### Issue Management

**Bug Reports:**
Use the bug report template and include:
- Environment details (OS, version, etc.)
- Steps to reproduce
- Expected vs actual behavior
- Error messages and logs
- Screenshots/videos if applicable

**Feature Requests:**
Use the feature request template and include:
- Use case description
- Proposed solution
- Alternative solutions considered
- Implementation complexity estimate

**Priority Labels:**
- `priority/critical`: Security issues, data loss, crashes
- `priority/high`: Major functionality issues
- `priority/medium`: Feature requests, minor bugs
- `priority/low`: Documentation, code cleanup

## Code Style and Standards

### Go Code Style

Follow standard Go conventions:

```go
// Good: Clear function documentation
// CreateProject creates a new project with the given request data.
// It validates the request, generates a unique ID, and stores the project
// in the database. Returns the created project or an error.
func (a *App) CreateProject(request CreateProjectRequest) (*Project, error) {
    // Validate input
    if request.Name == "" {
        return nil, a.createAPIError(
            "validation_error",
            "INVALID_PROJECT_NAME", 
            "Project name cannot be empty",
            map[string]string{"field": "name"},
        )
    }
    
    // Create project with generated ID
    project := &Project{
        ID:          generateUUID(),
        Name:        request.Name,
        Description: request.Description,
        OpenAPISpec: request.OpenAPISpec,
        Template:    request.Template,
        Status:      "active",
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    // Store in database
    if err := a.projectRepo.Create(project); err != nil {
        return nil, a.createAPIError(
            "database_error",
            "PROJECT_CREATE_FAILED",
            "Failed to create project",
            map[string]string{"error": err.Error()},
        )
    }
    
    // Emit project created event
    runtime.EventsEmit(a.ctx, "project:created", project)
    
    return project, nil
}

// Good: Clear struct documentation and tags
type Project struct {
    ID              string    `json:"id" db:"id"`
    Name            string    `json:"name" db:"name"`
    Description     string    `json:"description" db:"description"`
    OpenAPISpec     string    `json:"openAPISpec" db:"openapi_spec"`
    OutputDirectory string    `json:"outputDirectory" db:"output_directory"`
    Template        string    `json:"template" db:"template"`
    Status          string    `json:"status" db:"status"`
    CreatedAt       time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
    LastGenerated   time.Time `json:"lastGenerated,omitempty" db:"last_generated"`
}
```

**Error Handling:**
```go
// Good: Consistent error handling with context
func (s *Service) ProcessFile(path string) error {
    if !fileExists(path) {
        return fmt.Errorf("file not found: %s", path)
    }
    
    content, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("failed to read file %s: %w", path, err)
    }
    
    if err := s.validate(content); err != nil {
        return fmt.Errorf("validation failed for %s: %w", path, err)
    }
    
    return nil
}

// Good: Custom error types
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Value   string `json:"value"`
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error in field '%s': %s", e.Field, e.Message)
}
```

### TypeScript/React Style

Follow modern React and TypeScript best practices:

```typescript
// Good: Clear interface definitions
interface ProjectManagerProps {
  initialProjects?: Project[];
  onProjectCreate?: (project: Project) => void;
  onProjectUpdate?: (project: Project) => void;
  onProjectDelete?: (projectId: string) => void;
}

// Good: Functional component with hooks
const ProjectManager: React.FC<ProjectManagerProps> = ({
  initialProjects = [],
  onProjectCreate,
  onProjectUpdate,
  onProjectDelete,
}) => {
  const [projects, setProjects] = useState<Project[]>(initialProjects);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const { call } = useAPI();
  
  // Load projects on mount
  useEffect(() => {
    loadProjects();
  }, []);
  
  const loadProjects = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await call(() => window.go.app.App.GetProjects());
      if (response.data) {
        setProjects(response.data);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load projects');
    } finally {
      setLoading(false);
    }
  };
  
  const handleCreateProject = async (projectData: CreateProjectRequest) => {
    try {
      const response = await call(() => window.go.app.App.CreateProject(projectData));
      if (response.data) {
        setProjects(prev => [...prev, response.data!]);
        onProjectCreate?.(response.data);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create project');
    }
  };
  
  if (loading) {
    return <LoadingSpinner />;
  }
  
  if (error) {
    return (
      <ErrorMessage 
        message={error} 
        onRetry={loadProjects}
      />
    );
  }
  
  return (
    <div className="project-manager">
      <ProjectList 
        projects={projects}
        onUpdate={handleUpdateProject}
        onDelete={handleDeleteProject}
      />
      <CreateProjectForm onSubmit={handleCreateProject} />
    </div>
  );
};

export default ProjectManager;
```

**Custom Hooks:**
```typescript
// Good: Reusable custom hook
interface UseAsyncState<T> {
  data: T | null;
  loading: boolean;
  error: string | null;
  execute: (...args: any[]) => Promise<void>;
  reset: () => void;
}

function useAsyncState<T>(
  asyncFunction: (...args: any[]) => Promise<T>
): UseAsyncState<T> {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const execute = useCallback(async (...args: any[]) => {
    setLoading(true);
    setError(null);
    
    try {
      const result = await asyncFunction(...args);
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  }, [asyncFunction]);
  
  const reset = useCallback(() => {
    setData(null);
    setError(null);
    setLoading(false);
  }, []);
  
  return { data, loading, error, execute, reset };
}
```

### CSS/Styling Guidelines

Use CSS modules and follow BEM methodology:

```css
/* ProjectManager.module.css */
.projectManager {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1rem;
}

.projectManager__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.projectManager__title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.projectManager__createButton {
  background-color: var(--color-primary);
  color: white;
  border: none;
  border-radius: 0.5rem;
  padding: 0.5rem 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.projectManager__createButton:hover {
  background-color: var(--color-primary-dark);
}

.projectManager__createButton:disabled {
  background-color: var(--color-gray-400);
  cursor: not-allowed;
}

.projectManager__list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1rem;
}

.projectManager__emptyState {
  text-align: center;
  padding: 2rem;
  color: var(--color-text-secondary);
}
```

## Debugging

### Backend Debugging

**VS Code Launch Configuration (.vscode/launch.json):**
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug MCPWeaver",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}",
      "args": [],
      "env": {
        "WAILS_DEV_MODE": "true",
        "MCPWEAVER_LOG_LEVEL": "debug"
      },
      "showLog": true
    }
  ]
}
```

**Logging:**
```go
// Use structured logging
import "log/slog"

func (a *App) CreateProject(request CreateProjectRequest) (*Project, error) {
    slog.Info("Creating project", 
        "name", request.Name,
        "template", request.Template,
    )
    
    project, err := a.doCreateProject(request)
    if err != nil {
        slog.Error("Failed to create project",
            "name", request.Name,
            "error", err,
        )
        return nil, err
    }
    
    slog.Info("Project created successfully",
        "id", project.ID,
        "name", project.Name,
    )
    
    return project, nil
}
```

**Debug Helpers:**
```go
// Add debug endpoints in development
//go:build debug
func (a *App) GetDebugInfo() map[string]interface{} {
    return map[string]interface{}{
        "goroutines": runtime.NumGoroutine(),
        "memory": getMemoryStats(),
        "projects": len(a.projectRepo.GetAll()),
        "settings": a.settings,
    }
}

func getMemoryStats() map[string]interface{} {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return map[string]interface{}{
        "alloc": m.Alloc / 1024 / 1024,        // MB
        "totalAlloc": m.TotalAlloc / 1024 / 1024, // MB
        "sys": m.Sys / 1024 / 1024,           // MB
        "numGC": m.NumGC,
    }
}
```

### Frontend Debugging

**React Developer Tools:**
- Install React Developer Tools browser extension
- Use Profiler to identify performance bottlenecks
- Inspect component state and props

**Browser DevTools:**
```javascript
// Add debug helpers to window
if (process.env.NODE_ENV === 'development') {
  window.mcpweaver = {
    // Access app state
    getAppState: () => store.getState(),
    
    // Test API calls
    testAPI: {
      createProject: (data) => window.go.app.App.CreateProject(data),
      getProjects: () => window.go.app.App.GetProjects(),
    },
    
    // Performance monitoring
    performance: {
      measureAPI: async (name, fn) => {
        const start = performance.now();
        const result = await fn();
        const end = performance.now();
        console.log(`${name} took ${end - start}ms`);
        return result;
      }
    }
  };
}
```

**Error Boundaries:**
```typescript
class ErrorBoundary extends React.Component<
  React.PropsWithChildren<{}>,
  { hasError: boolean; error?: Error }
> {
  constructor(props: React.PropsWithChildren<{}>) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Error caught by boundary:', error, errorInfo);
    
    // Report error in development
    if (process.env.NODE_ENV === 'development') {
      window.go?.app?.App?.ReportError({
        type: 'react_error',
        message: error.message,
        stack: error.stack,
        componentStack: errorInfo.componentStack,
        timestamp: new Date().toISOString()
      });
    }
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="error-boundary">
          <h2>Something went wrong</h2>
          <details>
            <summary>Error details</summary>
            <pre>{this.state.error?.stack}</pre>
          </details>
        </div>
      );
    }

    return this.props.children;
  }
}
```

## Performance Optimization

### Backend Performance

**Memory Management:**
```go
// Use object pools for frequent allocations
var projectPool = sync.Pool{
    New: func() interface{} {
        return &Project{}
    },
}

func (r *ProjectRepository) Create(project *Project) error {
    // Use pooled object
    p := projectPool.Get().(*Project)
    defer projectPool.Put(p)
    
    // Copy data to pooled object
    *p = *project
    
    // Process...
    return nil
}

// Optimize database queries
func (r *ProjectRepository) GetProjectsWithStats() ([]*ProjectWithStats, error) {
    // Use single query instead of N+1
    query := `
        SELECT p.*, 
               COUNT(g.id) as generation_count,
               MAX(g.created_at) as last_generation
        FROM projects p
        LEFT JOIN generations g ON p.id = g.project_id
        GROUP BY p.id
        ORDER BY p.updated_at DESC
    `
    
    // Implementation...
    return results, nil
}
```

**Caching:**
```go
// Implement LRU cache for validation results
type ValidationCache struct {
    cache *lru.Cache
    mutex sync.RWMutex
}

func NewValidationCache(size int) *ValidationCache {
    cache, _ := lru.New(size)
    return &ValidationCache{
        cache: cache,
    }
}

func (c *ValidationCache) Get(key string) (*ValidationResult, bool) {
    c.mutex.RLock()
    defer c.mutex.RUnlock()
    
    if value, ok := c.cache.Get(key); ok {
        return value.(*ValidationResult), true
    }
    return nil, false
}

func (c *ValidationCache) Set(key string, result *ValidationResult) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    c.cache.Add(key, result)
}
```

### Frontend Performance

**Component Optimization:**
```typescript
// Use React.memo for expensive components
const ProjectCard = React.memo<ProjectCardProps>(({ project, onUpdate, onDelete }) => {
  return (
    <div className="project-card">
      {/* Component content */}
    </div>
  );
}, (prevProps, nextProps) => {
  // Custom comparison function
  return (
    prevProps.project.id === nextProps.project.id &&
    prevProps.project.updatedAt === nextProps.project.updatedAt
  );
});

// Use useMemo for expensive calculations
const ProjectStats: React.FC<{ projects: Project[] }> = ({ projects }) => {
  const stats = useMemo(() => {
    return {
      total: projects.length,
      active: projects.filter(p => p.status === 'active').length,
      totalSize: projects.reduce((acc, p) => acc + p.openAPISpec.length, 0),
    };
  }, [projects]);

  return <div>{/* Render stats */}</div>;
};

// Use useCallback for event handlers
const ProjectManager: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);

  const handleProjectUpdate = useCallback((updatedProject: Project) => {
    setProjects(prev => 
      prev.map(p => p.id === updatedProject.id ? updatedProject : p)
    );
  }, []);

  const handleProjectDelete = useCallback((projectId: string) => {
    setProjects(prev => prev.filter(p => p.id !== projectId));
  }, []);

  return (
    <div>
      {projects.map(project => (
        <ProjectCard
          key={project.id}
          project={project}
          onUpdate={handleProjectUpdate}
          onDelete={handleProjectDelete}
        />
      ))}
    </div>
  );
};
```

**Lazy Loading:**
```typescript
// Lazy load heavy components
const ProjectEditor = React.lazy(() => import('./ProjectEditor'));
const ServerGenerator = React.lazy(() => import('./ServerGenerator'));

const App: React.FC = () => {
  return (
    <Router>
      <Suspense fallback={<LoadingSpinner />}>
        <Routes>
          <Route path="/projects/:id/edit" element={<ProjectEditor />} />
          <Route path="/projects/:id/generate" element={<ServerGenerator />} />
        </Routes>
      </Suspense>
    </Router>
  );
};

// Virtual scrolling for large lists
import { FixedSizeList } from 'react-window';

const ProjectList: React.FC<{ projects: Project[] }> = ({ projects }) => {
  const Row = ({ index, style }: { index: number; style: React.CSSProperties }) => (
    <div style={style}>
      <ProjectCard project={projects[index]} />
    </div>
  );

  return (
    <FixedSizeList
      height={600}
      itemCount={projects.length}
      itemSize={150}
    >
      {Row}
    </FixedSizeList>
  );
};
```

## Security Considerations

### Input Validation

**Backend Validation:**
```go
// Sanitize and validate all inputs
func (a *App) CreateProject(request CreateProjectRequest) (*Project, error) {
    // Validate required fields
    if strings.TrimSpace(request.Name) == "" {
        return nil, a.createValidationError("Project name is required", nil, nil)
    }
    
    // Sanitize inputs
    request.Name = sanitizeString(request.Name)
    request.Description = sanitizeString(request.Description)
    
    // Validate OpenAPI spec
    if request.OpenAPISpec != "" {
        if err := validateOpenAPISpec(request.OpenAPISpec); err != nil {
            return nil, a.createValidationError("Invalid OpenAPI specification", nil, nil)
        }
    }
    
    // Validate output directory path
    if request.OutputDirectory != "" {
        if !isValidPath(request.OutputDirectory) {
            return nil, a.createValidationError("Invalid output directory path", nil, nil)
        }
    }
    
    // Continue with creation...
}

func sanitizeString(input string) string {
    // Remove control characters
    input = strings.Map(func(r rune) rune {
        if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
            return -1
        }
        return r
    }, input)
    
    // Trim whitespace
    return strings.TrimSpace(input)
}

func isValidPath(path string) bool {
    // Check for path traversal attempts
    if strings.Contains(path, "..") {
        return false
    }
    
    // Check for absolute paths pointing outside allowed directories
    absPath, err := filepath.Abs(path)
    if err != nil {
        return false
    }
    
    // Ensure path is within allowed directories
    allowedPaths := []string{
        "/tmp",
        "/var/tmp",
        os.TempDir(),
        ".",
    }
    
    for _, allowed := range allowedPaths {
        allowedAbs, _ := filepath.Abs(allowed)
        if strings.HasPrefix(absPath, allowedAbs) {
            return true
        }
    }
    
    return false
}
```

**Frontend Validation:**
```typescript
// Input sanitization and validation
import DOMPurify from 'dompurify';

interface ValidationRule {
  required?: boolean;
  minLength?: number;
  maxLength?: number;
  pattern?: RegExp;
  custom?: (value: string) => string | null;
}

interface ValidationRules {
  [key: string]: ValidationRule;
}

class InputValidator {
  static validate(data: Record<string, any>, rules: ValidationRules): Record<string, string> {
    const errors: Record<string, string> = {};
    
    for (const [field, rule] of Object.entries(rules)) {
      const value = data[field];
      
      if (rule.required && (!value || value.toString().trim() === '')) {
        errors[field] = `${field} is required`;
        continue;
      }
      
      if (value) {
        const stringValue = value.toString();
        
        if (rule.minLength && stringValue.length < rule.minLength) {
          errors[field] = `${field} must be at least ${rule.minLength} characters`;
        }
        
        if (rule.maxLength && stringValue.length > rule.maxLength) {
          errors[field] = `${field} must not exceed ${rule.maxLength} characters`;
        }
        
        if (rule.pattern && !rule.pattern.test(stringValue)) {
          errors[field] = `${field} format is invalid`;
        }
        
        if (rule.custom) {
          const customError = rule.custom(stringValue);
          if (customError) {
            errors[field] = customError;
          }
        }
      }
    }
    
    return errors;
  }
  
  static sanitizeInput(input: string): string {
    // Remove HTML tags and scripts
    return DOMPurify.sanitize(input, { ALLOWED_TAGS: [] });
  }
  
  static validateProjectData(data: CreateProjectRequest): Record<string, string> {
    const rules: ValidationRules = {
      name: {
        required: true,
        minLength: 1,
        maxLength: 100,
        pattern: /^[a-zA-Z0-9\s\-_]+$/,
      },
      description: {
        maxLength: 500,
      },
      template: {
        pattern: /^[a-zA-Z0-9\-_]+$/,
      },
    };
    
    return this.validate(data, rules);
  }
}
```

### File Security

```go
// Secure file operations
func (a *App) ReadFile(path string) (string, error) {
    // Validate file path
    if !isValidPath(path) {
        return "", a.createAPIError(
            "security_error",
            "INVALID_FILE_PATH",
            "File path is not allowed",
            map[string]string{"path": path},
        )
    }
    
    // Check file size limits
    info, err := os.Stat(path)
    if err != nil {
        return "", a.createFileSystemError("Failed to stat file", path, "read")
    }
    
    const maxFileSize = 10 * 1024 * 1024 // 10MB
    if info.Size() > maxFileSize {
        return "", a.createAPIError(
            "file_too_large",
            "FILE_SIZE_LIMIT_EXCEEDED",
            "File size exceeds maximum allowed size",
            map[string]string{
                "size": fmt.Sprintf("%d", info.Size()),
                "maxSize": fmt.Sprintf("%d", maxFileSize),
            },
        )
    }
    
    // Read file with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    content, err := readFileWithContext(ctx, path)
    if err != nil {
        return "", a.createFileSystemError("Failed to read file", path, "read")
    }
    
    return content, nil
}

func readFileWithContext(ctx context.Context, path string) (string, error) {
    done := make(chan struct{})
    var content string
    var err error
    
    go func() {
        defer close(done)
        data, readErr := os.ReadFile(path)
        if readErr != nil {
            err = readErr
            return
        }
        content = string(data)
    }()
    
    select {
    case <-ctx.Done():
        return "", ctx.Err()
    case <-done:
        return content, err
    }
}
```

## Release Process

### Version Management

We use **Semantic Versioning** (SemVer):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Workflow

1. **Prepare Release Branch**
   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b release/v1.2.0
   ```

2. **Update Version Numbers**
   ```bash
   # Update version in multiple files
   ./scripts/update-version.sh 1.2.0
   
   # Files to update:
   # - wails.json (productVersion)
   # - package.json (version)
   # - main.go (version variable)
   # - docs/INSTALLATION.md (download links)
   ```

3. **Update Changelog**
   ```bash
   # Generate changelog
   conventional-changelog -p angular -i CHANGELOG.md -s
   
   # Edit and clean up the changelog
   ```

4. **Final Testing**
   ```bash
   # Run full test suite
   go test ./...
   npm test --prefix frontend
   
   # Build for all platforms
   wails build -platform windows/amd64,darwin/universal,linux/amd64
   
   # Test installers
   ./scripts/test-installers.sh
   ```

5. **Merge and Tag**
   ```bash
   # Merge to main
   git checkout main
   git merge release/v1.2.0
   
   # Create tag
   git tag -a v1.2.0 -m "Release version 1.2.0"
   git push origin main
   git push origin v1.2.0
   ```

6. **Automated Release**
   The GitHub Actions workflow will:
   - Build for all platforms
   - Sign binaries
   - Create installers
   - Upload to GitHub Releases
   - Update documentation
   - Notify team

### Hotfix Process

For critical issues requiring immediate release:

1. **Create Hotfix Branch**
   ```bash
   git checkout main
   git checkout -b hotfix/v1.2.1
   ```

2. **Fix Issue**
   ```bash
   # Make minimal changes to fix the issue
   # Update version numbers
   # Add tests if possible
   ```

3. **Release**
   ```bash
   # Merge to main and develop
   git checkout main
   git merge hotfix/v1.2.1
   git checkout develop
   git merge hotfix/v1.2.1
   
   # Tag and push
   git tag -a v1.2.1 -m "Hotfix version 1.2.1"
   git push origin main develop
   git push origin v1.2.1
   ```

### Release Checklist

**Pre-release:**
- [ ] All tests pass
- [ ] Documentation updated
- [ ] Version numbers bumped
- [ ] Changelog updated
- [ ] Build artifacts tested
- [ ] Security scan passed

**Release:**
- [ ] Tag created
- [ ] GitHub release created
- [ ] Binaries uploaded
- [ ] Checksums verified
- [ ] Installation tested on all platforms

**Post-release:**
- [ ] Documentation site updated
- [ ] Team notified
- [ ] Social media announcement
- [ ] Update tracking metrics

---

## Getting Help

### Documentation Resources

- [User Guide](USER_GUIDE.md) - End-user documentation
- [API Reference](API.md) - Complete API documentation
- [Installation Guide](INSTALLATION.md) - Platform-specific installation
- [Troubleshooting](TROUBLESHOOTING.md) - Common issues and solutions

### Community Support

- **GitHub Issues**: [Report bugs and request features](https://github.com/matoval/MCPWeaver/issues)
- **GitHub Discussions**: [Community discussions](https://github.com/matoval/MCPWeaver/discussions)
- **Discord**: [Join our development Discord](https://discord.gg/mcpweaver)

### Contributing

We welcome all types of contributions:
- **Code**: Features, bug fixes, performance improvements
- **Documentation**: Guides, API docs, examples
- **Testing**: Bug reports, test cases, QA
- **Design**: UI/UX improvements, graphics, icons
- **Translation**: Internationalization support

Thank you for contributing to MCPWeaver! 🚀