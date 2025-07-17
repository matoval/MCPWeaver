# MCPWeaver User Interface Specification

## Overview

This document defines the user interface design, layout, and interaction patterns for MCPWeaver desktop application, focusing on simplicity, efficiency, and user-friendly experience.

## UI Design Principles

### Core Design Philosophy
- **Simplicity First**: Minimal cognitive load, clear hierarchy
- **Efficiency Focus**: Fast access to common operations
- **Progressive Disclosure**: Advanced features available but not overwhelming
- **Consistent Interaction**: Familiar patterns throughout the application
- **Responsive Feedback**: Immediate visual feedback for all actions

### Visual Design Language
- **Clean and Modern**: Flat design with subtle shadows and gradients
- **Contextual Color**: Color used meaningfully to convey status and actions
- **Typography**: Clear, readable fonts with proper hierarchy
- **Spacing**: Generous whitespace for better readability
- **Icons**: Consistent icon system for navigation and actions

## Application Layout

### Main Window Structure
```
┌─────────────────────────────────────────────────────────────┐
│ Menu Bar                                                    │
├─────────────────────────────────────────────────────────────┤
│ Toolbar                                                     │
├─────────────────────────────────────────────────────────────┤
│ ┌─────────────┐ ┌─────────────────────────────────────────┐ │
│ │             │ │                                         │ │
│ │   Sidebar   │ │           Main Content Area             │ │
│ │             │ │                                         │ │
│ │             │ │                                         │ │
│ └─────────────┘ └─────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│ Status Bar                                                  │
└─────────────────────────────────────────────────────────────┘
```

### Window Dimensions
- **Minimum Size**: 1024x768 pixels
- **Default Size**: 1200x800 pixels
- **Resizable**: Yes, with minimum constraints
- **Maximizable**: Yes
- **Multi-monitor Support**: Yes

## Component Specifications

### Menu Bar
```typescript
interface MenuStructure {
  file: {
    label: "File";
    items: [
      { label: "New Project", shortcut: "Ctrl+N", action: "newProject" },
      { label: "Open Project", shortcut: "Ctrl+O", action: "openProject" },
      { label: "Recent Projects", submenu: RecentProjectsMenu },
      { separator: true },
      { label: "Import OpenAPI Spec", shortcut: "Ctrl+I", action: "importSpec" },
      { label: "Export Server", shortcut: "Ctrl+E", action: "exportServer" },
      { separator: true },
      { label: "Settings", shortcut: "Ctrl+,", action: "openSettings" },
      { separator: true },
      { label: "Exit", shortcut: "Ctrl+Q", action: "exit" }
    ];
  };
  edit: {
    label: "Edit";
    items: [
      { label: "Undo", shortcut: "Ctrl+Z", action: "undo" },
      { label: "Redo", shortcut: "Ctrl+Y", action: "redo" },
      { separator: true },
      { label: "Cut", shortcut: "Ctrl+X", action: "cut" },
      { label: "Copy", shortcut: "Ctrl+C", action: "copy" },
      { label: "Paste", shortcut: "Ctrl+V", action: "paste" },
      { separator: true },
      { label: "Find", shortcut: "Ctrl+F", action: "find" },
      { label: "Replace", shortcut: "Ctrl+H", action: "replace" }
    ];
  };
  view: {
    label: "View";
    items: [
      { label: "Zoom In", shortcut: "Ctrl++", action: "zoomIn" },
      { label: "Zoom Out", shortcut: "Ctrl+-", action: "zoomOut" },
      { label: "Reset Zoom", shortcut: "Ctrl+0", action: "resetZoom" },
      { separator: true },
      { label: "Toggle Sidebar", shortcut: "Ctrl+B", action: "toggleSidebar" },
      { label: "Toggle Status Bar", action: "toggleStatusBar" },
      { separator: true },
      { label: "Activity Log", shortcut: "Ctrl+L", action: "showActivityLog" },
      { label: "Performance Metrics", action: "showMetrics" }
    ];
  };
  tools: {
    label: "Tools";
    items: [
      { label: "Validate Spec", shortcut: "F5", action: "validateSpec" },
      { label: "Generate Server", shortcut: "F6", action: "generateServer" },
      { label: "Test Server", shortcut: "F7", action: "testServer" },
      { separator: true },
      { label: "Template Manager", action: "manageTemplates" },
      { label: "Clear Cache", action: "clearCache" }
    ];
  };
  help: {
    label: "Help";
    items: [
      { label: "User Guide", shortcut: "F1", action: "showUserGuide" },
      { label: "API Documentation", action: "showApiDocs" },
      { label: "Keyboard Shortcuts", action: "showShortcuts" },
      { separator: true },
      { label: "Report Issue", action: "reportIssue" },
      { label: "About MCPWeaver", action: "showAbout" }
    ];
  };
}
```

