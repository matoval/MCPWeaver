# Command Line Interface Design

## Overview

MCPWeaver provides a clean, intuitive command-line interface that follows established CLI conventions while prioritizing ease of use and clear feedback. The design emphasizes simplicity with powerful functionality accessible through well-structured commands.

## Command Structure

### Primary Commands

The CLI uses a subcommand architecture for clear organization:

```bash
mcpweaver <command> [arguments] [flags]
```

#### Core Commands

##### Generate Command

```bash
mcpweaver generate <openapi-spec> [--output <directory>] [--verbose]
```

- **Purpose**: Convert OpenAPI specification to MCP server
- **Primary Input**: OpenAPI specification file (required positional argument)
- **Output Control**: `--output` flag specifies target directory (default: current directory)
- **Debugging**: `--verbose` flag enables detailed processing information

##### Validate Command

```bash
mcpweaver validate <openapi-spec>
```

- **Purpose**: Validate OpenAPI specification without generation
- **Features**: Comprehensive validation with line-number error reporting
- **Exit Codes**: 0 for valid, 2 for validation errors

##### Version Command

```bash
mcpweaver version
```

- **Purpose**: Display version information and build details
- **Output**: Version number, build date, Go version, commit hash

### Help System

#### General Help

```bash
mcpweaver --help
mcpweaver -h
```

#### Command-Specific Help

```bash
mcpweaver generate --help
mcpweaver validate --help
```

Each help output includes:

- Command description and purpose
- Usage patterns with examples
- Flag descriptions and default values
- Common use cases and tips

## User Experience Design

### Interactive Features

#### Endpoint Selection Interface

During generation, users are presented with an interactive selection interface:

```bash
Found 12 endpoints in specification:

 ✓ GET /users - List all users
 ✓ POST /users - Create new user  
 ✓ GET /users/{id} - Get user by ID
 ✗ PUT /users/{id} - Update user (deprecated - skipped)
 ✓ DELETE /users/{id} - Delete user
 ✓ GET /products - List products
 ✓ POST /products - Create product
 ...

Select endpoints to include:
[a] Select all  [n] Select none  [Enter] Continue with selection
```

#### Processing Feedback

Clear progress indication for long operations:

```bash
Processing OpenAPI specification...
✓ Parsing specification (spec.yaml)
✓ Validating OpenAPI format  
✓ Analyzing 12 endpoints
✓ Generating MCP server code
✓ Creating test suite
✓ Writing documentation

Generation complete! Files created:
  server.py - FastMCP server (8 endpoints)
  test_server.py - Test suite with mocked responses
  README.md - Setup and usage documentation
  requirements.txt - Python dependencies
  
Skipped 4 endpoints (deprecated or invalid)
```

### Error Handling and Reporting

#### Error Message Format

```bash
Error: Invalid OpenAPI specification
  File: api.yaml
  Line: 23
  Issue: Missing required 'responses' field for GET /users
  
Suggestion: Add a responses section to define expected return values
```

#### Warning System

```bash
Warning: Endpoint GET /deprecated marked as deprecated - skipping
Warning: Complex schema 'RecursiveNode' simplified to 'dict' type
Warning: Multiple auth schemes found - using first (ApiKeyAuth)
```

#### Exit Codes

- `0`: Success
- `1`: General error (file not found, permission denied)
- `2`: OpenAPI validation error
- `3`: Generation error (template issues, file I/O problems)

### Output Formatting

#### Standard Mode (Default)

- Minimal output focused on essential information
- Progress indicators for long operations
- Clear success/error messages
- Summary of generated files

#### Verbose Mode (`--verbose`)

- Detailed processing information
- Timing information for each stage
- Detailed validation results
- Template rendering details
- Debug information for troubleshooting

## Flag Design

### Global Flags

- `--help, -h`: Show help information
- `--verbose, -v`: Enable verbose output
- `--version`: Show version information

### Command-Specific Flags

#### Generate Command Flags

