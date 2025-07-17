# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

MCPWeaver is an open-source desktop application that transforms OpenAPI specifications into Model Context Protocol (MCP) servers. This is a **specifications-only repository** - the actual implementation code will be in a separate repository.

## Repository Structure

This repository contains comprehensive specification documents that define the complete system architecture, API contracts, and implementation requirements:

```
MCPWeaver/
├── specs/
│   ├── PROJECT-SPECIFICATION.md      # Overall project vision and requirements
│   ├── ARCHITECTURE-SPECIFICATION.md # Technical architecture (Wails-based)
│   ├── API-SPECIFICATION.md          # Internal API contracts
│   ├── DATA-MODELS-SPECIFICATION.md  # Database schemas and data structures
│   ├── OBSERVABILITY-SPECIFICATION.md # Monitoring and logging
│   └── UI-SPECIFICATION.md           # User interface design
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

**Note**: This is a specifications repository. The actual implementation will use:

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

All errors follow consistent format:
```go
type APIError struct {
    Type      string            `json:"type"`
    Code      string            `json:"code"`
    Message   string            `json:"message"`
    Details   map[string]string `json:"details,omitempty"`
    Timestamp time.Time         `json:"timestamp"`
}
```

## Security Requirements

- **Input Validation**: Strict OpenAPI spec validation
- **File System Security**: Restricted file access patterns
- **Template Security**: Sandboxed template execution
- **Local Storage**: Encrypted sensitive data

## Implementation Phases

1. **Phase 1**: Foundation (Weeks 1-2) - Core architecture
2. **Phase 2**: Core MVP (Weeks 3-4) - Basic functionality
3. **Phase 3**: Polish (Weeks 5-6) - Error handling, optimization
4. **Phase 4**: Release (Weeks 7-8) - Packaging, documentation

## Testing Requirements

- **Unit Tests**: >90% code coverage
- **Integration Tests**: Full workflow testing
- **Performance Tests**: Meet all performance requirements
- **Cross-Platform Tests**: Windows, macOS, Linux validation

## Contributing Guidelines

When implementing MCPWeaver:

1. **Follow specifications exactly** - all implementation details are defined
2. **Implement incrementally** - follow the development phases
3. **Test thoroughly** - meet coverage and performance requirements
4. **Document extensively** - maintain specification alignment
5. **Validate continuously** - ensure spec compliance throughout development

## Related Projects

This project reuses core functionality from the existing **openapi2mcp** project, adapting it for desktop use with Wails framework.

**Local Development**: The openapi2mcp repository is located at `../openapi2mcp` relative to this repository on the local machine. This contains the core parsing and generation logic that will be extracted and adapted for the desktop application.