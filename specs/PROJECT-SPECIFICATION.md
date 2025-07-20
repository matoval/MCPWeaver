# MCPWeaver Project Specification

## Overview

MCPWeaver is an open-source desktop application that transforms OpenAPI specifications into Model Context Protocol (MCP) servers. This project pivots from the hosted SaaS openapi2mcp to a lightweight, user-controlled desktop application that runs entirely on the user's machine.

## Project Vision

Create the simplest, most efficient desktop application for generating MCP servers from OpenAPI specifications, with minimal complexity and maximum user control.

## Core Principles

- **Open Source**: Fully open-source with MIT license
- **Desktop First**: No cloud dependencies, everything runs locally
- **Lightweight**: Minimal resource usage and fast startup
- **Simple**: Intuitive interface with minimal learning curve
- **Efficient**: Fast processing with real-time feedback
- **Portable**: Single binary distribution with no external dependencies

## Technical Requirements

### Platform Support
- **Primary**: Windows, macOS, Linux
- **Architecture**: x64 and ARM64 support
- **Minimum System**: 2GB RAM, 100MB disk space

### Core Functionality
- OpenAPI 3.0 specification parsing and validation
- MCP server code generation (Go-based)
- Real-time generation progress tracking
- Generated server testing and validation
- Project management and history
- Template customization support

### User Experience Requirements
- **Startup Time**: < 2 seconds cold start
- **Processing Time**: < 5 seconds for typical OpenAPI specs
- **Memory Usage**: < 50MB base memory footprint
- **Responsiveness**: UI remains responsive during generation
- **Error Recovery**: Graceful handling of invalid specifications

## Feature Specification

### Core Features (MVP)
1. **OpenAPI Import**: File selection and URL input for OpenAPI specs
2. **Spec Validation**: Real-time validation with detailed error reporting
3. **MCP Generation**: One-click generation of complete MCP servers
4. **Preview**: Generated code preview before saving
5. **Export**: Save generated servers to chosen directory
6. **History**: Track previously generated servers

### Advanced Features (Post-MVP)
1. **Template Customization**: Custom server templates
2. **Batch Processing**: Multiple spec processing
3. **Integration Testing**: Built-in MCP server testing
4. **Configuration Profiles**: Reusable generation settings
5. **Plugin System**: Extensible architecture for custom processors

## Technical Architecture

### Application Stack
- **Backend**: Go 1.21+ (reusing openapi2mcp core)
- **Frontend**: Electron with React/TypeScript
- **Database**: SQLite for local project storage
- **Communication**: Go-Electron IPC bridge
- **Packaging**: Electron Builder for cross-platform distribution

### Core Components
- **Parser Engine**: OpenAPI specification parsing (kin-openapi)
- **Mapping Engine**: OpenAPI to MCP tool transformation
- **Generator Engine**: Template-based code generation
- **UI Controller**: Desktop application interface
- **Project Manager**: Local project and history management
- **Validator**: Specification and generated code validation

## Data Architecture

### Local Storage
- **Database**: SQLite for projects, settings, and history
- **File System**: Generated servers and templates
- **Configuration**: JSON-based application settings
- **Logs**: Structured logging for debugging

### Data Models
- **Project**: OpenAPI spec, generation settings, output location
- **Generation**: Job status, progress, results, errors
- **Template**: Custom generation templates and configurations
- **Settings**: User preferences and application configuration

## User Interface Specification

### Main Application Window
- **Menu Bar**: File, Edit, View, Tools, Help
- **Toolbar**: Quick access to common operations
- **Main Panel**: Tabbed interface for projects
- **Status Bar**: Current operation status and progress
- **Side Panel**: Project navigation and tool palette

### Key User Flows
1. **New Project**: OpenAPI spec → Validation → Configuration → Generation
2. **Edit Project**: Modify settings → Re-generate → Compare results
3. **View History**: Browse past generations → Re-open projects
4. **Export Server**: Generation → Testing → Save to directory

