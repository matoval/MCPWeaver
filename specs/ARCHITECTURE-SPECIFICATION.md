# MCPWeaver Desktop Application Architecture Specification

## Architecture Overview

MCPWeaver follows a native desktop application architecture using Wails v2, with a Go backend and modern web frontend, providing native performance with web UI flexibility.

```
┌─────────────────────────────────────────────────────────────┐
│                    MCPWeaver Desktop App                    │
├─────────────────────────────────────────────────────────────┤
│  Frontend (React/TypeScript)                               │
│  ├─ React Components (UI Layer)                           │
│  ├─ State Management (Context/Zustand)                    │
│  └─ Wails Runtime API                                     │
├─────────────────────────────────────────────────────────────┤
│  Wails Runtime                                             │
│  ├─ Context Bridge                                        │
│  ├─ Event System                                          │
│  └─ Native OS Integration                                 │
├─────────────────────────────────────────────────────────────┤
│  Backend Services (Go)                                     │
│  ├─ App Context (Main Application)                        │
│  ├─ Parser Service (OpenAPI Processing)                   │
│  ├─ Generator Service (MCP Server Generation)             │
│  ├─ Validator Service (Spec & Code Validation)            │
│  └─ Project Manager (Local Project Management)            │
├─────────────────────────────────────────────────────────────┤
│  Data Layer                                                │
│  ├─ SQLite Database (Projects, History, Settings)         │
│  ├─ File System (Generated Servers, Templates)            │
│  └─ Configuration (JSON Settings)                         │
└─────────────────────────────────────────────────────────────┘
```

## Application Structure

### Directory Layout
```
MCPWeaver/
├── app.go                # Main Wails application
├── wails.json           # Wails configuration
├── build/               # Build configuration
├── internal/            # Go backend services
│   ├── parser/          # OpenAPI parsing (from openapi2mcp)
│   ├── generator/       # MCP server generation
│   ├── validator/       # Validation services
│   ├── project/         # Project management
│   └── database/        # SQLite database layer
├── frontend/            # React frontend
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── services/    # API service layer
│   │   ├── stores/      # State management
│   │   └── types/       # TypeScript types
│   ├── public/          # Static assets
│   └── package.json     # Frontend dependencies
├── templates/           # MCP server generation templates
└── assets/              # Application assets and icons
```

## Component Architecture

### 1. Wails Application Entry Point

```go
// app.go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/wailsapp/wails/v2"
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// App struct
type App struct {
    ctx         context.Context
    parser      *parser.Service
    generator   *generator.Service
    validator   *validator.Service
    project     *project.Service
    database    *database.DB
}

// NewApp creates a new App application struct
func NewApp() *App {
    return &App{}
}

// OnStartup is called when the app starts
func (a *App) OnStartup(ctx context.Context) {
    a.ctx = ctx
    
    // Initialize database
    db, err := database.Open("mcpweaver.db")
    if err != nil {
        log.Fatal("Failed to open database:", err)
    }
    a.database = db
    
    // Initialize services
    a.parser = parser.New()
    a.generator = generator.New()
    a.validator = validator.New()
    a.project = project.New(db)
}

// OnShutdown is called when the app is shutting down
func (a *App) OnShutdown(ctx context.Context) {
    if a.database != nil {
        a.database.Close()
    }
}

func main() {
    // Create an instance of the app structure
    app := NewApp()

    // Create application with options
    err := wails.Run(&options.App{
        Title:  "MCPWeaver",
        Width:  1200,
        Height: 800,
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        OnStartup:  app.OnStartup,
        OnShutdown: app.OnShutdown,
    })

    if err != nil {
        println("Error:", err.Error())
    }
}
```

### 2. Frontend Layer (React/TypeScript)

#### Main App Component
```typescript
// frontend/src/App.tsx
import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ProjectProvider } from './contexts/ProjectContext';
import { MainLayout } from './components/layouts/MainLayout';
import { ProjectList } from './components/ProjectList';
import { ProjectEditor } from './components/ProjectEditor';
import { GenerationView } from './components/GenerationView';

const App: React.FC = () => {
  return (
    <ProjectProvider>
      <Router>
        <MainLayout>
          <Routes>
            <Route path="/" element={<ProjectList />} />
            <Route path="/project/:id" element={<ProjectEditor />} />
            <Route path="/generate/:id" element={<GenerationView />} />
          </Routes>
        </MainLayout>
      </Router>
    </ProjectProvider>
  );
};

export default App;
```