### Toolbar
```typescript
interface ToolbarButton {
  id: string;
  label: string;
  icon: string;
  tooltip: string;
  action: string;
  shortcut?: string;
  enabled: boolean;
  variant: 'primary' | 'secondary' | 'danger' | 'success';
}

const toolbarButtons: ToolbarButton[] = [
  {
    id: 'new-project',
    label: 'New Project',
    icon: 'plus',
    tooltip: 'Create a new project (Ctrl+N)',
    action: 'newProject',
    shortcut: 'Ctrl+N',
    enabled: true,
    variant: 'primary'
  },
  {
    id: 'open-project',
    label: 'Open',
    icon: 'folder-open',
    tooltip: 'Open existing project (Ctrl+O)',
    action: 'openProject',
    shortcut: 'Ctrl+O',
    enabled: true,
    variant: 'secondary'
  },
  {
    id: 'import-spec',
    label: 'Import',
    icon: 'download',
    tooltip: 'Import OpenAPI specification (Ctrl+I)',
    action: 'importSpec',
    shortcut: 'Ctrl+I',
    enabled: true,
    variant: 'secondary'
  },
  {
    id: 'validate-spec',
    label: 'Validate',
    icon: 'check-circle',
    tooltip: 'Validate OpenAPI specification (F5)',
    action: 'validateSpec',
    shortcut: 'F5',
    enabled: false, // Enabled when project is selected
    variant: 'secondary'
  },
  {
    id: 'generate-server',
    label: 'Generate',
    icon: 'code',
    tooltip: 'Generate MCP server (F6)',
    action: 'generateServer',
    shortcut: 'F6',
    enabled: false, // Enabled when project is validated
    variant: 'success'
  },
  {
    id: 'export-server',
    label: 'Export',
    icon: 'upload',
    tooltip: 'Export generated server (Ctrl+E)',
    action: 'exportServer',
    shortcut: 'Ctrl+E',
    enabled: false, // Enabled when generation is complete
    variant: 'primary'
  }
];
```

### Sidebar
```typescript
interface SidebarSection {
  id: string;
  title: string;
  icon: string;
  collapsible: boolean;
  defaultExpanded: boolean;
  content: React.ComponentType;
}

const sidebarSections: SidebarSection[] = [
  {
    id: 'projects',
    title: 'Projects',
    icon: 'folder',
    collapsible: true,
    defaultExpanded: true,
    content: ProjectList
  },
  {
    id: 'recent',
    title: 'Recent',
    icon: 'clock',
    collapsible: true,
    defaultExpanded: false,
    content: RecentProjectsList
  },
  {
    id: 'templates',
    title: 'Templates',
    icon: 'template',
    collapsible: true,
    defaultExpanded: false,
    content: TemplateList
  },
  {
    id: 'activity',
    title: 'Activity',
    icon: 'activity',
    collapsible: true,
    defaultExpanded: false,
    content: ActivityLog
  }
];
```

### Main Content Area

#### Project Dashboard
```typescript
interface ProjectDashboard {
  project: Project;
  validation: ValidationResult;
  generationHistory: GenerationJob[];
  metrics: ProjectMetrics;
}

const ProjectDashboard: React.FC<ProjectDashboard> = ({ project, validation, generationHistory, metrics }) => {
  return (
    <div className="project-dashboard">
      <ProjectHeader project={project} />
      <div className="dashboard-grid">
        <ProjectInfoCard project={project} />
        <ValidationStatusCard validation={validation} />
        <GenerationHistoryCard history={generationHistory} />
        <MetricsCard metrics={metrics} />
      </div>
      <ProjectActions project={project} />
    </div>
  );
};
```

#### Generation Progress View
```typescript
interface GenerationProgressView {
  job: GenerationJob;
  onCancel: () => void;
  onShowDetails: () => void;
}

const GenerationProgressView: React.FC<GenerationProgressView> = ({ job, onCancel, onShowDetails }) => {
  return (
    <div className="generation-progress">
      <ProgressHeader job={job} />
      <ProgressIndicator progress={job.progress} />
      <StepDetails currentStep={job.currentStep} />
      <ProgressMetrics metrics={job.metrics} />
      <ActionButtons onCancel={onCancel} onShowDetails={onShowDetails} />
    </div>
  );
};
```