## Observability Requirements

### User-Facing Observability
- **Real-time Progress**: Generation progress with detailed steps
- **Error Reporting**: Clear error messages with suggested fixes
- **Performance Metrics**: Generation time, spec complexity metrics
- **Success Indicators**: Visual confirmation of successful generation
- **Resource Usage**: Memory and CPU usage during operations

### Internal Observability
- **Structured Logging**: JSON-formatted logs with correlation IDs
- **Metrics Collection**: Operation timing, error rates, resource usage
- **Health Checks**: Component status and dependency validation
- **Debug Mode**: Detailed tracing for development and troubleshooting

## Performance Requirements

### Generation Performance
- **Small Specs** (< 10 endpoints): < 1 second
- **Medium Specs** (10-100 endpoints): < 3 seconds
- **Large Specs** (100+ endpoints): < 10 seconds
- **Memory Usage**: Linear scaling with spec complexity

### UI Performance
- **Startup Time**: < 2 seconds cold start
- **UI Responsiveness**: < 100ms UI response time
- **File Operations**: < 500ms for file I/O operations
- **Background Processing**: Non-blocking with progress indication

## Security Requirements

### Input Security
- **File Validation**: Strict OpenAPI spec validation
- **Path Traversal Protection**: Secure file system access
- **URL Validation**: Safe remote spec fetching
- **Template Security**: Sandboxed template execution

### Data Security
- **Local Storage**: Encrypted sensitive data
- **File Permissions**: Appropriate file system permissions
- **Process Isolation**: Sandboxed generation processes
- **Update Security**: Signed application updates

## Auto-Update System

### Overview
MCPWeaver includes a comprehensive auto-update system that enables seamless application updates with minimal user intervention while maintaining security and reliability.

### Core Features
- **Automatic Update Detection**: Periodic checks for new releases via GitHub API
- **Secure Download**: HTTPS-only downloads with checksum verification
- **Digital Signature Verification**: Cryptographic validation of update authenticity
- **Rollback Capability**: Automatic backup and rollback for failed updates
- **User Control**: Configurable update preferences and scheduling
- **Progress Tracking**: Real-time download and installation progress
- **Bandwidth Management**: Configurable download speed limits

### Update Flow
1. **Detection**: Background check for newer versions on GitHub releases
2. **Notification**: User notification with release notes and options
3. **Download**: Secure download of new version with progress tracking
4. **Verification**: Checksum and signature validation
5. **Backup**: Create backup of current version before update
6. **Installation**: Replace current executable with new version
7. **Restart**: Application restart to complete update
8. **Rollback**: Automatic rollback if update fails

### Security Model
- **HTTPS Only**: All update communications over encrypted channels
- **Checksum Verification**: SHA256 verification of downloaded files
- **Digital Signatures**: RSA signature verification for authenticity
- **Backup Protection**: Secure backup storage with integrity checks
- **Sandboxed Operations**: Isolated update processes
- **User Consent**: No automatic installation without user approval

### Update Scheduling
- **Immediate**: Install updates as soon as available
- **Daily**: Check and install at specified time daily
- **Weekly**: Check and install on specified day/time weekly
- **Monthly**: Check and install on specified date monthly
- **Manual**: User-initiated updates only

### Configuration Options
- **Auto-Check**: Enable/disable automatic update checking
- **Auto-Download**: Enable/disable automatic update downloading
- **Auto-Install**: Enable/disable automatic installation
- **Update Channel**: Stable, beta, alpha release channels
- **Pre-Release**: Include pre-release versions
- **Bandwidth Limit**: Maximum download speed in bytes/second
- **Retry Policy**: Configurable retry attempts and delays

### Rollback System
- **Automatic Backup**: Current version backed up before each update
- **Multiple Versions**: Keep last 5 versions for rollback options
- **Validation**: Backup integrity verification before rollback
- **User Interface**: Easy rollback through application settings
- **Emergency Recovery**: Command-line rollback for critical failures

