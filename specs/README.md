# MCPWeaver - Specification Documents

## Overview

This repository contains comprehensive specification documents for MCPWeaver, an open-source desktop application that transforms OpenAPI specifications into Model Context Protocol (MCP) servers.

## Project Vision

MCPWeaver represents a pivot from the hosted SaaS openapi2mcp to a lightweight, user-controlled desktop application that runs entirely on the user's machine, built with Go backend and Wails framework.

## Specification Documents

### 1. [Project Specification](PROJECT-SPECIFICATION.md)
- **Purpose**: Overall project vision, requirements, and goals
- **Key Sections**:
  - Project vision and core principles
  - Technical requirements and platform support
  - Feature specification (MVP and advanced)
  - Performance and security requirements
  - Success metrics and timeline

### 2. [Architecture Specification](ARCHITECTURE-SPECIFICATION.md)
- **Purpose**: Technical architecture and system design
- **Key Sections**:
  - Wails-based desktop application architecture
  - Component structure and communication patterns
  - Backend services (Parser, Generator, Validator, Project Manager)
  - Database schema and data layer
  - Error handling and recovery strategies

### 3. [API Specification](API-SPECIFICATION.md)
- **Purpose**: Internal API contracts and communication interfaces
- **Key Sections**:
  - Wails context API methods
  - Data structures and type definitions
  - Event system for real-time updates
  - Error handling and validation
  - Rate limiting and concurrency management

### 4. [Data Models Specification](DATA-MODELS-SPECIFICATION.md)
- **Purpose**: Data structures, database schemas, and models
- **Key Sections**:
  - SQLite database schema
  - Core data models (Project, Generation, Template, etc.)
  - OpenAPI and MCP model definitions
  - Validation rules and constraints
  - Data relationships and constraints

### 5. [Observability Specification](OBSERVABILITY-SPECIFICATION.md)
- **Purpose**: Monitoring, logging, and user-facing observability
- **Key Sections**:
  - User-centric observability principles
  - Real-time status and progress tracking
  - Performance metrics and error reporting
  - Structured logging and health checks
  - Lightweight monitoring approach

### 6. [User Interface Specification](UI-SPECIFICATION.md)
- **Purpose**: UI design, layout, and interaction patterns
- **Key Sections**:
  - Application layout and component structure
  - Menu system and toolbar design
  - Theme system and responsive design
  - Accessibility and keyboard navigation
  - Animation and performance considerations

## Key Design Principles

### Technical Excellence
- **Wails Framework**: Native performance with web UI flexibility
- **Go Backend**: Reusing proven openapi2mcp core functionality
- **SQLite Database**: Local storage with no external dependencies
- **Single Binary**: Self-contained distribution

### User Experience
- **Simplicity First**: Minimal cognitive load and clear hierarchy
- **Efficiency Focus**: Fast access to common operations
- **Real-time Feedback**: Immediate visual feedback for all actions
- **Offline Operation**: No cloud dependencies

### Development Approach
- **Spec-Driven Development**: Implementation follows detailed specifications
- **Iterative Development**: Regular MVP iterations with user feedback
- **High Code Quality**: >90% test coverage and comprehensive documentation
- **Open Source**: MIT license with community contribution support

## Development Phases

### Phase 1: Foundation (Weeks 1-2)
- Core Go backend extraction from openapi2mcp
- Basic Wails application shell
- Local SQLite database setup
- Architecture implementation

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

## Technology Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Wails v2
- **Database**: SQLite with structured schema
- **Libraries**: kin-openapi (reused from openapi2mcp)

### Frontend
- **Framework**: React with TypeScript
- **State Management**: Context API / Zustand
- **UI Components**: Custom component library
- **Styling**: SCSS with CSS custom properties
- **Icons**: Lucide React

### Build and Distribution
- **Build Tool**: Wails CLI
- **Packaging**: Cross-platform (Windows, macOS, Linux)
- **Distribution**: GitHub Releases
- **Updates**: Secure in-app update mechanism

## Getting Started

These specifications provide the foundation for implementing MCPWeaver. To begin development:

1. **Review all specification documents** to understand the complete system
2. **Set up development environment** with Go 1.21+ and Wails v2
3. **Extract core functionality** from existing openapi2mcp codebase
4. **Implement architecture** following the specifications
5. **Build iteratively** starting with Phase 1 components

## Contributing

This project will be open source under MIT license. Contributions are welcome following the specifications and maintaining the core principles of simplicity and efficiency.

## Next Steps

1. **Create new repository** for MCPWeaver implementation
2. **Initialize Wails project** with proper structure
3. **Extract and adapt** core functionality from openapi2mcp
4. **Implement specifications** in order of priority
5. **Test thoroughly** with comprehensive test suite
6. **Document extensively** for users and contributors

---

*These specifications represent a comprehensive foundation for building MCPWeaver as a focused, efficient desktop application that delivers maximum value with minimal complexity.*