#### Wails API Service Layer
```typescript
// frontend/src/services/wailsApi.ts
import { 
  CreateProject, 
  GetProjects, 
  GenerateServer, 
  ValidateSpec 
} from '../wailsjs/go/main/App';

export class WailsApiService {
  // Project Management
  async createProject(request: CreateProjectRequest): Promise<Project> {
    try {
      const project = await CreateProject(request);
      return project;
    } catch (error) {
      throw new Error(`Failed to create project: ${error}`);
    }
  }

  async getProjects(): Promise<Project[]> {
    try {
      const projects = await GetProjects();
      return projects;
    } catch (error) {
      throw new Error(`Failed to get projects: ${error}`);
    }
  }

  // Generation
  async generateServer(projectId: string): Promise<GenerationJob> {
    try {
      const job = await GenerateServer(projectId);
      return job;
    } catch (error) {
      throw new Error(`Failed to start generation: ${error}`);
    }
  }

  // Validation
  async validateSpec(specPath: string): Promise<ValidationResult> {
    try {
      const result = await ValidateSpec(specPath);
      return result;
    } catch (error) {
      throw new Error(`Failed to validate spec: ${error}`);
    }
  }
}

export const wailsApi = new WailsApiService();
```

#### React Context for State Management
```typescript
// frontend/src/contexts/ProjectContext.tsx
import React, { createContext, useContext, useReducer, useEffect } from 'react';
import { wailsApi } from '../services/wailsApi';
import { EventsOn } from '../wailsjs/runtime/runtime';

interface ProjectState {
  projects: Project[];
  activeProject: Project | null;
  isLoading: boolean;
  error: string | null;
}

const ProjectContext = createContext<ProjectState | undefined>(undefined);

export const ProjectProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [state, dispatch] = useReducer(projectReducer, initialState);

  useEffect(() => {
    // Load initial projects
    loadProjects();

    // Set up event listeners for real-time updates
    EventsOn('project:created', (project: Project) => {
      dispatch({ type: 'PROJECT_CREATED', payload: project });
    });

    EventsOn('project:updated', (project: Project) => {
      dispatch({ type: 'PROJECT_UPDATED', payload: project });
    });

    EventsOn('generation:progress', (progress: GenerationProgress) => {
      dispatch({ type: 'GENERATION_PROGRESS', payload: progress });
    });

    return () => {
      // Cleanup event listeners
      EventsOff('project:created');
      EventsOff('project:updated');
      EventsOff('generation:progress');
    };
  }, []);

  const loadProjects = async () => {
    try {
      dispatch({ type: 'LOADING_START' });
      const projects = await wailsApi.getProjects();
      dispatch({ type: 'PROJECTS_LOADED', payload: projects });
    } catch (error) {
      dispatch({ type: 'ERROR', payload: error.message });
    }
  };

  return (
    <ProjectContext.Provider value={state}>
      {children}
    </ProjectContext.Provider>
  );
};
```

### 3. Backend Services (Go)

#### Wails-Compatible Service Methods
```go
// internal/app/project_methods.go
package main

import (
    "context"
    "fmt"
    "github.com/wailsapp/wails/v2/pkg/runtime"
)

// CreateProject creates a new project
func (a *App) CreateProject(request CreateProjectRequest) (*Project, error) {
    project, err := a.project.Create(request)
    if err != nil {
        return nil, fmt.Errorf("failed to create project: %w", err)
    }

    // Emit event to frontend
    runtime.EventsEmit(a.ctx, "project:created", project)
    
    return project, nil
}

// GetProjects returns all projects
func (a *App) GetProjects() ([]*Project, error) {
    projects, err := a.project.GetAll()
    if err != nil {
        return nil, fmt.Errorf("failed to get projects: %w", err)
    }
    
    return projects, nil
}

// GenerateServer starts MCP server generation
func (a *App) GenerateServer(projectId string) (*GenerationJob, error) {
    job, err := a.generator.StartGeneration(a.ctx, projectId)
    if err != nil {
        return nil, fmt.Errorf("failed to start generation: %w", err)
    }

    // Set up progress monitoring
    go a.monitorGenerationProgress(job)
    
    return job, nil
}

// ValidateSpec validates an OpenAPI specification
func (a *App) ValidateSpec(specPath string) (*ValidationResult, error) {
    result, err := a.validator.ValidateFile(specPath)
    if err != nil {
        return nil, fmt.Errorf("failed to validate spec: %w", err)
    }
    
    return result, nil
}

// SelectFile opens a file dialog
func (a *App) SelectFile(filters []FileFilter) (string, error) {
    options := runtime.OpenDialogOptions{
        Title:   "Select OpenAPI Specification",
        Filters: convertFilters(filters),
    }
    
    filePath, err := runtime.OpenFileDialog(a.ctx, options)
    if err != nil {
        return "", fmt.Errorf("failed to open file dialog: %w", err)
    }
    
    return filePath, nil
}

// monitorGenerationProgress monitors and emits progress updates
func (a *App) monitorGenerationProgress(job *GenerationJob) {
    for progress := range job.ProgressChan {
        runtime.EventsEmit(a.ctx, "generation:progress", progress)
    }
    
    // Emit completion event
    runtime.EventsEmit(a.ctx, "generation:completed", job)
}
```