#### Code Preview
```typescript
interface CodePreviewProps {
  files: GeneratedFile[];
  selectedFile: string;
  onFileSelect: (file: string) => void;
  onSave: () => void;
}

const CodePreview: React.FC<CodePreviewProps> = ({ files, selectedFile, onFileSelect, onSave }) => {
  return (
    <div className="code-preview">
      <FileTree files={files} selectedFile={selectedFile} onFileSelect={onFileSelect} />
      <CodeEditor file={selectedFile} readOnly={true} />
      <PreviewActions onSave={onSave} />
    </div>
  );
};
```

### Status Bar
```typescript
interface StatusBarProps {
  status: ApplicationStatus;
  activeOperations: number;
  systemHealth: SystemHealth;
  onStatusClick: () => void;
}

const StatusBar: React.FC<StatusBarProps> = ({ status, activeOperations, systemHealth, onStatusClick }) => {
  return (
    <div className="status-bar">
      <StatusIndicator status={status} onClick={onStatusClick} />
      <OperationCounter count={activeOperations} />
      <SystemHealthIndicator health={systemHealth} />
      <div className="status-spacer" />
      <ResourceUsage memory={systemHealth.memoryUsage} cpu={systemHealth.cpuUsage} />
      <AppVersion />
    </div>
  );
};
```

## Modal Dialogs and Overlays

### New Project Dialog
```typescript
interface NewProjectDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (project: CreateProjectRequest) => void;
}

const NewProjectDialog: React.FC<NewProjectDialogProps> = ({ isOpen, onClose, onSubmit }) => {
  return (
    <Dialog isOpen={isOpen} onClose={onClose} title="Create New Project">
      <ProjectForm onSubmit={onSubmit} />
    </Dialog>
  );
};
```

### Settings Dialog
```typescript
interface SettingsDialogProps {
  isOpen: boolean;
  onClose: () => void;
  settings: AppSettings;
  onSave: (settings: AppSettings) => void;
}

const SettingsDialog: React.FC<SettingsDialogProps> = ({ isOpen, onClose, settings, onSave }) => {
  return (
    <Dialog isOpen={isOpen} onClose={onClose} title="Settings" size="large">
      <SettingsTabs>
        <SettingsTab id="general" title="General" icon="settings">
          <GeneralSettings settings={settings.general} />
        </SettingsTab>
        <SettingsTab id="generation" title="Generation" icon="code">
          <GenerationSettings settings={settings.generation} />
        </SettingsTab>
        <SettingsTab id="editor" title="Editor" icon="edit">
          <EditorSettings settings={settings.editor} />
        </SettingsTab>
        <SettingsTab id="appearance" title="Appearance" icon="palette">
          <AppearanceSettings settings={settings.appearance} />
        </SettingsTab>
      </SettingsTabs>
    </Dialog>
  );
};
```

## Responsive Design

### Breakpoints
```scss
$breakpoints: (
  mobile: 768px,
  tablet: 1024px,
  desktop: 1200px,
  large: 1440px
);

@media (max-width: 768px) {
  .sidebar {
    transform: translateX(-100%);
    &.open {
      transform: translateX(0);
    }
  }
  
  .main-content {
    margin-left: 0;
  }
}

@media (max-width: 1024px) {
  .toolbar {
    .toolbar-button {
      .button-label {
        display: none;
      }
    }
  }
}
```

### Layout Adaptations
- **< 768px**: Mobile-first layout with collapsible sidebar
- **768px - 1024px**: Tablet layout with icon-only toolbar
- **1024px - 1200px**: Standard desktop layout
- **> 1200px**: Large desktop with expanded panels

## Theme System

