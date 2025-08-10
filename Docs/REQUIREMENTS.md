# Project Requirements

## Overview

MCPWeaver is a CLI tool that transforms OpenAPI specifications into Model Context Protocol (MCP) servers. This document defines the functional and non-functional requirements for the project.

## Functional Requirements

### Core Features

#### OpenAPI Support

- **Versions**: Support all OpenAPI versions (2.0, 3.0+)
- **HTTP Methods**: Support GET, POST, PUT, DELETE, PATCH operations
- **Authentication**: Handle API keys, basic auth, and bearer tokens
- **Data Types**: Support all OpenAPI data types with proper validation
- **Schema Support**: Handle complex schemas with nested objects and arrays
- **Server Definitions**: Process base URLs and server configurations

#### MCP Server Generation

- **Framework**: Generate Python servers using FastMCP
- **Capabilities**: Support MCP tools, resources, and prompts
- **Validation**: Include comprehensive error handling and input validation
- **Testing**: Generate complete test suites with mocked data
- **Documentation**: Create usage documentation and setup guides

#### Input/Output

- **Input**: Accept local OpenAPI specification files
- **Output**: Generate complete server directory structure
- **File Types**: Create Python server, tests, README, and executable binary
- **Customization**: Interactive endpoint selection during generation

### User Experience

#### CLI Interface

- **Commands**: Simple subcommand structure (`generate`, `validate`, `version`)
- **Options**: Minimal required parameters with sensible defaults
- **Feedback**: Processing indicators and clear success/error messages
- **Help**: Comprehensive help system with usage examples

#### Interactive Features

- **Endpoint Selection**: Allow users to choose which endpoints to include
- **Validation Feedback**: Display line numbers for errors and warnings
- **Output Summary**: Show generated files and skipped endpoints

### Quality Requirements

#### Error Handling

- **Validation**: Strict OpenAPI validation with fail-fast approach
- **Error Context**: Provide line numbers and file paths for debugging
- **Categorization**: Clear distinction between warnings and fatal errors
- **Recovery**: Skip invalid endpoints with user notification

#### Performance

- **Processing Time**: Generate servers in under 5 seconds for typical APIs
- **Memory Usage**: Efficient processing without persistent state
- **Startup**: Instant CLI tool startup time

## Non-Functional Requirements

### Technical Constraints

#### Implementation

- **Language**: Pure Go 1.21+ with no CGO dependencies
- **Architecture**: Pipeline-based (parse → transform → generate)
- **Dependencies**: Minimal external dependencies, embedded templates
- **Offline Operation**: Complete functionality without network connectivity

#### Platform Support

- **Operating Systems**: Linux and macOS for MVP
- **Architecture**: 64-bit systems (x86_64, ARM64)
- **Distribution**: Single binary with no external runtime dependencies

### Development Standards

#### Code Quality

- **Test Coverage**: Minimum 80% code coverage
- **Error Handling**: All errors must be handled with appropriate context
- **Logging**: Structured logging with debug mode support
- **Documentation**: Comprehensive inline code documentation

#### Testing Strategy

- **Unit Tests**: Individual component testing with mock data
- **Integration Tests**: End-to-end pipeline testing with real OpenAPI specs
- **Generated Code**: Validate that generated MCP servers work correctly
- **Edge Cases**: Handle malformed and complex OpenAPI specifications

## Scope Limitations

### MVP Exclusions

- URL-based OpenAPI spec fetching
- Multiple output formats (TypeScript, Node.js)
- Configuration file support
- Custom template systems
- OpenAPI callbacks and webhooks
- Deprecated endpoint processing
- Multiple specification processing in single run

### Future Considerations

- Additional MCP frameworks
- Custom authentication methods
- Advanced template customization
- Real-time specification watching
- Web-based interface
- Package manager integrations

## Success Criteria

### Minimum Viable Product

A CLI tool that can:

1. Parse any valid OpenAPI specification
2. Allow interactive endpoint selection
3. Generate a working Python FastMCP server
4. Include comprehensive tests with mocked data
5. Provide clear documentation for setup and usage
6. Handle errors gracefully with debugging information

### Quality Metrics

- **Reliability**: Handle 95% of real-world OpenAPI specifications correctly
- **Performance**: Process typical APIs (10-50 endpoints) in under 5 seconds
- **Usability**: New users can generate working servers in under 2 minutes
- **Maintainability**: Clear architecture supporting future enhancements

### Validation Criteria

- Generated servers pass all included tests
- Error messages provide actionable debugging information
- Documentation enables independent user success
- Binary distribution works across target platforms

## Dependencies

### Required Libraries

- **OpenAPI Parsing**: kin-openapi or go-swagger (to be determined)
- **CLI Framework**: Cobra for command structure and user interaction
- **Templates**: Go's text/template for code generation
- **HTTP Client**: Standard library for minimal dependencies

### Development Tools

- **Testing**: Standard Go testing framework
- **Build**: Standard Go build tools
- **Linting**: golangci-lint for code quality
- **Documentation**: godoc for API documentation

This requirements document serves as the foundation for development decisions and project scope definition.
