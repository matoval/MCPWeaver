# MCPWeaver Frequently Asked Questions

This FAQ covers the most common questions about MCPWeaver. For detailed troubleshooting, see our [Troubleshooting Guide](TROUBLESHOOTING.md).

## Table of Contents

- [General Questions](#general-questions)
- [Installation and Setup](#installation-and-setup)
- [Using MCPWeaver](#using-mcpweaver)
- [OpenAPI and MCP](#openapi-and-mcp)
- [Generation and Output](#generation-and-output)
- [Performance and Limitations](#performance-and-limitations)
- [Troubleshooting](#troubleshooting)
- [Development and Contributing](#development-and-contributing)
- [Licensing and Commercial Use](#licensing-and-commercial-use)

## General Questions

### What is MCPWeaver?

**Q: What does MCPWeaver do?**

A: MCPWeaver is a desktop application that converts OpenAPI 3.0 specifications into Model Context Protocol (MCP) servers. It provides a simple, fast way to transform REST API documentation into MCP-compatible servers that can be used with AI/LLM applications.

**Q: What is the Model Context Protocol (MCP)?**

A: MCP is a protocol developed by Anthropic that allows AI models to securely connect to external data sources and tools. It provides a standardized way for LLMs to interact with APIs, databases, and other services while maintaining security and user control.

**Q: How is MCPWeaver different from openapi2mcp?**

A: MCPWeaver is a desktop GUI application built on top of the openapi2mcp core functionality. Key differences:

- **User Interface**: Graphical interface vs command-line
- **Project Management**: Built-in project tracking and history
- **Real-time Validation**: Instant feedback on OpenAPI specs
- **Template Customization**: Visual template management
- **Cross-platform**: Native desktop app for Windows, macOS, and Linux
- **Performance Monitoring**: Built-in performance tracking

**Q: Is MCPWeaver free to use?**

A: Yes! MCPWeaver is open-source software licensed under AGPL v3. It's free for personal and commercial use. The source code is available on GitHub for transparency and community contributions.

**Q: What platforms does MCPWeaver support?**

A: MCPWeaver supports:
- **Windows**: Windows 10/11 (x64)
- **macOS**: macOS 10.15+ (Intel and Apple Silicon)
- **Linux**: Ubuntu 18.04+, Fedora, Arch, and other distributions (x64)

### Who Should Use MCPWeaver?

**Q: Who is MCPWeaver designed for?**

A: MCPWeaver is designed for:
- **API Developers**: Converting existing REST APIs to MCP
- **AI/LLM Engineers**: Integrating APIs with AI applications
- **DevOps Teams**: Modernizing legacy API infrastructure
- **Product Managers**: Enabling AI features for existing services
- **Researchers**: Exploring AI-API integration patterns

**Q: Do I need programming experience to use MCPWeaver?**

A: Basic understanding of APIs is helpful, but MCPWeaver is designed to be user-friendly:
- **No Coding Required**: Generate servers through the GUI
- **Visual Validation**: See errors and suggestions in real-time
- **Templates**: Pre-built templates for common scenarios
- **Documentation**: Comprehensive guides and examples

However, you'll need some technical knowledge to:
- Deploy the generated MCP servers
- Understand OpenAPI specifications
- Configure the generated code for your environment

## Installation and Setup

### System Requirements

**Q: What are the minimum system requirements?**

A: **Minimum Requirements**:
- **OS**: Windows 10, macOS 10.15, or Linux (Ubuntu 18.04+)
- **RAM**: 2 GB available memory
- **Storage**: 100 MB free disk space
- **CPU**: Any modern processor (x64 or ARM64)

**Recommended Requirements**:
- **RAM**: 4 GB or more
- **Storage**: 500 MB for app and project data
- **CPU**: Multi-core processor for faster generation
- **Display**: 1920x1080 or higher resolution

**Q: Can I run MCPWeaver on older systems?**

A: MCPWeaver requires relatively modern systems due to its web-based UI components:
- **Windows**: Must support WebView2 (Windows 10 version 1903+)
- **macOS**: Requires modern WebKit (macOS 10.15+)
- **Linux**: Needs GTK 3.20+ and WebKitGTK 2.24+

For older systems, consider using the original openapi2mcp CLI tool.

### Installation Issues

**Q: Why does Windows block MCPWeaver installation?**

A: Windows Defender SmartScreen may block unsigned executables. This is normal for new applications:
1. Click "More info" in the SmartScreen dialog
2. Click "Run anyway" to proceed
3. Verify the download came from the official GitHub releases page

Future releases will include extended validation (EV) code signing to eliminate this warning.

**Q: How do I install MCPWeaver without administrator privileges?**

A: Use the portable installation method:
1. Download the ZIP archive instead of the installer
2. Extract to a folder in your user directory
3. Run MCPWeaver.exe directly
4. Create a `portable.txt` file in the same directory for portable mode

**Q: Can I install MCPWeaver on a network drive or USB?**

A: Yes, using portable mode:
1. Use the portable ZIP/tar.gz distribution
2. Create `portable.txt` file in the application directory
3. All settings and data will be stored relative to the application
4. Ensure the drive has write permissions

### First-Time Setup

**Q: What should I do after installing MCPWeaver?**

A: **Quick Start Steps**:
1. **Launch** MCPWeaver from your applications menu
2. **Complete** the welcome wizard (choose theme, set default directories)
3. **Test** with the built-in Petstore API example
4. **Import** your first OpenAPI specification
5. **Generate** your first MCP server
6. **Review** the generated code and documentation

**Q: Where does MCPWeaver store its data?**

A: **Data Locations** (varies by platform):

**Windows**:
- **Settings**: `%APPDATA%\MCPWeaver\`
- **Projects**: `%APPDATA%\MCPWeaver\projects\`
- **Logs**: `%APPDATA%\MCPWeaver\logs\`
- **Cache**: `%LOCALAPPDATA%\MCPWeaver\cache\`

**macOS**:
- **Settings**: `~/Library/Application Support/MCPWeaver/`
- **Projects**: `~/Library/Application Support/MCPWeaver/projects/`
- **Logs**: `~/Library/Logs/MCPWeaver/`
- **Cache**: `~/Library/Caches/MCPWeaver/`

**Linux**:
- **Settings**: `~/.config/MCPWeaver/`
- **Projects**: `~/.local/share/MCPWeaver/projects/`
- **Logs**: `~/.local/share/MCPWeaver/logs/`
- **Cache**: `~/.cache/MCPWeaver/`

## Using MCPWeaver

### Basic Usage

**Q: How do I import an OpenAPI specification?**

A: **Three Ways to Import**:

1. **From File**:
   - Click "Import Spec" → "From File"
   - Select your `.json`, `.yaml`, or `.yml` file
   - MCPWeaver will validate and load the specification

2. **From URL**:
   - Click "Import Spec" → "From URL"
   - Enter the URL (e.g., `https://petstore3.swagger.io/api/v3/openapi.json`)
   - MCPWeaver will download and validate the specification

3. **Direct Input**:
   - Click "Import Spec" → "Paste Content"
   - Paste the OpenAPI specification directly
   - MCPWeaver will parse and validate the content

**Q: What OpenAPI versions are supported?**

A: **Supported Versions**:
- ✅ **OpenAPI 3.0.x** (fully supported)
- ✅ **OpenAPI 3.1.x** (fully supported)
- ❌ **OpenAPI 2.0** (Swagger) - not supported

**Migration Path**: Use the [Swagger Editor](https://editor.swagger.io/) to convert OpenAPI 2.0 specs to 3.0+ format.

**Q: How do I fix validation errors in my OpenAPI spec?**

A: **Common Fixes**:

1. **Missing Required Fields**:
   ```yaml
   # Add required fields
   openapi: "3.0.0"
   info:
     title: "My API"
     version: "1.0.0"
   paths: {}
   ```

2. **Invalid Path Format**:
   ```yaml
   # Paths must start with /
   paths:
     "/users":  # ✅ Correct
       get: {}
     "users":   # ❌ Invalid
       get: {}
   ```

3. **Missing Operation IDs**:
   ```yaml
   # Add operationId for code generation
   paths:
     "/users":
       get:
         operationId: "getUsers"  # Required for generation
   ```

Use the built-in validator suggestions for specific fixes.

### Project Management

**Q: How do I organize multiple API projects?**

A: **Project Organization Tips**:

1. **Create Separate Projects** for each API
2. **Use Descriptive Names** (e.g., "User Service API v2")
3. **Add Descriptions** to explain the project purpose
4. **Use Tags** to categorize projects (e.g., "internal", "public", "v1")
5. **Set Output Directories** to keep generated code organized

**Q: Can I share projects with team members?**

A: **Sharing Options**:

1. **Export Project**:
   - Right-click project → "Export Project"
   - Share the `.mcpweaver` file
   - Team members can import using "Import Project"

2. **Share OpenAPI Specs**:
   - Export the OpenAPI specification
   - Team members can import and create their own projects

3. **Version Control**:
   - Store OpenAPI specs in Git repositories
   - Share repository URLs for import
   - Use MCPWeaver to generate from shared specs

**Q: How do I backup my projects?**

A: **Backup Strategies**:

1. **Built-in Export**:
   - Settings → "Export All Projects"
   - Creates a backup file with all projects

2. **Manual Backup**:
   - Copy the entire MCPWeaver data directory
   - Includes projects, settings, and cache

3. **Selective Backup**:
   - Export individual projects as needed
   - Store OpenAPI specs in version control

### Generation and Templates

**Q: What gets generated when I create an MCP server?**

A: **Generated Files**:
- `server.go` - Main MCP server implementation
- `handlers.go` - API endpoint handlers
- `types.go` - Data type definitions
- `client.go` - HTTP client for API calls
- `config.go` - Configuration management
- `main.go` - Server entry point
- `go.mod` - Go module file
- `README.md` - Setup and usage instructions
- `docker/` - Docker configuration files
- `tests/` - Basic test files

**Q: Can I customize the generated code?**

A: **Customization Options**:

1. **Templates**: Create custom generation templates
2. **Settings**: Adjust generation parameters
3. **Post-Generation**: Modify generated code as needed
4. **Hooks**: Add custom code hooks (planned feature)

**Q: How do I use custom templates?**

A: **Template Management**:

1. **Access Templates**:
   - Settings → "Templates" → "Manage Templates"

2. **Create Custom Template**:
   - Click "New Template"
   - Choose base template to modify
   - Edit using Go template syntax

3. **Apply Template**:
   - Select template during project creation
   - Or change template in project settings

Templates use Go's `text/template` syntax with custom functions for OpenAPI data access.

## OpenAPI and MCP

### OpenAPI Specifications

**Q: What makes a good OpenAPI specification for MCP generation?**

A: **Best Practices**:

1. **Complete Operation IDs**:
   ```yaml
   paths:
     "/users":
       get:
         operationId: "listUsers"  # Required
       post:
         operationId: "createUser" # Required
   ```

2. **Descriptive Schemas**:
   ```yaml
   components:
     schemas:
       User:
         type: object
         properties:
           id:
             type: string
             description: "Unique user identifier"
   ```

3. **Response Examples**:
   ```yaml
   responses:
     "200":
       description: "Success"
       content:
         application/json:
           schema:
             $ref: "#/components/schemas/User"
           example:
             id: "123"
             name: "John Doe"
   ```

4. **Error Responses**:
   ```yaml
   responses:
     "404":
       description: "User not found"
     "400":
       description: "Invalid request"
   ```

**Q: What OpenAPI features are not supported?**

A: **Limitations**:
- **Callbacks**: Not directly mapped to MCP
- **Links**: Not implemented in current version
- **Webhooks**: MCP is request-response based
- **Complex Authentication**: Only basic auth and API keys
- **File Uploads**: Limited support for multipart/form-data

**Q: How does MCPWeaver handle authentication in OpenAPI specs?**

A: **Supported Authentication**:

1. **API Key Authentication**:
   ```yaml
   components:
     securitySchemes:
       ApiKeyAuth:
         type: apiKey
         in: header
         name: X-API-Key
   ```

2. **HTTP Basic Authentication**:
   ```yaml
   components:
     securitySchemes:
       BasicAuth:
         type: http
         scheme: basic
   ```

3. **Bearer Token Authentication**:
   ```yaml
   components:
     securitySchemes:
       BearerAuth:
         type: http
         scheme: bearer
   ```

**Not Currently Supported**:
- OAuth 2.0 flows
- OpenID Connect
- Complex custom authentication

### MCP Concepts

**Q: How does MCPWeaver map OpenAPI operations to MCP tools?**

A: **Mapping Strategy**:

1. **HTTP Methods → MCP Tools**:
   - `GET /users` → `listUsers` tool
   - `POST /users` → `createUser` tool
   - `GET /users/{id}` → `getUser` tool
   - `PUT /users/{id}` → `updateUser` tool
   - `DELETE /users/{id}` → `deleteUser` tool

2. **Parameters → Tool Arguments**:
   - Path parameters become required arguments
   - Query parameters become optional arguments
   - Request body becomes structured argument

3. **Responses → Tool Results**:
   - Success responses become tool results
   - Error responses become error handling

**Q: What MCP capabilities does the generated server support?**

A: **Supported MCP Features**:
- ✅ **Tools**: All OpenAPI operations become MCP tools
- ✅ **Resources**: API data exposed as MCP resources
- ✅ **Prompts**: Generated prompts for common operations
- ✅ **Sampling**: Request/response examples for LLMs
- ⚠️ **Logging**: Basic logging (enhanced logging planned)
- ❌ **Notifications**: Not currently implemented

**Q: How do I configure the generated MCP server?**

A: **Configuration Options**:

1. **Environment Variables**:
   ```bash
   export API_BASE_URL="https://api.example.com"
   export API_KEY="your-api-key"
   export MCP_SERVER_PORT="3000"
   ```

2. **Configuration File** (`config.yaml`):
   ```yaml
   api:
     baseURL: "https://api.example.com"
     apiKey: "${API_KEY}"
     timeout: "30s"
   
   mcp:
     port: 3000
     logLevel: "info"
   ```

3. **Command-Line Flags**:
   ```bash
   ./mcp-server --api-url="https://api.example.com" --port=3000
   ```

## Performance and Limitations

### Performance

**Q: How fast is MCPWeaver?**

A: **Performance Benchmarks**:
- **Startup Time**: < 2 seconds (target: < 1.5s)
- **Memory Usage**: < 50MB for typical use
- **Small APIs** (< 10 endpoints): < 1 second generation
- **Medium APIs** (10-100 endpoints): < 3 seconds generation
- **Large APIs** (100+ endpoints): < 10 seconds generation

**Q: What affects generation performance?**

A: **Performance Factors**:

1. **Specification Size**: Larger specs take longer
2. **Template Complexity**: Custom templates add overhead
3. **System Resources**: CPU, RAM, and disk speed
4. **Antivirus Software**: Real-time scanning slows file operations
5. **Network Latency**: For URL imports
6. **Disk Type**: SSD vs HDD for output directory

**Q: How can I improve performance?**

A: **Optimization Tips**:

1. **System Optimization**:
   - Use SSD storage for output directory
   - Ensure sufficient RAM (4GB+)
   - Close unnecessary applications
   - Exclude MCPWeaver from antivirus real-time scanning

2. **Specification Optimization**:
   - Remove unnecessary examples and descriptions
   - Use `$ref` for repeated schemas
   - Split large APIs into smaller specs

3. **Application Settings**:
   - Enable generation caching
   - Disable real-time validation for large specs
   - Use local files instead of URL imports

### Limitations

**Q: What are the current limitations of MCPWeaver?**

A: **Known Limitations**:

1. **OpenAPI Support**:
   - OpenAPI 2.0 (Swagger) not supported
   - Some advanced features not implemented
   - Complex authentication schemes limited

2. **File Size Limits**:
   - Maximum OpenAPI spec size: 10MB
   - Memory usage scales with spec complexity
   - Very large APIs may require splitting

3. **Platform Limitations**:
   - Requires modern operating systems
   - WebView2/WebKit dependency
   - Limited offline functionality

4. **Generation Limits**:
   - Go language output only (other languages planned)
   - Template customization requires technical knowledge
   - Some edge cases in complex schemas

**Q: Are there any security considerations?**

A: **Security Features**:
- ✅ **Local Processing**: No data sent to external servers
- ✅ **Input Validation**: Strict validation of OpenAPI specs
- ✅ **Sandboxed Templates**: Template execution is isolated
- ✅ **File System Restrictions**: Limited file access patterns

**Security Limitations**:
- Generated code inherits API security model
- Users responsible for securing generated servers
- Template security depends on user-created content

## Troubleshooting

### Common Issues

**Q: MCPWeaver won't start. What should I do?**

A: **Quick Fixes**:
1. **Restart your computer** (clears file locks and memory)
2. **Check system requirements** (OS version, available RAM)
3. **Delete database file** (forces clean restart)
4. **Reinstall MCPWeaver** (fixes corrupted installation)

For detailed troubleshooting, see our [Troubleshooting Guide](TROUBLESHOOTING.md).

**Q: My OpenAPI spec shows validation errors but works in other tools. Why?**

A: **Possible Causes**:

1. **Stricter Validation**: MCPWeaver enforces OpenAPI 3.0+ compliance
2. **Required Fields**: Some tools are more lenient about missing fields
3. **Cache Issues**: Clear validation cache and retry
4. **Version Differences**: Ensure you're using OpenAPI 3.0+, not 2.0

**Q: Generation fails with "Permission Denied" error. How do I fix this?**

A: **Permission Solutions**:

1. **Choose Different Output Directory**:
   - Use your home directory or Documents folder
   - Avoid system directories like Program Files

2. **Fix Directory Permissions**:
   ```bash
   # Windows (run as Administrator)
   icacls "C:\path\to\output" /grant %USERNAME%:F
   
   # macOS/Linux
   chmod 755 /path/to/output
   sudo chown $USER:$USER /path/to/output
   ```

3. **Run as Administrator/Root** (temporary solution)

**Q: The generated MCP server doesn't compile. What's wrong?**

A: **Common Compilation Issues**:

1. **Invalid Go Identifiers**: OpenAPI names must be valid Go identifiers
2. **Reserved Keywords**: Avoid Go reserved words in schema names
3. **Missing Dependencies**: Ensure Go modules are properly initialized
4. **Template Errors**: Reset to default template if using custom templates

### Getting Help

**Q: Where can I get help with MCPWeaver?**

A: **Support Channels**:

1. **Documentation**:
   - [User Guide](USER_GUIDE.md) - Complete usage instructions
   - [Troubleshooting Guide](TROUBLESHOOTING.md) - Problem resolution
   - [API Documentation](API.md) - Technical reference

2. **Community Support**:
   - [GitHub Issues](https://github.com/matoval/MCPWeaver/issues) - Bug reports and feature requests
   - [GitHub Discussions](https://github.com/matoval/MCPWeaver/discussions) - General questions
   - [Stack Overflow](https://stackoverflow.com/questions/tagged/mcpweaver) - Programming questions

3. **Professional Support**:
   - Enterprise support available for commercial users
   - Custom development and integration services
   - Training and consultation

**Q: How do I report a bug?**

A: **Bug Report Steps**:

1. **Search Existing Issues**: Check if the bug is already reported
2. **Use Bug Report Template**: Provides structure for essential information
3. **Include Diagnostic Information**:
   - Operating system and version
   - MCPWeaver version
   - Steps to reproduce
   - Error messages and logs
   - Sample OpenAPI specification (if relevant)

4. **Generate Diagnostic Report**: Help → "Generate Diagnostic Report"

## Development and Contributing

### Contributing

**Q: How can I contribute to MCPWeaver?**

A: **Contribution Types**:

1. **Code Contributions**:
   - Bug fixes and feature improvements
   - New templates and generators
   - Performance optimizations
   - Test coverage improvements

2. **Documentation**:
   - User guides and tutorials
   - API documentation
   - Code examples and samples
   - Translation to other languages

3. **Testing and Feedback**:
   - Bug reports and reproduction cases
   - Feature requests and use cases
   - Performance testing and benchmarks
   - Usability feedback

4. **Community Support**:
   - Answer questions in discussions
   - Help other users troubleshoot issues
   - Share templates and configurations
   - Write blog posts and tutorials

**Q: How do I set up a development environment?**

A: **Development Setup**:

1. **Prerequisites**:
   - Go 1.23+ (latest stable recommended)
   - Node.js 18+ (for frontend development)
   - Wails CLI v2.10.1+
   - Git for version control

2. **Clone and Setup**:
   ```bash
   git clone https://github.com/matoval/MCPWeaver.git
   cd MCPWeaver
   go mod tidy
   cd frontend && npm install && cd ..
   ```

3. **Run Development Server**:
   ```bash
   wails dev
   ```

See the [Developer Guide](DEVELOPER.md) for detailed instructions.

**Q: What programming languages/technologies does MCPWeaver use?**

A: **Technology Stack**:

**Backend (Go)**:
- Go 1.23+ for core application logic
- SQLite for local data storage
- Standard library for most functionality
- Third-party libraries for OpenAPI parsing

**Frontend (React/TypeScript)**:
- React 18 for user interface
- TypeScript for type safety
- Modern CSS for styling
- Wails runtime for backend communication

**Framework**:
- Wails v2 for desktop application framework
- Cross-platform native desktop integration
- WebView for modern UI rendering

### Extending MCPWeaver

**Q: Can I create custom code generators?**

A: **Customization Options**:

1. **Custom Templates**:
   - Modify existing Go templates
   - Add new file templates
   - Create language-specific generators

2. **Plugin System** (planned):
   - Runtime plugin loading
   - Custom generation pipelines
   - Extended language support

3. **Fork and Modify**:
   - Full source code access under AGPL v3
   - Modify core generation logic
   - Add new features and capabilities

**Q: How do I create a custom template?**

A: **Template Creation**:

1. **Access Template Manager**:
   - Settings → "Templates" → "Create New Template"

2. **Template Syntax**:
   - Uses Go `text/template` syntax
   - Access OpenAPI data via template variables
   - Custom functions for common operations

3. **Example Template**:
   ```go
   // Generated MCP server for {{.Info.Title}}
   package main
   
   import (
       "context"
       "log"
   )
   
   {{range .Paths}}
   // {{.OperationID}} handles {{.HTTPMethod}} {{.Path}}
   func {{.OperationID}}(ctx context.Context, args map[string]interface{}) (interface{}, error) {
       // Implementation here
       return nil, nil
   }
   {{end}}
   ```

4. **Testing Templates**:
   - Use with simple OpenAPI specs first
   - Test generation and compilation
   - Iterate and refine

## Licensing and Commercial Use

### Licensing

**Q: What license does MCPWeaver use?**

A: MCPWeaver is licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)**.

**What this means**:
- ✅ **Free to use** for personal and commercial projects
- ✅ **Free to modify** and distribute
- ✅ **Source code access** guaranteed
- ⚠️ **Copyleft license** - derivative works must also be open source
- ⚠️ **Network use trigger** - if you modify and run as a service, you must share modifications

**Q: Can I use MCPWeaver in commercial projects?**

A: **Yes, commercial use is allowed**, but with important considerations:

**Allowed**:
- Using MCPWeaver to generate MCP servers for your business
- Generating servers for commercial APIs
- Including generated code in commercial products
- Using MCPWeaver in corporate environments

**Requirements**:
- Generated code is yours to use however you want
- If you modify MCPWeaver itself and distribute it, you must share modifications
- If you run a modified version as a service, you must provide source code

**Q: What about the generated code? What license does it have?**

A: **Generated Code Ownership**:
- Generated MCP servers belong to you
- No licensing restrictions on generated code
- You can use generated code in any project (commercial or open source)
- You can license generated code however you choose

**Q: Do I need to pay for MCPWeaver?**

A: **MCPWeaver is completely free**:
- No licensing fees
- No subscription costs  
- No usage restrictions
- No user limits

**Optional Paid Services**:
- Priority support for enterprise users
- Custom development services
- Training and consulting

### Enterprise Use

**Q: Is MCPWeaver suitable for enterprise use?**

A: **Enterprise Considerations**:

**Advantages**:
- Local processing (no cloud dependencies)
- Full source code access for security reviews
- Active development and community support
- Professional support options available

**Considerations**:
- AGPL license requires understanding of copyleft implications
- Self-hosted deployment and maintenance
- Support primarily through community channels

**Q: Do you offer enterprise licenses?**

A: Currently, MCPWeaver is only available under AGPL v3. However, we're considering:
- Commercial licenses for specific use cases
- Enterprise support packages
- Custom development agreements

Contact us at enterprise@mcpweaver.dev for enterprise discussions.

**Q: Can I get professional support?**

A: **Support Options**:

**Community Support** (Free):
- GitHub issues and discussions
- Documentation and guides
- Community forums and chat

**Professional Support** (Paid):
- Priority issue resolution
- Custom feature development
- Integration assistance and consulting
- Training for development teams
- SLA-backed support agreements

Contact support@mcpweaver.dev for professional support inquiries.

---

## Quick Reference

### Essential Links
- **Download**: [GitHub Releases](https://github.com/matoval/MCPWeaver/releases/latest)
- **Documentation**: [User Guide](USER_GUIDE.md) | [API Docs](API.md) | [Installation](INSTALLATION.md)
- **Support**: [Issues](https://github.com/matoval/MCPWeaver/issues) | [Discussions](https://github.com/matoval/MCPWeaver/discussions)
- **Community**: [Discord](https://discord.gg/mcpweaver) | [Stack Overflow](https://stackoverflow.com/questions/tagged/mcpweaver)

### Quick Commands
```bash
# Check version
mcpweaver --version

# Generate diagnostic report
mcpweaver --diagnostic

# Reset settings
mcpweaver --reset-settings

# Clear all caches
mcpweaver --clear-cache
```

### Emergency Contacts
- **Security Issues**: security@mcpweaver.dev
- **Enterprise Support**: enterprise@mcpweaver.dev
- **General Support**: support@mcpweaver.dev

---

**Still have questions?** Check our [complete documentation](README.md) or [join the community](https://github.com/matoval/MCPWeaver/discussions) for help! 🚀