### Color Palette
```scss
:root {
  // Primary Colors
  --primary-50: #f0f9ff;
  --primary-100: #e0f2fe;
  --primary-500: #0ea5e9;
  --primary-600: #0284c7;
  --primary-700: #0369a1;
  
  // Secondary Colors
  --secondary-50: #f8fafc;
  --secondary-100: #f1f5f9;
  --secondary-500: #64748b;
  --secondary-600: #475569;
  --secondary-700: #334155;
  
  // Status Colors
  --success-500: #10b981;
  --warning-500: #f59e0b;
  --error-500: #ef4444;
  --info-500: #3b82f6;
  
  // Neutral Colors
  --gray-50: #f9fafb;
  --gray-100: #f3f4f6;
  --gray-200: #e5e7eb;
  --gray-300: #d1d5db;
  --gray-400: #9ca3af;
  --gray-500: #6b7280;
  --gray-600: #4b5563;
  --gray-700: #374151;
  --gray-800: #1f2937;
  --gray-900: #111827;
}
```

### Dark Theme
```scss
[data-theme="dark"] {
  --background: var(--gray-900);
  --foreground: var(--gray-100);
  --surface: var(--gray-800);
  --surface-hover: var(--gray-700);
  --border: var(--gray-700);
  --text-primary: var(--gray-100);
  --text-secondary: var(--gray-400);
  --text-muted: var(--gray-500);
}
```

### Light Theme
```scss
[data-theme="light"] {
  --background: var(--gray-50);
  --foreground: var(--gray-900);
  --surface: #ffffff;
  --surface-hover: var(--gray-100);
  --border: var(--gray-200);
  --text-primary: var(--gray-900);
  --text-secondary: var(--gray-600);
  --text-muted: var(--gray-500);
}
```

## Typography System

### Font Hierarchy
```scss
// Font Families
$font-family-primary: "Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
$font-family-mono: "JetBrains Mono", "Fira Code", Consolas, "Liberation Mono", Monaco, monospace;

// Font Weights
$font-weight-normal: 400;
$font-weight-medium: 500;
$font-weight-semibold: 600;
$font-weight-bold: 700;

// Font Sizes
$font-size-xs: 0.75rem;   // 12px
$font-size-sm: 0.875rem;  // 14px
$font-size-base: 1rem;    // 16px
$font-size-lg: 1.125rem;  // 18px
$font-size-xl: 1.25rem;   // 20px
$font-size-2xl: 1.5rem;   // 24px
$font-size-3xl: 1.875rem; // 30px
$font-size-4xl: 2.25rem;  // 36px

// Line Heights
$line-height-tight: 1.25;
$line-height-normal: 1.5;
$line-height-relaxed: 1.75;
```

### Typography Classes
```scss
.text-h1 {
  font-size: $font-size-4xl;
  font-weight: $font-weight-bold;
  line-height: $line-height-tight;
}

.text-h2 {
  font-size: $font-size-3xl;
  font-weight: $font-weight-semibold;
  line-height: $line-height-tight;
}

.text-h3 {
  font-size: $font-size-2xl;
  font-weight: $font-weight-semibold;
  line-height: $line-height-tight;
}

.text-body {
  font-size: $font-size-base;
  font-weight: $font-weight-normal;
  line-height: $line-height-normal;
}

.text-caption {
  font-size: $font-size-sm;
  font-weight: $font-weight-normal;
  line-height: $line-height-normal;
  color: var(--text-secondary);
}

.text-code {
  font-family: $font-family-mono;
  font-size: $font-size-sm;
  font-weight: $font-weight-normal;
  line-height: $line-height-normal;
  background-color: var(--surface-hover);
  padding: 0.125rem 0.25rem;
  border-radius: 0.25rem;
}
```

## Icon System

### Icon Library
Using Lucide React icons for consistency and performance:

```typescript
import {
  Plus, Folder, FolderOpen, Download, Upload, Code, CheckCircle,
  AlertCircle, Clock, Settings, Activity, Template, Play, Pause,
  Stop, Trash2, Edit, Copy, ExternalLink, Search, Filter, Menu,
  X, ChevronDown, ChevronRight, ChevronLeft, ChevronUp,
  Info, Warning, Error, Success, Loading
} from 'lucide-react';

const iconMap = {
  'plus': Plus,
  'folder': Folder,
  'folder-open': FolderOpen,
  'download': Download,
  'upload': Upload,
  'code': Code,
  'check-circle': CheckCircle,
  'alert-circle': AlertCircle,
  'clock': Clock,
  'settings': Settings,
  'activity': Activity,
  'template': Template,
  'play': Play,
  'pause': Pause,
  'stop': Stop,
  'trash': Trash2,
  'edit': Edit,
  'copy': Copy,
  'external-link': ExternalLink,
  'search': Search,
  'filter': Filter,
  'menu': Menu,
  'x': X,
  'chevron-down': ChevronDown,
  'chevron-right': ChevronRight,
  'chevron-left': ChevronLeft,
  'chevron-up': ChevronUp,
  'info': Info,
  'warning': Warning,
  'error': Error,
  'success': Success,
  'loading': Loading
};
```