### Analytics and Monitoring
- **Update Metrics**: Track update success/failure rates
- **Performance Data**: Download speeds and installation times
- **Error Reporting**: Detailed error logs for troubleshooting
- **User Behavior**: Anonymous update preference analytics
- **System Information**: Platform and version compatibility data

### Offline Support
- **Local Updates**: Support for manual update file installation
- **Update Packages**: Downloadable update packages for offline installation
- **Network Detection**: Automatic handling of network availability
- **Deferred Updates**: Queue updates when offline, install when online

### Technical Implementation
- **GitHub Integration**: GitHub Releases API for version checking
- **HTTP Client**: Configurable HTTP client with timeout and retry
- **File Operations**: Secure file handling with proper permissions
- **Process Management**: Safe process restart and cleanup
- **Event System**: Real-time update events for UI integration

### Error Handling
- **Network Errors**: Graceful handling of connectivity issues
- **Download Failures**: Automatic retry with exponential backoff
- **Verification Failures**: Secure handling of invalid updates
- **Installation Errors**: Automatic rollback and error reporting
- **Rollback Failures**: Emergency recovery procedures

### User Experience
- **Non-Intrusive**: Updates don't interrupt active work
- **Progress Visibility**: Clear progress indication throughout process
- **User Control**: Full control over update timing and preferences
- **Notifications**: Informative but not overwhelming notifications
- **Settings Integration**: Update preferences in main application settings

## Distribution Strategy

### Release Channels
- **GitHub Releases**: Primary distribution channel
- **Package Managers**: brew, winget, apt (future)
- **Direct Download**: Platform-specific installers
- **Auto-Updates**: Secure in-app update mechanism

### Packaging Requirements
- **Single Binary**: Self-contained executable
- **Cross-Platform**: Windows (.exe), macOS (.dmg), Linux (.AppImage)
- **Code Signing**: Signed binaries for security
- **Minimal Dependencies**: No external runtime requirements

## Development Guidelines

### Code Quality
- **Test Coverage**: > 90% unit test coverage
- **Documentation**: Comprehensive API and user documentation
- **Code Review**: All changes require peer review
- **Linting**: Automated code quality checks

### Development Process
- **Spec-Driven Development**: Specifications before implementation
- **Iterative Development**: Regular MVP iterations
- **User Feedback**: Continuous user testing and feedback
- **Performance Testing**: Regular performance benchmarking

## Success Metrics

### Technical Metrics
- **Generation Success Rate**: > 95% for valid OpenAPI specs
- **Performance Benchmarks**: Meet all performance requirements
- **Error Recovery**: < 1% unrecoverable errors
- **Resource Efficiency**: Memory usage within specifications

### User Experience Metrics
- **Time to First Success**: < 2 minutes from install to first generation
- **User Satisfaction**: Measured through feedback and usage patterns
- **Adoption Rate**: GitHub stars, downloads, community engagement
- **Issue Resolution**: < 48 hours for critical issues

## Timeline and Milestones

### Phase 1: Foundation (Weeks 1-2)
- Architecture design and specifications
- Core Go backend extraction from openapi2mcp
- Basic Electron application shell
- Local SQLite database setup

### Phase 2: Core MVP (Weeks 3-4)
- OpenAPI parsing and validation
- MCP server generation
- Basic UI for import/export
- Real-time progress tracking

### Phase 3: Polish and Testing (Weeks 5-6)
- Error handling and recovery
- Performance optimization
- Comprehensive testing
- User experience improvements

### Phase 4: Release Preparation (Weeks 7-8)
- Cross-platform packaging
- Documentation and guides
- Security audit
- Initial release preparation

This specification provides a comprehensive foundation for building MCPWeaver as a focused, efficient desktop application that delivers maximum value with minimal complexity.