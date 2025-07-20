# MCPWeaver User Guide

Welcome to MCPWeaver! This comprehensive guide will help you transform OpenAPI specifications into Model Context Protocol (MCP) servers quickly and efficiently.

## Table of Contents

- [Getting Started](#getting-started)
- [Interface Overview](#interface-overview)
- [Working with Projects](#working-with-projects)
- [OpenAPI Import and Validation](#openapi-import-and-validation)
- [Generating MCP Servers](#generating-mcp-servers)
- [Templates and Customization](#templates-and-customization)
- [Project Management](#project-management)
- [Settings and Preferences](#settings-and-preferences)
- [Best Practices](#best-practices)
- [Advanced Usage](#advanced-usage)

## Getting Started

### What is MCPWeaver?

MCPWeaver is a desktop application that converts OpenAPI 3.0 specifications into fully functional Model Context Protocol (MCP) servers. MCP servers enable large language models (LLMs) to interact with APIs in a standardized way.

### What You'll Need

- **OpenAPI Specification**: A valid OpenAPI 3.0 specification (JSON or YAML)
- **MCPWeaver**: Installed on your Windows, macOS, or Linux system
- **Basic API Knowledge**: Understanding of REST APIs and OpenAPI specifications

### Your First MCP Server

Let's create your first MCP server in under 5 minutes:

#### Step 1: Launch MCPWeaver

1. Open MCPWeaver from your applications menu or desktop shortcut
2. The application will start and show the main dashboard
3. You'll see options to create a new project or open existing ones

#### Step 2: Create a New Project

1. Click **"New Project"** or use `Ctrl+N` (`Cmd+N` on macOS)
2. Enter a project name (e.g., "My First MCP Server")
3. Optionally add a description
4. Click **"Create Project"**

#### Step 3: Import Your OpenAPI Specification

You can import your OpenAPI spec in several ways:

**From File:**
1. Click **"Import from File"** or drag and drop your OpenAPI file
2. Select your `.json`, `.yaml`, or `.yml` file
3. MCPWeaver will automatically validate the specification

**From URL:**
1. Click **"Import from URL"**
2. Enter the URL of your OpenAPI specification
3. Click **"Import"** to fetch and validate

**Example OpenAPI URLs to Try:**
- Petstore API: `https://petstore3.swagger.io/api/v3/openapi.json`
- JSONPlaceholder: `https://jsonplaceholder.typicode.com/openapi.json`

#### Step 4: Review and Validate

1. MCPWeaver will automatically validate your OpenAPI specification
2. Any errors or warnings will be highlighted in the **Validation Panel**
3. Review the detected endpoints, schemas, and operations
4. Fix any validation issues before proceeding

#### Step 5: Generate Your MCP Server

1. Click **"Generate MCP Server"** once validation passes
2. Choose your output directory
3. Select any custom templates (optional)
4. Click **"Start Generation"**
5. Watch the real-time progress as your server is created

#### Step 6: Test Your Generated Server

1. Navigate to your output directory
2. Follow the README instructions in the generated server
3. Run `go mod tidy && go run .` to start your MCP server
4. Your server is now ready to connect to LLMs!

## Interface Overview

### Main Dashboard

The main dashboard provides an overview of your projects and quick access to common actions:

- **Recent Projects**: Your most recently worked on projects
- **Quick Actions**: Create new project, import spec, or open existing project
- **Performance Stats**: Memory usage, generation times, and system status
- **News and Tips**: Updates and helpful tips for using MCPWeaver

### Project Workspace

When working on a project, you'll see several key areas:

#### 1. Project Explorer (Left Panel)
- **Project Information**: Name, description, creation date
- **OpenAPI Specification**: Current spec details and status
- **Generation History**: Previous generations with timestamps
- **Files and Assets**: Associated project files

#### 2. Main Editor (Center)
- **Specification Viewer**: Read-only view of your OpenAPI spec
- **Validation Results**: Real-time validation feedback
- **Generation Progress**: Live updates during server generation
- **Code Preview**: Preview generated code before saving

#### 3. Properties Panel (Right)
- **Generation Settings**: Output format, templates, customizations
- **Template Variables**: Configure template-specific variables
- **Advanced Options**: Caching, optimization, and debugging settings

#### 4. Status Bar (Bottom)
- **Validation Status**: Overall spec health
- **Generation Status**: Current operation progress
- **Performance Metrics**: Memory usage and processing time
- **Connection Status**: File system and external service status

### Toolbar and Menus

#### File Menu
- **New Project** (`Ctrl+N`): Create a new project
- **Open Project** (`Ctrl+O`): Open an existing project
- **Save Project** (`Ctrl+S`): Save current project
- **Import Specification**: Import OpenAPI spec from file or URL
- **Export Results**: Export generated servers or project data
- **Recent Projects**: Quick access to recent projects

#### Edit Menu
- **Undo/Redo** (`Ctrl+Z`/`Ctrl+Y`): Undo or redo recent actions
- **Copy** (`Ctrl+C`): Copy selected content
- **Select All** (`Ctrl+A`): Select all content in current view
- **Find** (`Ctrl+F`): Search within specifications or generated code

#### Project Menu
- **Validate Specification**: Manual validation trigger
- **Generate Server**: Start MCP server generation
- **Clear Cache**: Clear validation and generation cache
- **Project Settings**: Configure project-specific settings

#### Tools Menu
- **Settings**: Application preferences and configuration
- **Template Manager**: Manage and customize generation templates
- **Performance Monitor**: View detailed performance metrics
- **Logs**: Access application logs and debugging information

#### Help Menu
- **User Guide**: This documentation
- **API Reference**: Complete API documentation
- **Keyboard Shortcuts**: List of all shortcuts
- **About**: Application version and information

## Working with Projects

### Project Structure

MCPWeaver organizes your work into projects. Each project contains:

- **OpenAPI Specification**: The source API specification
- **Generation History**: Record of all server generations
- **Custom Settings**: Project-specific generation preferences
- **Cache Data**: Validation and processing cache for faster operations

### Creating Projects

#### New Project Wizard

1. **Basic Information**:
   - Project name (required)
   - Description (optional)
   - Output directory preference
   - Default template selection

2. **Import Specification**:
   - Choose import method (file, URL, or paste)
   - Automatic validation
   - Error correction suggestions

3. **Initial Setup**:
   - Template selection
   - Generation preferences
   - Variable configuration

#### Project Templates

MCPWeaver provides several project templates:

- **Basic REST API**: Standard OpenAPI to MCP conversion
- **Microservice Gateway**: Multiple service aggregation
- **Legacy System Integration**: Older API modernization
- **Development Testing**: Testing and mock server setup

### Managing Projects

#### Project Dashboard

The project dashboard shows:
- Generation statistics
- Recent activity
- Validation status
- Performance metrics
- Quick actions

#### Project Settings

Configure project-specific options:

**General Settings**:
- Project name and description
- Default output directory
- Auto-save preferences
- Backup settings

**Generation Settings**:
- Default template
- Output format preferences
- Validation strictness
- Error handling behavior

**Performance Settings**:
- Cache settings
- Memory limits
- Timeout configurations
- Parallel processing options

#### Project Export/Import

**Export Project**:
1. Go to **File > Export Project**
2. Choose export format (MCPWeaver project file or ZIP archive)
3. Select what to include (specs, history, settings, cache)
4. Choose destination and export

**Import Project**:
1. Go to **File > Import Project**
2. Select project file or archive
3. Choose import location
4. Review and confirm import settings

## OpenAPI Import and Validation

### Supported Formats

MCPWeaver supports OpenAPI specifications in multiple formats:

#### File Formats
- **JSON** (`.json`): Standard JSON format
- **YAML** (`.yaml`, `.yml`): Human-readable YAML format
- **Compressed** (`.zip`): ZIP archives containing specifications

#### OpenAPI Versions
- **OpenAPI 3.0.x**: Full support with all features
- **OpenAPI 3.1.x**: Full support with JSON Schema compatibility
- **Swagger 2.0**: Limited support with automatic conversion

### Import Methods

#### From Local Files

1. **Drag and Drop**:
   - Simply drag your OpenAPI file into MCPWeaver
   - Automatic format detection
   - Instant validation feedback

2. **File Browser**:
   - Click **"Import from File"**
   - Use the file browser to select your specification
   - Support for multiple file selection

3. **Recent Files**:
   - Access recently imported files from the File menu
   - Quick re-import with remembered settings
   - Automatic change detection

#### From URLs

1. **Direct URL Import**:
   - Click **"Import from URL"**
   - Enter the URL of your OpenAPI specification
   - Support for authentication headers

2. **Common API Sources**:
   - GitHub repositories
   - API documentation sites
   - Development servers
   - Cloud storage services

#### From Clipboard

1. **Paste Content**:
   - Copy OpenAPI content from any source
   - Use **Edit > Paste Specification** or `Ctrl+V`
   - Automatic format detection

### Validation Process

#### Real-time Validation

MCPWeaver continuously validates your OpenAPI specification:

- **Syntax Validation**: JSON/YAML syntax errors
- **Schema Validation**: OpenAPI schema compliance
- **Semantic Validation**: Logical consistency checks
- **Best Practice Validation**: Industry standard recommendations

#### Validation Results

The validation panel shows different types of issues:

**Errors** (Red 🔴):
- Must be fixed before generation
- Prevent successful MCP server creation
- Include specific location and fix suggestions

**Warnings** (Yellow 🟡):
- Should be addressed for best results
- Won't prevent generation but may affect functionality
- Include optimization suggestions

**Info** (Blue 🔵):
- Informational messages
- Best practice recommendations
- Performance optimization tips

**Success** (Green 🟢):
- Validation passed successfully
- Ready for MCP server generation
- Performance and quality metrics

#### Common Validation Issues

**Syntax Errors**:
```yaml
# ❌ Invalid YAML syntax
paths:
  /users
    get:  # Missing colon
      summary: Get users

# ✅ Correct YAML syntax
paths:
  /users:
    get:
      summary: Get users
```

**Missing Required Fields**:
```yaml
# ❌ Missing required fields
openapi: 3.0.0
# Missing 'info' field

# ✅ Complete minimal spec
openapi: 3.0.0
info:
  title: My API
  version: 1.0.0
paths: {}
```

**Invalid References**:
```yaml
# ❌ Invalid reference
components:
  schemas:
    User:
      $ref: '#/components/schemas/NonExistentSchema'

# ✅ Valid reference
components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: integer
```

### Auto-Correction Features

MCPWeaver can automatically fix common issues:

#### Auto-Fix Options

- **Format Cleanup**: Standardize indentation and spacing
- **Reference Resolution**: Fix broken internal references
- **Schema Enhancement**: Add missing required fields with defaults
- **Best Practice Application**: Apply OpenAPI best practices

#### Manual Correction

For complex issues, MCPWeaver provides:

- **Interactive Error Panel**: Click on errors to jump to problematic sections
- **Fix Suggestions**: Specific recommendations for each issue
- **Example Corrections**: Show correct syntax for common problems
- **Documentation Links**: Direct links to OpenAPI specification docs

## Generating MCP Servers

### Generation Overview

MCPWeaver transforms your validated OpenAPI specification into a complete MCP server with:

- **MCP Protocol Implementation**: Full MCP specification compliance
- **API Client Code**: Generated client for your OpenAPI endpoints
- **Error Handling**: Robust error handling and recovery
- **Documentation**: Generated README and API documentation
- **Testing Tools**: Basic testing utilities and examples

### Generation Process

#### Step 1: Pre-Generation Setup

1. **Validate Specification**: Ensure all validation issues are resolved
2. **Configure Output**: Choose destination directory
3. **Select Template**: Pick appropriate generation template
4. **Set Variables**: Configure template-specific variables

#### Step 2: Generation Options

**Basic Options**:
- **Output Directory**: Where to save the generated server
- **Server Name**: Name for the generated MCP server
- **Package Name**: Go package name for the server
- **Module Path**: Go module path for dependencies

**Advanced Options**:
- **Template Customization**: Modify generation templates
- **Feature Flags**: Enable/disable specific features
- **Optimization Level**: Balance between performance and readability
- **Documentation Level**: Control documentation detail

#### Step 3: Generation Execution

1. **Start Generation**: Click "Generate MCP Server"
2. **Monitor Progress**: Watch real-time progress in the progress panel
3. **Review Results**: Examine generated files and logs
4. **Test Output**: Validate the generated server works correctly

### Generation Templates

#### Built-in Templates

**Standard MCP Server**:
- Full MCP protocol implementation
- REST API client generation
- Error handling and recovery
- Basic documentation

**Lightweight Server**:
- Minimal MCP implementation
- Reduced dependencies
- Faster startup time
- Essential features only

**Development Server**:
- Extended debugging features
- Detailed logging
- Test utilities
- Mock data support

**Production Server**:
- Performance optimizations
- Security enhancements
- Monitoring integration
- Production-ready configuration

#### Custom Templates

Create your own generation templates:

1. **Template Structure**: Understand template file organization
2. **Variable System**: Use template variables for customization
3. **Template Language**: Learn the template syntax (Go templates)
4. **Testing Templates**: Validate custom templates before use

### Generated Server Structure

A typical generated MCP server includes:

```
my-mcp-server/
├── main.go                 # Server entry point
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── README.md               # Setup and usage instructions
├── client/                 # Generated API client
│   ├── client.go          # Main client implementation
│   ├── models.go          # Data models
│   └── operations.go      # API operations
├── mcp/                   # MCP protocol implementation
│   ├── server.go          # MCP server
│   ├── handlers.go        # Request handlers
│   └── tools.go           # MCP tools definition
├── config/                # Configuration
│   ├── config.go          # Configuration structure
│   └── defaults.go        # Default values
├── docs/                  # Documentation
│   ├── API.md             # API documentation
│   └── SETUP.md           # Setup instructions
└── tests/                 # Testing utilities
    ├── client_test.go     # Client tests
    └── integration_test.go # Integration tests
```

### Post-Generation Steps

#### 1. Review Generated Code

- **Examine the README**: Understanding setup and usage
- **Review Configuration**: Check default settings
- **Inspect API Client**: Verify endpoint mapping is correct
- **Test MCP Tools**: Ensure tools are properly defined

#### 2. Build and Test

```bash
# Navigate to generated server directory
cd my-mcp-server

# Install dependencies
go mod tidy

# Build the server
go build -o mcp-server

# Run tests
go test ./...

# Start the server
./mcp-server
```

#### 3. Integration

- **Connect to LLM**: Configure your LLM to use the MCP server
- **Test Endpoints**: Verify all API endpoints work correctly
- **Monitor Performance**: Check response times and error rates
- **Deploy**: Deploy to your target environment

## Templates and Customization

### Template System Overview

MCPWeaver uses a flexible template system based on Go templates to generate MCP servers. Templates define how OpenAPI specifications are transformed into working code.

### Understanding Templates

#### Template Components

**Template Files**:
- **Main Template**: Core server implementation
- **Client Template**: API client generation
- **Model Template**: Data structure generation
- **Documentation Template**: README and docs generation

**Template Variables**:
- **Specification Data**: OpenAPI spec content
- **Generation Context**: Metadata about the generation process
- **User Variables**: Custom values set during generation
- **System Variables**: MCPWeaver-specific data

#### Template Syntax

MCPWeaver uses Go template syntax with additional helper functions:

```go
// Basic variable substitution
{{.PackageName}}

// Conditional logic
{{if .HasAuthentication}}
// Authentication code
{{end}}

// Loops
{{range .Endpoints}}
func {{.Name}}() {
    // Endpoint implementation
}
{{end}}

// Helper functions
{{.MethodName | toLower}}
{{.Description | escapeQuotes}}
```

### Managing Templates

#### Template Manager

Access the Template Manager from **Tools > Template Manager**:

1. **View Templates**: See all available templates
2. **Edit Templates**: Modify existing templates
3. **Create Templates**: Design new custom templates
4. **Import/Export**: Share templates with others
5. **Validate Templates**: Test templates before use

#### Template Categories

**System Templates** (Read-only):
- Built-in templates provided by MCPWeaver
- Regularly updated with new features
- Cannot be modified directly
- Use as base for custom templates

**User Templates** (Editable):
- Created by you or imported from others
- Fully customizable
- Can override system templates
- Stored in user configuration directory

**Project Templates** (Project-specific):
- Associated with specific projects
- Override user and system templates
- Perfect for project-specific requirements
- Included in project exports

### Creating Custom Templates

#### Template Creation Wizard

1. **Start New Template**:
   - Click "Create New Template" in Template Manager
   - Choose base template to start from
   - Enter template name and description

2. **Define Variables**:
   - Set custom variables for your template
   - Define default values
   - Add validation rules
   - Include helpful descriptions

3. **Edit Template Files**:
   - Modify main server template
   - Customize client generation
   - Adjust model templates
   - Update documentation templates

4. **Test Template**:
   - Validate template syntax
   - Test with sample OpenAPI specs
   - Check generated output
   - Fix any issues

#### Advanced Template Features

**Conditional Generation**:
```go
{{if .HasWebhooks}}
// Webhook handling code
type WebhookHandler struct {
    // Implementation
}
{{end}}
```

**Dynamic Imports**:
```go
import (
    "context"
    "fmt"
    {{if .RequiresAuth}}"crypto/tls"{{end}}
    {{range .AdditionalImports}}
    "{{.}}"
    {{end}}
)
```

**Helper Functions**:
```go
// Convert to camelCase
{{.FieldName | camelCase}}

// Convert to snake_case
{{.FieldName | snakeCase}}

// Escape for JSON
{{.Description | jsonEscape}}

// Format as Go comment
{{.Documentation | goComment}}
```

### Template Variables Reference

#### Specification Variables

```go
// Basic specification info
{{.Spec.Info.Title}}
{{.Spec.Info.Version}}
{{.Spec.Info.Description}}

// Server information
{{range .Spec.Servers}}
{{.URL}} - {{.Description}}
{{end}}

// Paths and operations
{{range .Spec.Paths}}
Path: {{.Path}}
{{range .Operations}}
  Method: {{.Method}}
  OperationID: {{.OperationID}}
{{end}}
{{end}}
```

#### Generation Context

```go
// Generation metadata
{{.GenerationTime}}     // When generation started
{{.MCPWeaverVersion}}   // MCPWeaver version
{{.ProjectName}}        // Current project name
{{.OutputDirectory}}    // Target output directory

// User preferences
{{.UserName}}           // System username
{{.WorkspacePath}}      // Current workspace
{{.TemplateVersion}}    // Template version
```

#### Custom Variables

Define your own variables in template configuration:

```yaml
# Template configuration
variables:
  author_name:
    type: string
    default: "Anonymous"
    description: "Author name for generated code"
  
  include_examples:
    type: boolean
    default: true
    description: "Include example usage in documentation"
  
  target_go_version:
    type: string
    default: "1.21"
    description: "Target Go version for generated code"
```

Use in templates:
```go
// Generated by {{.Variables.author_name}}
// Target Go version: {{.Variables.target_go_version}}

{{if .Variables.include_examples}}
// Example usage:
// client := NewClient("https://api.example.com")
{{end}}
```

## Project Management

### Project Organization

#### Project Storage

MCPWeaver stores projects in your user directory:

**Windows**: `%APPDATA%\MCPWeaver\projects\`
**macOS**: `~/Library/Application Support/MCPWeaver/projects/`
**Linux**: `~/.config/MCPWeaver/projects/`

Each project is stored in its own directory containing:
- Project metadata (`project.json`)
- OpenAPI specification cache
- Generation history
- User settings and preferences

#### Project Backup

**Automatic Backup**:
- MCPWeaver automatically backs up projects before major operations
- Backups are stored in the `backups/` subdirectory
- Configurable retention period (default: 30 days)
- Automatic cleanup of old backups

**Manual Backup**:
1. Right-click project in project list
2. Select "Create Backup"
3. Choose backup location
4. Include/exclude specific data (cache, history, settings)

#### Project Import/Export

**Export Projects**:
- Export single projects or entire workspace
- Choose between MCPWeaver format or generic ZIP
- Include or exclude cache data for smaller files
- Password protection for sensitive projects

**Import Projects**:
- Import MCPWeaver project files
- Import from ZIP archives
- Merge with existing projects or create new ones
- Conflict resolution for duplicate projects

### Workspace Management

#### Multiple Workspaces

MCPWeaver supports multiple workspaces for different contexts:

**Personal Workspace**: Your default workspace for personal projects
**Work Workspace**: Separate workspace for work-related projects
**Shared Workspace**: Collaborative workspace on shared storage
**Temporary Workspace**: For testing and experimental projects

**Switching Workspaces**:
1. Go to **File > Switch Workspace**
2. Select existing workspace or create new one
3. MCPWeaver will restart with the selected workspace
4. Recent workspaces are available in the File menu

#### Workspace Settings

Configure workspace-specific settings:

**Default Locations**:
- Project storage directory
- Template storage location
- Generated server output directory
- Backup storage location

**Collaboration Settings**:
- Shared template repositories
- Version control integration
- Team member permissions
- Sync settings for shared workspaces

### Project Collaboration

#### Sharing Projects

**Direct File Sharing**:
1. Export project as MCPWeaver file
2. Share file via email, cloud storage, or file sharing
3. Recipient imports file into their MCPWeaver
4. Project appears in their project list

**Version Control Integration**:
- Git integration for project versioning
- Commit project changes automatically
- Share projects via Git repositories
- Collaborate on OpenAPI specifications

#### Team Workflows

**Project Templates**:
- Create standardized project templates for your team
- Include common OpenAPI patterns and best practices
- Share templates via template repositories
- Ensure consistency across team projects

**Review Process**:
- Export generated servers for review
- Share validation results with team members
- Collaborate on OpenAPI specification improvements
- Track changes and maintain project history

### Project Analytics

#### Usage Statistics

MCPWeaver tracks project usage to help you understand your workflow:

**Generation Statistics**:
- Number of servers generated
- Average generation time
- Most used templates
- Error rates and common issues

**Performance Metrics**:
- Project loading times
- Validation performance
- Memory usage by project
- Cache effectiveness

**Activity Timeline**:
- Project creation and modification dates
- Generation history with timestamps
- Template usage over time
- Error and warning trends

#### Reporting

**Project Reports**:
- Generate detailed project reports
- Include validation results, generation history, and performance
- Export reports as PDF or HTML
- Share reports with team members or stakeholders

**Workspace Summary**:
- Overview of all projects in workspace
- Aggregate statistics and trends
- Resource usage analysis
- Recommendations for optimization

## Settings and Preferences

### Application Settings

Access settings via **Tools > Settings** or `Ctrl+,` (`Cmd+,` on macOS).

#### General Settings

**Appearance**:
- **Theme**: Light, Dark, or System (follows OS theme)
- **Font Size**: Adjust interface font size (8-24pt)
- **Font Family**: Choose interface font (system fonts available)
- **UI Scale**: Scale interface for high-DPI displays (100%-200%)

**Language and Region**:
- **Language**: Interface language (English, Spanish, French, German, etc.)
- **Date Format**: Regional date format preferences
- **Number Format**: Decimal separator and thousands separator
- **Time Zone**: For timestamps and scheduling

**Startup Behavior**:
- **Show Welcome Screen**: Display welcome screen on startup
- **Restore Last Session**: Reopen projects from previous session
- **Check for Updates**: Automatically check for MCPWeaver updates
- **Startup Project**: Automatically open specific project on startup

#### Performance Settings

**Memory Management**:
- **Memory Limit**: Maximum memory usage (512MB - 4GB)
- **Cache Size**: Size of validation and generation cache (100MB - 1GB)
- **Garbage Collection**: Automatic memory cleanup frequency
- **Background Processing**: Enable background validation and processing

**Processing Options**:
- **Parallel Processing**: Number of parallel processing threads (1-16)
- **Timeout Settings**: Operation timeout values (30s - 300s)
- **Priority Mode**: Prioritize responsiveness vs. throughput
- **Resource Monitoring**: Enable detailed resource usage tracking

**Network Settings**:
- **Connection Timeout**: Timeout for URL-based imports (10s - 120s)
- **Retry Attempts**: Number of retry attempts for failed operations (1-5)
- **Proxy Configuration**: HTTP/HTTPS proxy settings
- **SSL/TLS Settings**: Certificate validation and security options

#### Editor Settings

**Validation Behavior**:
- **Real-time Validation**: Validate as you type/import
- **Validation Strictness**: Strict, Standard, or Lenient
- **Auto-fix Enabled**: Automatically fix common issues
- **Warning Threshold**: Minimum severity level for warnings

**Display Options**:
- **Line Numbers**: Show line numbers in specification viewer
- **Syntax Highlighting**: Color-code OpenAPI specifications
- **Error Highlighting**: Highlight validation errors in-line
- **Minimap**: Show document minimap for navigation

**Auto-save Options**:
- **Auto-save Enabled**: Automatically save project changes
- **Save Interval**: How often to auto-save (30s - 600s)
- **Backup Before Save**: Create backup before saving changes
- **Save Location**: Where to save auto-saved files

### Project Settings

Configure settings specific to individual projects via **Project > Project Settings**.

#### Generation Defaults

**Output Configuration**:
- **Default Output Directory**: Where to generate servers by default
- **Package Name Pattern**: Template for Go package names
- **Module Path Pattern**: Template for Go module paths
- **File Naming Convention**: How to name generated files

**Template Preferences**:
- **Default Template**: Template to use for new generations
- **Template Variables**: Default values for template variables
- **Custom Template Path**: Path to project-specific templates
- **Template Inheritance**: Inherit from user or system templates

**Quality Settings**:
- **Code Quality Level**: Basic, Standard, or High-quality code generation
- **Documentation Level**: Minimal, Standard, or Comprehensive docs
- **Error Handling**: Basic, Robust, or Comprehensive error handling
- **Testing Level**: None, Basic, or Comprehensive test generation

#### Validation Preferences

**Validation Rules**:
- **Required Fields**: Enforce OpenAPI required fields
- **Best Practices**: Apply OpenAPI best practice validation
- **Custom Rules**: Project-specific validation rules
- **Ignored Warnings**: Suppress specific warning types

**Auto-correction**:
- **Auto-fix Enabled**: Enable automatic issue fixing
- **Fix Categories**: Which types of issues to auto-fix
- **Confirmation Required**: Ask before applying fixes
- **Backup Before Fix**: Create backup before auto-fixing

### Template Settings

#### Template Management

**Template Locations**:
- **System Templates**: Built-in templates (read-only)
- **User Templates**: Personal custom templates
- **Project Templates**: Project-specific templates
- **Shared Templates**: Team or organization templates

**Template Sources**:
- **Local Directories**: Load templates from local directories
- **Git Repositories**: Sync templates from Git repos
- **Template Registries**: Download from template registries
- **HTTP URLs**: Load templates from web URLs

**Template Updates**:
- **Auto-update System Templates**: Keep built-in templates current
- **Check Template Updates**: Periodically check for template updates
- **Update Notifications**: Notify when template updates are available
- **Backup Before Update**: Create backup before updating templates

#### Custom Template Configuration

**Development Settings**:
- **Template Editor**: Choose external editor for template development
- **Syntax Validation**: Real-time template syntax validation
- **Test Data**: Sample OpenAPI specs for template testing
- **Debug Mode**: Enable detailed template processing logs

**Template Variables**:
- **Global Variables**: Variables available to all templates
- **Project Variables**: Variables specific to current project
- **Environment Variables**: Use system environment variables
- **Computed Variables**: Auto-calculated values (dates, versions, etc.)

### Advanced Settings

#### Developer Options

**Debugging**:
- **Debug Mode**: Enable detailed logging and debugging features
- **Log Level**: Minimum log level to display (Error, Warning, Info, Debug)
- **Log File Location**: Where to save application logs
- **Performance Profiling**: Enable performance monitoring and profiling

**Experimental Features**:
- **Beta Features**: Enable experimental features in development
- **Feature Flags**: Toggle specific features on/off
- **Advanced APIs**: Access to advanced MCPWeaver APIs
- **Plugin Support**: Enable third-party plugin support

#### Integration Settings

**External Tools**:
- **Default Code Editor**: Editor to open generated code
- **Diff Tool**: Tool for comparing files and specifications
- **Terminal/Shell**: Default terminal for running commands
- **File Manager**: Default file manager for browsing output

**Version Control**:
- **Git Integration**: Enable Git integration for projects
- **Auto-commit**: Automatically commit project changes
- **Commit Message Templates**: Templates for commit messages
- **Ignored Files**: Files to ignore in version control

#### Security Settings

**Privacy**:
- **Analytics**: Send anonymous usage analytics
- **Error Reporting**: Automatically report crashes and errors
- **Update Checking**: Allow automatic update checking
- **Telemetry**: Send performance and usage data

**Data Protection**:
- **Encrypt Projects**: Encrypt sensitive project data
- **Secure Deletion**: Securely delete temporary files
- **Memory Protection**: Clear sensitive data from memory
- **Audit Logging**: Log security-relevant operations

### Settings Management

#### Import/Export Settings

**Export Settings**:
1. Go to **Tools > Settings > Export Settings**
2. Choose which settings to export (All, Application, Project, Templates)
3. Select export format (MCPWeaver settings file or JSON)
4. Save settings file to desired location

**Import Settings**:
1. Go to **Tools > Settings > Import Settings**
2. Select settings file to import
3. Choose merge strategy (Replace, Merge, or Selective)
4. Review changes before applying

#### Reset Settings

**Reset to Defaults**:
- Reset individual setting categories
- Reset all settings to factory defaults
- Keep user data while resetting preferences
- Create backup before resetting

**Settings Profiles**:
- Create multiple settings profiles (Work, Personal, Testing)
- Switch between profiles quickly
- Share settings profiles with team members
- Apply profile-specific customizations

## Best Practices

### OpenAPI Specification Best Practices

#### Well-Structured Specifications

**Use Descriptive Names**:
```yaml
# ❌ Poor naming
paths:
  /api/v1/u:
    get:
      operationId: getU
      
# ✅ Good naming  
paths:
  /api/v1/users:
    get:
      operationId: getUserList
      summary: Retrieve a list of users
```

**Provide Comprehensive Documentation**:
```yaml
paths:
  /users/{userId}:
    get:
      summary: Retrieve user by ID
      description: |
        Retrieves detailed information for a specific user identified by their unique ID.
        This endpoint requires authentication and returns user profile data including
        contact information, preferences, and account status.
      parameters:
        - name: userId
          in: path
          required: true
          description: The unique identifier for the user
          schema:
            type: integer
            format: int64
            minimum: 1
```

**Use Proper Data Types**:
```yaml
# ✅ Specific, well-defined schemas
components:
  schemas:
    User:
      type: object
      required:
        - id
        - email
        - createdAt
      properties:
        id:
          type: integer
          format: int64
          description: Unique user identifier
        email:
          type: string
          format: email
          description: User's email address
        createdAt:
          type: string
          format: date-time
          description: Account creation timestamp
        profile:
          $ref: '#/components/schemas/UserProfile'
```

#### Error Handling

**Define Standard Error Responses**:
```yaml
components:
  responses:
    BadRequest:
      description: Invalid request parameters
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
    NotFound:
      description: Resource not found
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
    
  schemas:
    Error:
      type: object
      required:
        - code
        - message
      properties:
        code:
          type: string
          description: Error code for programmatic handling
        message:
          type: string
          description: Human-readable error message
        details:
          type: object
          description: Additional error context
```

**Use Consistent Status Codes**:
```yaml
paths:
  /users:
    get:
      responses:
        '200':
          description: Successfully retrieved users
        '400':
          $ref: '#/components/responses/BadRequest'
        '401':
          description: Authentication required
        '403':
          description: Insufficient permissions
        '500':
          description: Internal server error
```

### MCPWeaver Usage Best Practices

#### Project Organization

**Logical Project Structure**:
- One project per API or related set of APIs
- Use descriptive project names and descriptions
- Organize projects by team, product, or environment
- Regular cleanup of unused or outdated projects

**Version Management**:
- Use semantic versioning for your APIs
- Track API changes through specification updates
- Maintain backwards compatibility when possible
- Document breaking changes clearly

**Documentation Strategy**:
- Keep OpenAPI specifications up-to-date with actual API
- Include examples for all endpoints
- Document authentication and authorization requirements
- Provide usage examples and integration guides

#### Generation Strategy

**Template Selection**:
- Use Standard template for most REST APIs
- Choose Lightweight template for simple APIs
- Select Development template for testing scenarios
- Use Production template for deployment-ready servers

**Output Organization**:
- Use consistent directory structures for generated servers
- Include project metadata in generated documentation
- Version generated servers alongside API versions
- Maintain separate outputs for different environments

**Quality Assurance**:
- Always validate specifications before generation
- Test generated servers with real API endpoints
- Review generated code for security and performance
- Document any manual modifications to generated code

#### Performance Optimization

**Specification Optimization**:
- Remove unused schemas and components
- Minimize deep nesting in data structures
- Use references ($ref) to avoid duplication
- Optimize large specifications for validation performance

**Generation Optimization**:
- Enable caching for repeated generations
- Use appropriate memory limits for large specifications
- Monitor generation performance and optimize as needed
- Clean up cache periodically to free disk space

**Resource Management**:
- Monitor MCPWeaver memory usage
- Close unused projects to free resources
- Configure appropriate timeouts for slow operations
- Use background processing for large generations

### Security Best Practices

#### Specification Security

**Authentication Documentation**:
```yaml
components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
      description: |
        JWT token-based authentication. Include the token in the Authorization
        header as 'Bearer <token>'. Tokens expire after 1 hour.

security:
  - BearerAuth: []
```

**Sensitive Data Handling**:
- Avoid including sensitive data in OpenAPI specifications
- Use generic examples instead of real data
- Document data privacy and compliance requirements
- Consider data classification and handling requirements

**API Security Considerations**:
- Document rate limiting and throttling
- Include CORS policy information
- Specify HTTPS requirements
- Document input validation requirements

#### Generated Server Security

**Code Review**:
- Review generated authentication handling
- Verify input validation implementation
- Check error message exposure
- Validate security header handling

**Deployment Security**:
- Use secure deployment practices
- Implement proper logging and monitoring
- Configure firewalls and network security
- Keep dependencies updated

### Integration Best Practices

#### LLM Integration

**Tool Definition Quality**:
- Use clear, descriptive tool names
- Provide comprehensive tool descriptions
- Include parameter validation rules
- Document expected response formats

**Error Handling**:
- Implement graceful error handling for API failures
- Provide meaningful error messages to LLMs
- Log errors for debugging and monitoring
- Include retry logic for transient failures

**Performance Considerations**:
- Optimize for common LLM usage patterns
- Implement appropriate caching strategies
- Monitor API usage and performance
- Set reasonable timeout values

#### Testing and Validation

**Testing Strategy**:
- Test generated servers with real API endpoints
- Validate MCP protocol compliance
- Test error scenarios and edge cases
- Perform integration testing with target LLMs

**Monitoring and Observability**:
- Implement logging for all MCP operations
- Monitor API call success/failure rates
- Track performance metrics
- Set up alerting for critical failures

**Continuous Integration**:
- Automate specification validation
- Include generated server testing in CI/CD
- Monitor for API specification changes
- Automatically update generated servers when appropriate

## Advanced Usage

### Custom Template Development

#### Template Architecture

Understanding MCPWeaver's template system enables powerful customizations:

**Template Hierarchy**:
```
System Templates (Built-in)
├── standard/           # Standard MCP server template
├── lightweight/        # Minimal MCP server template
├── development/        # Development-focused template
└── production/         # Production-ready template

User Templates (Custom)
├── my-company/         # Company-specific template
├── microservice/       # Microservice-focused template
└── legacy-integration/ # Legacy system integration

Project Templates (Project-specific)
├── project-alpha/      # Alpha project customizations
└── beta-testing/       # Beta testing modifications
```

**Template Components**:
```
template/
├── template.yaml       # Template metadata and configuration
├── server/            # Server implementation templates
│   ├── main.go.tmpl   # Main server entry point
│   ├── handlers.go.tmpl # MCP request handlers
│   └── client.go.tmpl # API client implementation
├── models/            # Data model templates
│   ├── models.go.tmpl # Generated data structures
│   └── types.go.tmpl  # Custom type definitions
├── docs/              # Documentation templates
│   ├── README.md.tmpl # Generated README
│   └── API.md.tmpl    # API documentation
└── tests/             # Test generation templates
    ├── client_test.go.tmpl
    └── integration_test.go.tmpl
```

#### Template Variables and Functions

**Advanced Variable Usage**:
```go
{{/* Basic variable access */}}
{{.Spec.Info.Title}}

{{/* Nested data structures */}}
{{range .Spec.Paths}}
  {{range .Operations}}
    Operation: {{.OperationID}}
    Method: {{.HTTPMethod}}
    {{range .Parameters}}
      Param: {{.Name}} ({{.Type}})
    {{end}}
  {{end}}
{{end}}

{{/* Conditional logic with complex conditions */}}
{{if and .HasAuthentication (not .IsReadOnly)}}
// Authentication required for write operations
{{end}}

{{/* Custom variable access */}}
{{.Variables.company_name}}
{{.Variables.target_environment}}
```

**Helper Functions**:
```go
{{/* String manipulation */}}
{{.OperationID | camelCase}}        # getUserList
{{.OperationID | pascalCase}}       # GetUserList
{{.OperationID | snakeCase}}        # get_user_list
{{.OperationID | kebabCase}}        # get-user-list

{{/* Type conversion */}}
{{.SchemaType | goType}}            # Convert OpenAPI type to Go type
{{.Parameter | goParameter}}        # Convert to Go parameter

{{/* Documentation helpers */}}
{{.Description | goComment}}        # Format as Go comment
{{.Description | markdownEscape}}   # Escape for Markdown

{{/* Code generation helpers */}}
{{.Schema | generateStruct}}        # Generate Go struct
{{.Operation | generateHandler}}    # Generate handler function
```

#### Creating Production-Ready Templates

**Enterprise Template Example**:
```yaml
# template.yaml
name: "enterprise-mcp-server"
description: "Production-ready MCP server with enterprise features"
version: "1.0.0"
author: "Your Company"
license: "MIT"

variables:
  company_name:
    type: string
    required: true
    description: "Company name for generated code"
  
  monitoring_enabled:
    type: boolean
    default: true
    description: "Enable Prometheus monitoring"
  
  auth_provider:
    type: string
    enum: ["oauth2", "jwt", "apikey"]
    default: "jwt"
    description: "Authentication provider"

  rate_limiting:
    type: object
    properties:
      enabled:
        type: boolean
        default: true
      requests_per_minute:
        type: integer
        default: 1000
        minimum: 1
        maximum: 10000

features:
  - monitoring
  - rate_limiting
  - authentication
  - logging
  - error_tracking
  - health_checks
  - graceful_shutdown
```

**Advanced Server Template**:
```go
// main.go.tmpl
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    {{if .Variables.monitoring_enabled}}
    "github.com/prometheus/client_golang/prometheus/promhttp"
    {{end}}
    
    "{{.ModulePath}}/internal/server"
    "{{.ModulePath}}/internal/config"
    {{if .Variables.monitoring_enabled}}
    "{{.ModulePath}}/internal/monitoring"
    {{end}}
)

func main() {
    cfg := config.Load()
    
    {{if .Variables.monitoring_enabled}}
    // Initialize monitoring
    monitoring.Setup(cfg.Monitoring)
    {{end}}
    
    // Create MCP server
    srv := server.New(cfg)
    
    // Setup graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan
        
        log.Println("Shutting down gracefully...")
        cancel()
    }()
    
    // Start server
    if err := srv.Run(ctx); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
```

### API Integration Patterns

#### Complex Authentication Scenarios

**Multi-Auth Support**:
```yaml
# OpenAPI specification with multiple auth schemes
components:
  securitySchemes:
    OAuth2:
      type: oauth2
      flows:
        authorizationCode:
          authorizationUrl: https://auth.example.com/oauth/authorize
          tokenUrl: https://auth.example.com/oauth/token
          scopes:
            read: Read access
            write: Write access
    
    ApiKey:
      type: apiKey
      in: header
      name: X-API-Key
    
    Bearer:
      type: http
      scheme: bearer

security:
  - OAuth2: [read, write]
  - ApiKey: []
  - Bearer: []
```

**Template Implementation**:
```go
// Generated authentication handling
type AuthConfig struct {
    OAuth2   *OAuth2Config  `json:"oauth2,omitempty"`
    APIKey   *APIKeyConfig  `json:"apiKey,omitempty"`
    Bearer   *BearerConfig  `json:"bearer,omitempty"`
}

func (c *Client) authenticate(req *http.Request) error {
    {{if .HasOAuth2}}
    if c.authConfig.OAuth2 != nil {
        token, err := c.getOAuth2Token()
        if err != nil {
            return fmt.Errorf("OAuth2 authentication failed: %w", err)
        }
        req.Header.Set("Authorization", "Bearer "+token)
        return nil
    }
    {{end}}
    
    {{if .HasAPIKey}}
    if c.authConfig.APIKey != nil {
        req.Header.Set("{{.APIKeyHeader}}", c.authConfig.APIKey.Key)
        return nil
    }
    {{end}}
    
    {{if .HasBearer}}
    if c.authConfig.Bearer != nil {
        req.Header.Set("Authorization", "Bearer "+c.authConfig.Bearer.Token)
        return nil
    }
    {{end}}
    
    return errors.New("no authentication method configured")
}
```

#### Advanced Error Handling

**Comprehensive Error Mapping**:
```go
// Error handling template
type APIError struct {
    Code       string            `json:"code"`
    Message    string            `json:"message"`
    Details    map[string]string `json:"details,omitempty"`
    StatusCode int               `json:"statusCode"`
    Retryable  bool              `json:"retryable"`
}

func (c *Client) handleError(resp *http.Response) error {
    switch resp.StatusCode {
    {{range .ErrorResponses}}
    case {{.StatusCode}}:
        return &APIError{
            Code:       "{{.Code}}",
            Message:    "{{.Message}}",
            StatusCode: {{.StatusCode}},
            Retryable:  {{.Retryable}},
        }
    {{end}}
    default:
        return &APIError{
            Code:       "UNKNOWN_ERROR",
            Message:    fmt.Sprintf("Unexpected status code: %d", resp.StatusCode),
            StatusCode: resp.StatusCode,
            Retryable:  false,
        }
    }
}
```

### Performance Optimization

#### Large Specification Handling

**Chunked Processing**:
- Process large specifications in chunks
- Implement progressive validation
- Use streaming for memory efficiency
- Cache intermediate results

**Memory Management**:
```go
// Optimized specification processing
type SpecProcessor struct {
    chunkSize    int
    maxMemory    int64
    cacheEnabled bool
}

func (p *SpecProcessor) ProcessLargeSpec(spec *openapi.Spec) error {
    // Monitor memory usage
    if p.getMemoryUsage() > p.maxMemory {
        return errors.New("memory limit exceeded")
    }
    
    // Process in chunks
    for chunk := range p.chunkOperations(spec.Paths) {
        if err := p.processChunk(chunk); err != nil {
            return fmt.Errorf("chunk processing failed: %w", err)
        }
        
        // Garbage collect between chunks
        runtime.GC()
    }
    
    return nil
}
```

#### Caching Strategies

**Multi-Level Caching**:
```go
type CacheManager struct {
    specCache        *sync.Map // In-memory specification cache
    validationCache  *sync.Map // Validation result cache
    generationCache  *sync.Map // Generation artifact cache
    diskCache        string    // Disk-based cache location
}

func (cm *CacheManager) GetValidationResult(specHash string) (*ValidationResult, bool) {
    // Check memory cache first
    if result, ok := cm.validationCache.Load(specHash); ok {
        return result.(*ValidationResult), true
    }
    
    // Check disk cache
    if result := cm.loadFromDisk(specHash); result != nil {
        cm.validationCache.Store(specHash, result)
        return result, true
    }
    
    return nil, false
}
```

### Automation and Scripting

#### CLI Integration

MCPWeaver supports command-line operation for automation:

```bash
# Automated specification validation
mcpweaver validate --spec api.yaml --output validation.json

# Batch generation
mcpweaver generate --project projects/ --template enterprise --output servers/

# Template management
mcpweaver template install --source https://github.com/company/templates.git
mcpweaver template list --format json

# Performance monitoring
mcpweaver monitor --project myproject --duration 24h --output metrics.json
```

#### CI/CD Integration

**GitHub Actions Example**:
```yaml
name: MCP Server Generation
on:
  push:
    paths: ['api-specs/**']

jobs:
  generate-servers:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup MCPWeaver
        uses: mcpweaver/setup-action@v1
        with:
          version: latest
      
      - name: Validate Specifications
        run: |
          for spec in api-specs/*.yaml; do
            mcpweaver validate --spec "$spec" --strict
          done
      
      - name: Generate MCP Servers
        run: |
          mcpweaver generate \
            --specs api-specs/ \
            --template production \
            --output generated-servers/ \
            --parallel 4
      
      - name: Upload Generated Servers
        uses: actions/upload-artifact@v4
        with:
          name: mcp-servers
          path: generated-servers/
```

#### Monitoring and Alerting

**Performance Monitoring**:
```go
// Monitoring integration
type Monitor struct {
    prometheus *prometheus.Registry
    logger     *zap.Logger
    alerts     chan Alert
}

func (m *Monitor) TrackGeneration(duration time.Duration, success bool) {
    // Track generation metrics
    generationDuration.Observe(duration.Seconds())
    
    if success {
        generationSuccess.Inc()
    } else {
        generationFailures.Inc()
        m.alerts <- Alert{
            Type:    "generation_failure",
            Message: "MCP server generation failed",
            Time:    time.Now(),
        }
    }
}
```

This concludes the comprehensive MCPWeaver User Guide. The guide covers everything from basic usage to advanced customization and automation, providing users with the knowledge they need to effectively use MCPWeaver for their API integration projects.