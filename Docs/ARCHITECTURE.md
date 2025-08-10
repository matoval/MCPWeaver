# Technical Architecture

## Overview

MCPWeaver follows a pipeline-based architecture designed for simplicity, maintainability, and extensibility. The system transforms OpenAPI specifications into Model Context Protocol (MCP) servers through a series of well-defined stages.

## Architectural Principles

### Design Philosophy

- **Simplicity First**: Monolithic CLI for MVP with clear separation of concerns
- **Pipeline Pattern**: Sequential processing stages with fail-fast error handling
- **Pure Go**: No CGO dependencies, single binary distribution
- **Offline Operation**: Complete functionality without network connectivity
- **Convention over Configuration**: Minimal configuration with sensible defaults

### Core Components

The architecture consists of four primary components:

```text
┌──────────┐    ┌─────────────┐    ┌─────────────┐    ┌───────────┐
│  Parser  │───▶│ Transformer │───▶│  Generator  │───▶│   CLI     │
└──────────┘    └─────────────┘    └─────────────┘    └───────────┘
```

## Component Architecture

### 1. Parser Component

**Responsibility**: OpenAPI specification parsing and validation

#### Key Features

- Support for OpenAPI 2.0, 3.0+ specifications
- Fail-fast validation with line number reporting
- Normalization to internal data model
- Schema validation and type checking

#### Implementation Details

- **Library**: kin-openapi for robust OpenAPI parsing
- **Error Handling**: Detailed error context with file locations
- **Validation**: Comprehensive spec validation before transformation
- **Normalization**: Convert all OpenAPI versions to consistent internal format

```go
type Parser interface {
    Parse(filename string) (*OpenAPISpec, error)
    Validate(spec *OpenAPISpec) error
    Normalize(spec *OpenAPISpec) (*InternalSpec, error)
}
```

### 2. Transformer Component

**Responsibility**: Convert OpenAPI operations to MCP server components

#### Transformation Logic

- **Operations to Tools**: Each HTTP operation becomes an MCP tool
- **Parameter Mapping**: OpenAPI parameters to MCP tool parameters
- **Authentication**: Extract auth schemes for server configuration
- **Schema Processing**: Complex types to Python type hints

#### Internal Data Model

```go
type MCPServer struct {
    Name        string
    Tools       []MCPTool
    Resources   []MCPResource
    Prompts     []MCPPrompt
    Auth        AuthConfig
}

type MCPTool struct {
    Name        string
    Description string
    Parameters  []Parameter
    HTTPConfig  HTTPOperation
}
```

#### Processing Pipeline

1. **Operation Analysis**: Examine each OpenAPI path/method combination
2. **Tool Generation**: Create MCP tools with appropriate names and descriptions
3. **Parameter Extraction**: Map path, query, and body parameters
4. **Authentication Processing**: Handle API keys, basic auth, bearer tokens
5. **Validation Rules**: Generate input validation logic

### 3. Generator Component

**Responsibility**: Generate Python FastMCP server code and supporting files

#### Template System

- **Engine**: Go's `text/template` for code generation
- **Templates**: Embedded templates for distribution simplicity
- **Structure**: Single FastMCP template with modular components

#### Generated Files

```text
output-directory/
├── server.py          # Main FastMCP server
├── test_server.py     # Test suite with mocked responses  
├── README.md          # Usage documentation
├── requirements.txt   # Python dependencies
└── mcp-server         # Compiled executable (optional)
```

#### Code Generation Features

- **Async Functions**: FastMCP-compatible async tool implementations
- **HTTP Client**: Complete httpx-based HTTP client code
- **Error Handling**: Comprehensive try/catch with meaningful errors
- **Input Validation**: Parameter validation with type checking
- **Authentication**: Placeholder code for API key/token handling

### 4. CLI Component

**Responsibility**: User interface and command orchestration

#### Command Structure

```bash
mcpweaver generate <spec.yaml> [--output <dir>] [--verbose]
mcpweaver validate <spec.yaml>
mcpweaver version
```

#### Interactive Features

- **Endpoint Selection**: User interface for choosing operations to include
- **Progress Indication**: Processing feedback for long operations
- **Error Reporting**: Clear error messages with debugging context
- **Output Summary**: Report of generated files and skipped operations

## Data Flow

### Pipeline Stages

1. **Input Validation**
   - Verify OpenAPI file exists and is readable
   - Basic format validation (YAML/JSON)

