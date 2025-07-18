# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

MCPWeaver is an open-source desktop application that transforms OpenAPI specifications into Model Context Protocol (MCP) servers. This repository contains the **active implementation** with both specifications and working code.

## Repository Structure

This repository contains both specification documents and the implemented codebase:

```
MCPWeaver/
├── specs/                            # Specification documents
│   ├── PROJECT-SPECIFICATION.md      # Overall project vision and requirements
│   ├── ARCHITECTURE-SPECIFICATION.md # Technical architecture (Wails-based)
│   ├── API-SPECIFICATION.md          # Internal API contracts
│   ├── DATA-MODELS-SPECIFICATION.md  # Database schemas and data structures
│   ├── OBSERVABILITY-SPECIFICATION.md # Monitoring and logging
│   └── UI-SPECIFICATION.md           # User interface design
├── internal/                         # Go backend implementation
│   ├── app/                          # Main application logic
│   │   ├── app.go                    # Application context and lifecycle
│   │   ├── files.go                  # File operations and I/O
│   │   ├── types.go                  # Type definitions
│   │   ├── errors.go                 # Error handling system
│   │   └── ...                       # Other app modules
│   ├── database/                     # Database layer
│   ├── generator/                    # Code generation engine
│   ├── mapping/                      # OpenAPI to MCP mapping
│   ├── parser/                       # OpenAPI parsing
│   └── validator/                    # Specification validation
├── frontend/                         # React/TypeScript frontend
│   ├── src/
│   │   ├── components/               # React components
│   │   ├── services/                 # API services
│   │   ├── types/                    # TypeScript types
│   │   └── ...                       # Other frontend modules
│   └── ...                           # Frontend config files
├── README.md                         # Project overview
└── LICENSE                          # AGPL v3 license
```

## Architecture Overview

MCPWeaver is designed as a **Wails v2 desktop application** with:

- **Backend**: Go 1.21+ (reusing openapi2mcp core functionality)
- **Frontend**: React/TypeScript with Wails runtime
- **Database**: SQLite for local project storage
- **Distribution**: Single binary cross-platform (Windows, macOS, Linux)

### Key Components

1. **Parser Service**: OpenAPI specification parsing and validation
2. **Generator Service**: MCP server code generation with progress tracking
3. **Validator Service**: Specification and generated code validation
4. **Project Manager**: Local project and history management
5. **Template Manager**: Custom generation templates

## Development Commands

The repository contains a fully functional Wails v2 application. Use these commands:

```bash
# Wails development
wails dev              # Development with hot reload
wails build            # Production build
wails build -platform windows/amd64,darwin/amd64,linux/amd64  # Cross-platform build

# Go development (for backend services)
go test ./...          # Run all tests
go build              # Build backend
go mod tidy           # Clean dependencies

# Frontend development
npm install           # Install dependencies
npm run build         # Build frontend
npm run test          # Run tests
```

## Key Design Principles

1. **Spec-Driven Development**: Implementation must follow detailed specifications
2. **Desktop-First**: No cloud dependencies, everything runs locally
3. **Lightweight**: Minimal resource usage and fast startup (<2 seconds)
4. **User-Controlled**: Users have full control over their data and generation
5. **Performance**: <5 seconds for typical OpenAPI spec generation

## Performance Requirements

- **Startup Time**: <2 seconds cold start
- **Memory Usage**: <50MB base footprint
- **Generation Speed**: 
  - Small specs (<10 endpoints): <1 second
  - Medium specs (10-100 endpoints): <3 seconds
  - Large specs (100+ endpoints): <10 seconds

## Data Models

### Core Entities
- **Project**: Contains OpenAPI spec, settings, and generation history
- **GenerationJob**: Tracks generation progress and results
- **Template**: Custom generation templates
- **ValidationResult**: Spec validation results with errors/warnings

### Database Schema
SQLite database with tables for projects, generations, templates, and settings. See `DATA-MODELS-SPECIFICATION.md` for complete schema.

## API Structure

The application uses **Wails context binding** for frontend-backend communication:

### Main API Methods
- `CreateProject(request)` - Create new project
- `GetProjects()` - List all projects
- `GenerateServer(projectId)` - Start MCP server generation
- `ValidateSpec(specPath)` - Validate OpenAPI specification
- `SelectFile(filters)` - Open file dialog

### Real-time Events
- `project:created` - Project creation events
- `generation:progress` - Real-time generation progress
- `generation:completed` - Generation completion
- `system:notification` - System notifications

## Error Handling

The application implements comprehensive error handling with multiple layers:

### Backend Error System
All errors use the `APIError` type with enhanced features:
```go
type APIError struct {
    Type          string            `json:"type"`
    Code          string            `json:"code"`
    Message       string            `json:"message"`
    Details       map[string]string `json:"details,omitempty"`
    Timestamp     time.Time         `json:"timestamp"`
    Suggestions   []string          `json:"suggestions,omitempty"`
    CorrelationID string            `json:"correlationId,omitempty"`
    Severity      ErrorSeverity     `json:"severity"`
    Recoverable   bool              `json:"recoverable"`
    RetryAfter    *time.Duration    `json:"retryAfter,omitempty"`
    Context       *ErrorContext     `json:"context,omitempty"`
}
```

### Frontend Error System
- **ErrorBoundary**: React component that catches JavaScript errors
- **Error Display**: User-friendly error messages with recovery options
- **Retry Logic**: Automatic retry with exponential backoff
- **Error Reporting**: Integrated GitHub issue creation for bugs