## Accessibility

### ARIA Implementation
```typescript
interface AccessibleButtonProps {
  'aria-label': string;
  'aria-describedby'?: string;
  'aria-pressed'?: boolean;
  'aria-expanded'?: boolean;
  'aria-disabled'?: boolean;
  role?: string;
  tabIndex?: number;
}

const AccessibleButton: React.FC<AccessibleButtonProps> = ({
  'aria-label': ariaLabel,
  'aria-describedby': ariaDescribedBy,
  'aria-pressed': ariaPressed,
  'aria-expanded': ariaExpanded,
  'aria-disabled': ariaDisabled,
  role = 'button',
  tabIndex = 0,
  ...props
}) => {
  return (
    <button
      aria-label={ariaLabel}
      aria-describedby={ariaDescribedBy}
      aria-pressed={ariaPressed}
      aria-expanded={ariaExpanded}
      aria-disabled={ariaDisabled}
      role={role}
      tabIndex={tabIndex}
      {...props}
    />
  );
};
```

### Keyboard Navigation
```typescript
const keyboardShortcuts = {
  'Ctrl+N': 'newProject',
  'Ctrl+O': 'openProject',
  'Ctrl+S': 'saveProject',
  'Ctrl+I': 'importSpec',
  'Ctrl+E': 'exportServer',
  'F5': 'validateSpec',
  'F6': 'generateServer',
  'F7': 'testServer',
  'Ctrl+F': 'find',
  'Ctrl+,': 'openSettings',
  'Ctrl+B': 'toggleSidebar',
  'Ctrl+L': 'showActivityLog',
  'Escape': 'closeModal',
  'Tab': 'navigateNext',
  'Shift+Tab': 'navigatePrevious',
  'Enter': 'activate',
  'Space': 'select',
  'ArrowUp': 'navigateUp',
  'ArrowDown': 'navigateDown',
  'ArrowLeft': 'navigateLeft',
  'ArrowRight': 'navigateRight'
};
```

## Animation and Transitions

### Transition System
```scss
// Transition Variables
$transition-fast: 150ms ease-in-out;
$transition-normal: 300ms ease-in-out;
$transition-slow: 500ms ease-in-out;

// Common Transitions
.transition-opacity {
  transition: opacity $transition-fast;
}

.transition-transform {
  transition: transform $transition-normal;
}

.transition-colors {
  transition: color $transition-fast, background-color $transition-fast, border-color $transition-fast;
}

.transition-all {
  transition: all $transition-normal;
}

// Animation Classes
.fade-in {
  animation: fadeIn $transition-normal ease-in-out;
}

.slide-in-right {
  animation: slideInRight $transition-normal ease-out;
}

.slide-in-left {
  animation: slideInLeft $transition-normal ease-out;
}

.scale-in {
  animation: scaleIn $transition-fast ease-out;
}

// Keyframes
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideInRight {
  from { transform: translateX(100%); }
  to { transform: translateX(0); }
}

@keyframes slideInLeft {
  from { transform: translateX(-100%); }
  to { transform: translateX(0); }
}

@keyframes scaleIn {
  from { transform: scale(0.9); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}
```

## Performance Considerations

### Optimization Strategies
- **Virtual Scrolling**: For large project lists and file trees
- **Lazy Loading**: Load components and data on demand
- **Memoization**: React.memo and useMemo for expensive operations
- **Code Splitting**: Dynamic imports for routes and features
- **Image Optimization**: Proper sizing and format selection
- **Bundle Analysis**: Regular bundle size monitoring

### Memory Management
```typescript
// Custom hook for cleanup
const useCleanup = (cleanup: () => void) => {
  useEffect(() => {
    return cleanup;
  }, [cleanup]);
};

// Event listener cleanup
const useEventListener = (event: string, handler: EventListener) => {
  useEffect(() => {
    document.addEventListener(event, handler);
    return () => document.removeEventListener(event, handler);
  }, [event, handler]);
};
```

This UI specification provides a comprehensive foundation for building a modern, accessible, and user-friendly desktop application interface for MCPWeaver while maintaining the focus on simplicity and efficiency.