2. **Parsing Stage**
   - Parse OpenAPI specification using kin-openapi
   - Validate against OpenAPI schema
   - Report parsing errors with line numbers

3. **Transformation Stage**
   - Normalize OpenAPI spec to internal model
   - Generate MCP tool definitions
   - Extract authentication configuration
   - Process schema definitions

4. **User Interaction**
   - Present available endpoints for selection
   - Allow filtering of operations to include
   - Confirm output directory and options

5. **Generation Stage**
   - Generate Python FastMCP server code
   - Create test suite with mocked data
   - Generate documentation and setup files
   - Optionally create executable binary

6. **Validation & Output**
   - Validate generated Python code syntax
   - Run generated tests to verify functionality
   - Report generation summary to user

### Error Handling Strategy

#### Error Categories

- **Parse Errors**: Invalid OpenAPI format, missing required fields
- **Validation Errors**: Schema violations, invalid references
- **Transformation Errors**: Unsupported features, complex schemas
- **Generation Errors**: Template rendering, file I/O issues

#### Error Context

- **Line Numbers**: Precise error locations in source files
- **File Paths**: Clear identification of problematic files
- **Suggestions**: Helpful guidance for common issues
- **Debug Mode**: Verbose output for troubleshooting

## Package Structure

```text
MCPWeaver/
├── main.go              # CLI entry point
├── cmd/                 # Command implementations
│   ├── generate.go
│   ├── validate.go
│   └── version.go
├── internal/            # Private application code
│   ├── parser/          # OpenAPI parsing
│   │   ├── service.go
│   │   └── types.go
│   ├── transformer/     # Spec transformation  
│   │   ├── service.go
│   │   └── mapping.go
│   ├── generator/       # Code generation
│   │   ├── service.go
│   │   └── templates/
│   └── common/          # Shared utilities
│       ├── errors.go
│       └── validation.go
├── templates/           # Embedded code templates
│   └── fastmcp/
│       ├── server.py.tmpl
│       ├── test.py.tmpl
│       └── readme.md.tmpl
└── testdata/           # Test fixtures
    ├── specs/
    └── expected/
```

## Technology Stack

### Core Dependencies

- **CLI Framework**: Cobra for command structure and help system
- **OpenAPI Parser**: kin-openapi for specification processing
- **Template Engine**: Go's text/template for code generation
- **HTTP Client**: Standard library for minimal dependencies

### Development Tools

- **Testing**: Standard Go testing framework with testify
- **Linting**: golangci-lint for code quality enforcement
- **Build**: Standard Go build tools with cross-compilation
- **Documentation**: godoc for API documentation

## Performance Considerations

### Design Decisions

- **Memory Management**: Load entire spec into memory for simplicity
- **Concurrency**: Single-threaded processing for MVP
- **Caching**: No caching to minimize complexity
- **Streaming**: Not implemented for MVP scope

### Performance Targets

- **Generation Time**: < 5 seconds for typical APIs (10-50 endpoints)
- **Memory Usage**: Minimal with no persistent state
- **Binary Size**: < 20MB single executable
- **Startup Time**: Instant CLI tool launch

## Extensibility Points

### Future Enhancements

- **Multiple Frameworks**: Additional MCP server frameworks
- **Custom Templates**: User-provided code templates
- **Output Formats**: TypeScript, Node.js, other languages
- **Authentication**: Advanced auth method support
- **Configuration**: YAML/JSON configuration files

### Plugin Architecture (Future)

- **Parser Plugins**: Support for API specification formats beyond OpenAPI
- **Generator Plugins**: Custom code generation backends
- **Transformer Plugins**: Custom transformation logic
- **Validator Plugins**: Additional validation rules

## Testing Strategy

### Test Categories

- **Unit Tests**: Individual component testing with mocks
- **Integration Tests**: End-to-end pipeline testing
- **Generated Code Tests**: Validation of generated server functionality
- **Edge Case Tests**: Malformed specs, complex schemas

### Test Data

- **Real-world APIs**: Public OpenAPI specifications from popular services
- **Edge Cases**: Complex schemas, authentication variations
- **Error Cases**: Invalid specs, missing fields, malformed data

### Quality Metrics

- **Coverage Target**: 80% minimum code coverage
- **Error Handling**: All error paths tested and validated
- **Performance**: Generation time benchmarks for various spec sizes
- **Compatibility**: Testing across OpenAPI version variations

This architecture provides a solid foundation for the MVP while maintaining flexibility for future enhancements and scaling.