- `--output, -o <directory>`: Output directory for generated files
- `--force, -f`: Overwrite existing files without confirmation (future)
- `--dry-run`: Show what would be generated without creating files (future)

#### Validate Command Flags

- `--strict`: Enable strict validation mode (future)
- `--format <json|text>`: Output format for validation results (future)

## Usage Examples

### Basic Usage Scenarios

#### Simple Generation

```bash
# Generate MCP server in current directory
mcpweaver generate api.yaml

# Generate in specific directory
mcpweaver generate api.yaml --output ./my-server
```

#### Validation Workflow

```bash
# Validate specification first
mcpweaver validate api.yaml

# If valid, proceed with generation
mcpweaver generate api.yaml --output ./server
```

#### Development and Debugging

```bash
# Verbose mode for troubleshooting
mcpweaver generate api.yaml --verbose

# Check tool version
mcpweaver version
```

### Advanced Usage Patterns

#### CI/CD Integration

```bash
# Validate in CI pipeline
mcpweaver validate openapi.yaml || exit 1

# Generate server for deployment
mcpweaver generate openapi.yaml --output ./build/server
```

#### Multiple Environment Workflow

```bash
# Development server
mcpweaver generate dev-api.yaml --output ./dev-server

# Production server  
mcpweaver generate prod-api.yaml --output ./prod-server
```

## Input Handling

### File Path Resolution

- **Absolute paths**: Supported (`/path/to/spec.yaml`)
- **Relative paths**: Resolved from current working directory (`./api.yaml`)
- **Current directory**: Simple filename (`api.yaml`)

### Supported File Formats

- **YAML**: `.yaml`, `.yml` extensions
- **JSON**: `.json` extension
- **Auto-detection**: Based on file content if extension is ambiguous

### Input Validation

- File existence and readability checks
- Basic format validation (valid YAML/JSON)
- OpenAPI specification validation
- Clear error messages for common issues

## Output Organization

### Generated Directory Structure

```bash
output-directory/
├── server.py          # Main FastMCP server
├── test_server.py     # Comprehensive test suite  
├── README.md          # Setup and usage guide
├── requirements.txt   # Python dependencies
└── .env.example       # Environment variables template
```

### File Naming Conventions

- **Consistent naming**: snake_case for Python files
- **Descriptive names**: Clear purpose from filename
- **No conflicts**: Avoid overwriting important files
- **Standard extensions**: Conventional file extensions

## Configuration Philosophy

### Convention Over Configuration

- Sensible defaults for all options
- Minimal required parameters
- Standard output formats and locations
- Predictable behavior across different environments

### No Configuration Files (MVP)

- All settings embedded in command-line interface
- Reduces complexity and setup requirements
- Enables immediate use without additional configuration
- Future enhancement opportunity for advanced users

## Accessibility and Usability

### User-Friendly Design

- **Clear command names**: Self-explanatory command purposes
- **Logical flag names**: Intuitive short and long flag options
- **Helpful error messages**: Actionable guidance for resolving issues
- **Progress feedback**: Visual indicators for long-running operations

### Development Experience

- **Fast feedback**: Quick validation and error reporting
- **Reproducible builds**: Consistent output across environments
- **Debug support**: Verbose mode for troubleshooting
- **Integration friendly**: Suitable for automation and CI/CD

## Future Enhancements

### Planned CLI Improvements

- **Configuration files**: YAML/JSON configuration support
- **Watch mode**: Automatic regeneration on spec changes
- **Multi-spec support**: Batch processing of multiple specifications
- **Custom templates**: User-provided code generation templates
- **Shell completion**: Bash/zsh/fish completion scripts

### Advanced Features

- **Interactive mode**: Full TUI for complex configurations
- **Plugin system**: Extension points for custom functionality
- **Output formats**: Additional target languages and frameworks
- **Integration hooks**: Git hooks, IDE plugins, build tool integration

This CLI design provides a solid foundation for the MVP while maintaining extensibility for future enhancements and user needs.