### Error Categories
- **Validation Errors**: User input validation with suggestions
- **File System Errors**: File operation failures with recovery guidance
- **Network Errors**: Connection issues with retry mechanisms
- **Generation Errors**: Code generation failures with context
- **Internal Errors**: System errors with debugging information

## Security Requirements

- **Input Validation**: Strict OpenAPI spec validation
- **File System Security**: Restricted file access patterns
- **Template Security**: Sandboxed template execution
- **Local Storage**: Encrypted sensitive data

## Implementation Status

The project has progressed through multiple phases:

### ✅ Completed Features
- **Core Architecture**: Wails v2 application structure
- **File Import/Export System**: OpenAPI spec import from files and URLs
- **Error Handling System**: Comprehensive error management with recovery
- **Type System**: Complete type definitions for all entities
- **Database Layer**: SQLite integration with repositories
- **Project Management**: Basic project CRUD operations
- **Validation System**: OpenAPI specification validation
- **Recent Files**: File history tracking
- **Progress Tracking**: Real-time operation progress

### 🚧 In Progress
- **Code Generation**: MCP server generation from OpenAPI specs
- **Frontend Components**: React UI components
- **Template System**: Custom generation templates

### 📋 Planned
- **Advanced UI Features**: Complete user interface
- **Performance Optimization**: Meeting performance requirements
- **Cross-Platform Testing**: Windows, macOS, Linux validation
- **Documentation**: User guides and API documentation

## Code Patterns and Best Practices

### Backend (Go)
1. **Error Handling**: Always use `createAPIError()` for consistent error responses
2. **Helper Functions**: Extract common logic into reusable helper functions (e.g., `fileExists()`, `ensureDir()`)
3. **Context Validation**: Check `a.ctx` for nil before Wails runtime calls
4. **Resource Cleanup**: Always clean up temporary files and close connections
5. **Structured Types**: Use comprehensive type definitions with JSON tags

### Frontend (React/TypeScript)
1. **Error Boundaries**: Wrap components with `ErrorBoundary` for error recovery
2. **Type Safety**: Use TypeScript interfaces that match backend types
3. **Error Reporting**: Integrate with backend error reporting system
4. **Progress Tracking**: Show progress for long-running operations
5. **Retry Logic**: Implement automatic retry with exponential backoff

### File Operations
1. **Validation**: Always validate file existence and permissions
2. **Helper Functions**: Use centralized helpers (`fileExists`, `dirExists`, `ensureDir`)
3. **Progress Tracking**: Emit progress events for file operations
4. **Error Context**: Include file paths and operation details in errors
5. **Cleanup**: Remove temporary files after operations

## Testing Requirements

- **Unit Tests**: >90% code coverage (especially for helper functions)
- **Integration Tests**: Full workflow testing
- **Performance Tests**: Meet all performance requirements
- **Cross-Platform Tests**: Windows, macOS, Linux validation
- **Error Handling Tests**: Comprehensive error scenario testing

## Contributing Guidelines

When working on MCPWeaver:

1. **Follow Existing Patterns**: Use established error handling and code organization patterns
2. **Extract Common Logic**: Create helper functions for repeated operations
3. **Test Thoroughly**: Include unit tests for new functionality
4. **Document Changes**: Update specifications and documentation
5. **Maintain Consistency**: Follow existing naming conventions and structure

## API Methods Reference

### File Operations (internal/app/files.go)

- `SelectFile(filters)` - Opens file selection dialog
- `SelectDirectory(title)` - Opens directory selection dialog  
- `SaveFile(content, defaultPath, filters)` - Opens save file dialog
- `ReadFile(path)` - Reads file content
- `WriteFile(path, content)` - Writes content to file
- `FileExists(path)` - Checks if file exists
- `ImportOpenAPISpec(filePath)` - Imports OpenAPI spec from file
- `ImportOpenAPISpecFromURL(url)` - Imports OpenAPI spec from URL
- `ExportGeneratedServer(projectID, targetDir)` - Exports generated server
- `AddRecentFile(filePath, fileType)` - Adds file to recent files
- `GetRecentFiles()` - Gets list of recent files
- `ClearRecentFiles()` - Clears recent files list

### Error Handling (internal/app/errors.go)

- `CreateError(errType, code, message, options...)` - Creates APIError
- `CreateValidationError(message, details, suggestions)` - Creates validation error
- `CreateNetworkError(message, details)` - Creates network error
- `CreateFileSystemError(message, filePath, operation)` - Creates file system error
- `CreateGenerationError(message, projectID, step)` - Creates generation error
- `CreateInternalError(message, err)` - Creates internal error
- `CreateErrorCollection(operation, totalItems)` - Creates error collection

### Helper Functions

- `fileExists(path)` - Internal file existence check
- `dirExists(path)` - Internal directory existence and writability check
- `ensureDir(path)` - Creates directory if it doesn't exist
- `convertFilters(filters)` - Converts FileFilter to Wails format

## Event System

### Emitted Events

- `system:startup` - Application startup
- `system:ready` - DOM ready
- `system:shutdown` - Application shutdown
- `system:error` - Error occurred
- `system:notification` - System notification
- `file:progress` - File operation progress
- `project:created` - Project creation
- `generation:progress` - Generation progress
- `generation:completed` - Generation completion

## Related Projects

This project reuses core functionality from the existing **openapi2mcp** project, adapting it for desktop use with Wails framework.

**Local Development**: The openapi2mcp repository is located at `../openapi2mcp` relative to this repository on the local machine. This contains the core parsing and generation logic that will be extracted and adapted for the desktop application.