#### Generator Service with Progress Tracking
```go
// internal/generator/service.go
package generator

import (
    "context"
    "fmt"
    "time"
)

type Service struct {
    templateManager *TemplateManager
}

type GenerationJob struct {
    ID           string                `json:"id"`
    ProjectID    string                `json:"projectId"`
    Status       GenerationStatus      `json:"status"`
    Progress     float64               `json:"progress"`
    CurrentStep  string                `json:"currentStep"`
    StartTime    time.Time             `json:"startTime"`
    EndTime      *time.Time            `json:"endTime"`
    Results      *GenerationResults    `json:"results"`
    Errors       []GenerationError     `json:"errors"`
    ProgressChan chan GenerationProgress `json:"-"`
}

type GenerationProgress struct {
    JobID       string  `json:"jobId"`
    Progress    float64 `json:"progress"`
    CurrentStep string  `json:"currentStep"`
    Message     string  `json:"message"`
}

func (s *Service) StartGeneration(ctx context.Context, projectId string) (*GenerationJob, error) {
    job := &GenerationJob{
        ID:           generateID(),
        ProjectID:    projectId,
        Status:       StatusStarted,
        StartTime:    time.Now(),
        ProgressChan: make(chan GenerationProgress, 10),
    }

    go s.processGeneration(ctx, job)
    return job, nil
}

func (s *Service) processGeneration(ctx context.Context, job *GenerationJob) {
    defer close(job.ProgressChan)
    
    steps := []struct {
        progress float64
        message  string
        action   func() error
    }{
        {0.25, "Parsing OpenAPI specification...", s.parseSpec},
        {0.50, "Mapping endpoints to MCP tools...", s.mapEndpoints},
        {0.75, "Generating MCP server code...", s.generateCode},
        {1.00, "Validating generated server...", s.validateGenerated},
    }

    for _, step := range steps {
        select {
        case <-ctx.Done():
            job.Status = StatusCancelled
            return
        default:
            job.Progress = step.progress
            job.CurrentStep = step.message
            
            job.ProgressChan <- GenerationProgress{
                JobID:       job.ID,
                Progress:    step.progress,
                CurrentStep: step.message,
                Message:     step.message,
            }

            if err := step.action(); err != nil {
                job.Status = StatusFailed
                job.Errors = append(job.Errors, GenerationError{
                    Type:    "generation",
                    Message: err.Error(),
                })
                return
            }
        }
    }

    job.Status = StatusCompleted
    now := time.Now()
    job.EndTime = &now
}
```

### 4. Wails Configuration

```json
// wails.json
{
  "name": "MCPWeaver",
  "outputfilename": "mcpweaver",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev": "npm run dev",
  "frontend:dev:serverUrl": "http://localhost:3000",
  "author": {
    "name": "MCPWeaver Team",
    "email": "team@mcpweaver.dev"
  },
  "info": {
    "productName": "MCPWeaver",
    "productVersion": "1.0.0",
    "copyright": "Copyright © 2024 MCPWeaver Team",
    "companyName": "MCPWeaver"
  },
  "nsisType": "multiple",
  "obfuscated": false,
  "garbleargs": ""
}
```

### 5. Build Configuration

```json
// build/windows/info.json
{
  "fixed": {
    "file_version": "1.0.0",
    "product_version": "1.0.0"
  },
  "info": {
    "0000": {
      "ProductName": "MCPWeaver",
      "CompanyName": "MCPWeaver Team",
      "FileDescription": "OpenAPI to MCP Server Generator",
      "InternalName": "mcpweaver",
      "LegalCopyright": "Copyright © 2024 MCPWeaver Team",
      "OriginalFilename": "mcpweaver.exe",
      "ProductVersion": "1.0.0"
    }
  }
}
```

## Performance Considerations

### Native Performance Benefits
- **Single Binary**: No Node.js runtime overhead
- **Native OS Integration**: Direct system API access
- **Memory Efficiency**: Go's garbage collector and efficient memory management
- **Fast Startup**: Native application startup times
- **Cross-Platform**: Single codebase for all platforms

### Wails-Specific Optimizations
- **Context Binding**: Efficient data transfer between Go and frontend
- **Event System**: Real-time updates without polling
- **Resource Embedding**: Assets embedded in binary
- **Hot Reload**: Fast development iteration

## Security Architecture

### Wails Security Features
- **Native Process**: No browser security model limitations
- **Direct File Access**: Native file system permissions
- **OS Integration**: System-level security controls
- **Code Signing**: Native application signing support

### Application Security
- **Input Validation**: Comprehensive validation of all inputs
- **File System Security**: Restricted file access patterns
- **Template Security**: Sandboxed template execution
- **Update Security**: Signed application updates

## Development and Build Process

### Development Setup
```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Initialize project
wails init -n MCPWeaver -t vanilla

# Development with hot reload
wails dev

# Build for production
wails build
```

### Cross-Platform Building
```bash
# Build for all platforms
wails build -platform windows/amd64,darwin/amd64,linux/amd64

# Build for specific platform
wails build -platform windows/amd64
```

This Wails-based architecture provides native performance, simplified deployment, and seamless integration between Go backend services and modern web frontend, making it ideal for MCPWeaver's requirements.