// Type definitions for MCPWeaver frontend

export interface Project {
  id: string;
  name: string;
  specPath: string;
  specUrl: string;
  outputPath: string;
  settings: ProjectSettings;
  status: ProjectStatus;
  createdAt: string;
  updatedAt: string;
  lastGenerated?: string;
  generationCount: number;
}

export interface ProjectSettings {
  packageName: string;
  serverPort: number;
  enableLogging: boolean;
  logLevel: string;
  customTemplates: string[];
}

export type ProjectStatus = 'created' | 'validating' | 'ready' | 'generating' | 'error';

export interface GenerationJob {
  id: string;
  projectId: string;
  status: GenerationStatus;
  progress: number;
  currentStep: string;
  startTime: string;
  endTime?: string;
  results?: GenerationResults;
  errors?: GenerationError[];
  warnings?: string[];
}

export type GenerationStatus = 
  | 'started' 
  | 'parsing' 
  | 'mapping' 
  | 'generating' 
  | 'validating' 
  | 'completed' 
  | 'failed' 
  | 'cancelled';

export interface GenerationResults {
  serverPath: string;
  generatedFiles: GeneratedFile[];
  mcpTools: MCPTool[];
  statistics: GenerationStats;
}

export interface GeneratedFile {
  path: string;
  type: string;
  size: number;
  linesOfCode: number;
}

export interface GenerationStats {
  totalEndpoints: number;
  generatedTools: number;
  processingTime: number;
  specComplexity: string;
  templateVersion: string;
}

export interface GenerationError {
  type: string;
  message: string;
  details: string;
  suggestions?: string[];
  location?: ErrorLocation;
}

export interface ErrorLocation {
  file: string;
  line: number;
  column: number;
}

export interface MCPTool {
  name: string;
  description: string;
  inputSchema: InputSchema;
  method: string;
  path: string;
  baseURL: string;
}

export interface InputSchema {
  type: string;
  properties: Record<string, Property>;
  required: string[];
}

export interface Property {
  type: string;
  description?: string;
  example?: any;
  enum?: string[];
  format?: string;
  items?: Property;
}

export interface GenerationProgress {
  jobId: string;
  progress: number;
  currentStep: string;
  message: string;
  timestamp: string;
  processingRate?: number;
  estimatedTimeRemaining?: number;
  filesGenerated?: number;
  errorCount?: number;
  warningCount?: number;
  memoryUsage?: number;
}

export interface ProgressMetrics {
  startTime: string;
  currentTime: string;
  elapsedTime: number;
  processingRate: number;
  estimatedTimeRemaining: number;
  filesGenerated: number;
  totalFiles: number;
  errorCount: number;
  warningCount: number;
  memoryUsage: number;
  cpuUsage: number;
}

export interface ProgressHistoryEntry {
  jobId: string;
  projectId: string;
  projectName: string;
  status: GenerationStatus;
  progress: number;
  startTime: string;
  endTime?: string;
  duration?: number;
  success: boolean;
  errorMessage?: string;
  statistics?: GenerationStats;
}

export interface Notification {
  id: string;
  type: 'info' | 'success' | 'warning' | 'error';
  title: string;
  message: string;
  timestamp: string;
  actions?: string[];
  read: boolean;
}

export interface SystemInfo {
  version: string;
  os: string;
  arch: string;
  platform: string;
  totalMemory: number;
  availableMemory: number;
}

export interface AppSettings {
  theme: string;
  language: string;
  autoSave: boolean;
  defaultOutputPath: string;
  recentProjects: string[];
  recentFiles: string[];
  windowSettings: WindowSettings;
  editorSettings: EditorSettings;
  generationSettings: GenerationSettings;
  notificationSettings: NotificationSettings;
  appearanceSettings: AppearanceSettings;
  updateSettings: UpdateSettings;
}

export interface WindowSettings {
  width: number;
  height: number;
  maximized: boolean;
  x: number;
  y: number;
}

export interface EditorSettings {
  fontSize: number;
  fontFamily: string;
  tabSize: number;
  wordWrap: boolean;
  lineNumbers: boolean;
  syntaxHighlight: boolean;
}

export interface GenerationSettings {
  defaultTemplate: string;
  enableValidation: boolean;
  autoOpenOutput: boolean;
  showAdvancedOptions: boolean;
  backupOnGenerate: boolean;
  customTemplates: string[];
  performanceMode: boolean;
  maxWorkers: number;
}

export interface NotificationSettings {
  enableDesktopNotifications: boolean;
  enableSoundNotifications: boolean;
  notificationPosition: string;
  notificationDuration: number;
  soundVolume: number;
  showGenerationProgress: boolean;
  showErrorNotifications: boolean;
  showSuccessNotifications: boolean;
}

export interface AppearanceSettings {
  uiTheme: string;
  accentColor: string;
  windowOpacity: number;
  showAnimation: boolean;
  reducedMotion: boolean;
  fontScale: number;
  compactMode: boolean;
  showSidebar: boolean;
  sidebarPosition: string;
  showStatusBar: boolean;
  showToolbar: boolean;
}

export interface UpdateSettings {
  enableAutoUpdate: boolean;
  updateChannel: string;
  checkInterval: number;
  downloadInBackground: boolean;
  notifyOnAvailable: boolean;
  autoInstallSecurity: boolean;
  lastChecked?: string;
  lastUpdateInstalled?: string;
}

export interface ValidationResult {
  valid: boolean;
  errors: ValidationError[];
  warnings: ValidationWarning[];
  suggestions: string[];
  specInfo?: SpecInfo;
  validationTime: number;
  cacheHit: boolean;
  validatedAt: string;
}

export interface ValidationError {
  type: string;
  message: string;
  path: string;
  line?: number;
  column?: number;
  severity: string;
  code: string;
  location?: ErrorLocation;
}

export interface ValidationWarning {
  type: string;
  message: string;
  path: string;
  suggestion: string;
}

export interface SpecInfo {
  version: string;
  title: string;
  description: string;
  operationCount: number;
  schemaCount: number;
  securitySchemes: SecurityScheme[];
  servers: ServerInfo[];
}

export interface SecurityScheme {
  type: string;
  name: string;
  description: string;
}

export interface ServerInfo {
  url: string;
  description: string;
}

// Event types for Wails events
export interface WailsEvent<T = any> {
  data: T;
  timestamp: string;
}

export interface GenerationStartedEvent extends WailsEvent<GenerationJob> {}
export interface GenerationProgressEvent extends WailsEvent<GenerationProgress> {}
export interface GenerationCompletedEvent extends WailsEvent<GenerationJob> {}
export interface GenerationFailedEvent extends WailsEvent<{ jobId: string; type: string; message: string }> {}
export interface GenerationCancelledEvent extends WailsEvent<GenerationJob> {}
export interface ProjectUpdatedEvent extends WailsEvent<Project> {}
export interface SystemNotificationEvent extends WailsEvent<Notification> {}

// Error Handling Types
export interface APIError {
  type: string;
  code: string;
  message: string;
  details?: Record<string, string>;
  timestamp: string;
  suggestions?: string[];
  correlationId?: string;
  severity: 'low' | 'medium' | 'high' | 'critical';
  recoverable: boolean;
  retryAfter?: number;
  context?: ErrorContext;
}

export interface ErrorContext {
  operation?: string;
  component?: string;
  projectId?: string;
  userId?: string;
  sessionId?: string;
  requestId?: string;
  stackTrace?: string;
  metadata?: Record<string, string>;
}

export interface ErrorCollection {
  errors: APIError[];
  warnings: APIError[];
  operation: string;
  totalItems: number;
  failedItems: number;
  timestamp: string;
}

export interface RetryPolicy {
  maxRetries: number;
  initialDelay: number;
  maxDelay: number;
  backoffMultiplier: number;
  jitterEnabled: boolean;
  retryableErrors: string[];
}

export interface RetryResult {
  success: boolean;
  attempts: number;
  lastError?: APIError;
  totalDelay: number;
  startTime: string;
  endTime: string;
}

// File Import/Export Types
export interface ImportResult {
  content: string;
  valid: boolean;
  specInfo?: SpecInfo;
  errors?: string[];
  warnings?: string[];
  importedFrom: 'file' | 'url';
  filePath?: string;
  sourceUrl?: string;
  fileSize: number;
  importedAt: string;
}

export interface ExportResult {
  projectId: string;
  projectName: string;
  targetDir: string;
  exportedFiles: ExportedFile[];
  totalFiles: number;
  totalSize: number;
  exportedAt: string;
}

export interface ExportedFile {
  name: string;
  path: string;
  size: number;
  modifiedTime: string;
}

export interface FileOperationProgress {
  operationId: string;
  type: 'import' | 'export';
  progress: number;
  currentFile: string;
  totalFiles: number;
  processedFiles: number;
  startTime: string;
  elapsedTime: number;
  estimatedRemaining: number;
}

export interface FileFilter {
  displayName: string;
  pattern: string;
  extensions: string[];
}

export interface RecentFile {
  path: string;
  name: string;
  size: number;
  lastAccessed: string;
  type: 'spec' | 'export';
}

// Template System Types
export interface Template {
  id: string;
  name: string;
  description: string;
  version: string;
  author: string;
  type: TemplateType;
  path: string;
  isBuiltIn: boolean;
  variables: TemplateVariable[];
  createdAt: string;
  updatedAt: string;
}

export type TemplateType = 'default' | 'custom' | 'plugin';

export interface TemplateVariable {
  name: string;
  description: string;
  type: string;
  defaultValue: string;
  required: boolean;
  options?: string[];
  validation?: string;
}

export interface CreateTemplateRequest {
  name: string;
  description: string;
  version: string;
  author: string;
  type: TemplateType;
  path: string;
  variables?: TemplateVariable[];
}

export interface UpdateTemplateRequest {
  name?: string;
  description?: string;
  version?: string;
  author?: string;
  type?: TemplateType;
  path?: string;
  variables?: TemplateVariable[];
}

export interface TemplateValidationResult {
  valid: boolean;
  errors?: TemplateError[];
  warnings?: TemplateWarning[];
  suggestions?: string[];
  performance?: TemplatePerformance;
  dependencies?: TemplateDependency[];
}

export interface TemplateError {
  type: string;
  message: string;
  line?: number;
  column?: number;
  severity: string;
}

export interface TemplateWarning {
  type: string;
  message: string;
  line?: number;
  suggestion?: string;
}

export interface TemplatePerformance {
  renderTime: string; // Duration as string
  memoryUsage: number;
  complexity: string;
  cacheHit: boolean;
}

export interface TemplateDependency {
  name: string;
  version: string;
  required: boolean;
  type: string;
}

export interface TemplateTestRequest {
  templateId: string;
  testData: Record<string, any>;
  options?: TemplateTestOptions;
}

export interface TemplateTestOptions {
  validateOutput?: boolean;
  measurePerformance?: boolean;
  generateReport?: boolean;
}

export interface TemplateTestResult {
  success: boolean;
  output?: string;
  errors?: TemplateError[];
  warnings?: TemplateWarning[];
  performance?: TemplatePerformance;
  report?: TemplateTestReport;
}

export interface TemplateTestReport {
  templateId: string;
  testExecutedAt: string;
  executionTime: string;
  outputSize: number;
  variablesUsed: string[];
  functionsUsed: string[];
  recommendations: string[];
}

export interface TemplateImportRequest {
  source: string;
  path?: string;
  url?: string;
  marketplaceId?: string;
  options?: TemplateImportOptions;
}

export interface TemplateImportOptions {
  overwriteExisting?: boolean;
  validateOnly?: boolean;
  includeDependencies?: boolean;
  targetType?: TemplateType;
}

export interface TemplateExportRequest {
  templateId: string;
  format: string;
  targetPath: string;
  options?: TemplateExportOptions;
}

export interface TemplateExportOptions {
  includeDocumentation?: boolean;
  includeExamples?: boolean;
  includeDependencies?: boolean;
  minify?: boolean;
}

export interface TemplateMarketplaceItem {
  id: string;
  name: string;
  description: string;
  version: string;
  author: string;
  type: TemplateType;
  tags: string[];
  rating: number;
  downloads: number;
  createdAt: string;
  updatedAt: string;
  license: string;
  repository?: string;
  homePage?: string;
  screenshots?: string[];
  variables: TemplateVariable[];
  dependencies: TemplateDependency[];
}

export interface TemplateSearchRequest {
  query?: string;
  type?: TemplateType;
  tags?: string[];
  author?: string;
  minRating?: number;
  sortBy?: string;
  sortOrder?: string;
  limit?: number;
  offset?: number;
}

export interface TemplateSearchResult {
  items: TemplateMarketplaceItem[];
  total: number;
  limit: number;
  offset: number;
  hasMore: boolean;
  searchTime: string;
}

// Activity Log and Monitoring Types
export interface ActivityLogEntry {
  id: string;
  timestamp: string;
  level: LogLevel;
  component: string;
  operation: string;
  message: string;
  details?: string;
  duration?: string; // Duration as string
  projectId?: string;
  userAction: boolean;
  metadata?: Record<string, any>;
}

export type LogLevel = 'debug' | 'info' | 'warn' | 'error' | 'fatal';

export interface LogFilter {
  level?: LogLevel;
  component?: string;
  operation?: string;
  projectId?: string;
  userAction?: boolean;
  startTime?: string;
  endTime?: string;
  search?: string;
  limit?: number;
}

export interface LogSearchRequest {
  query: string;
  filter: LogFilter;
  limit: number;
  offset: number;
}

export interface LogSearchResult {
  entries: ActivityLogEntry[];
  total: number;
  hasMore: boolean;
  searchTime: string; // Duration as string
}

export interface LogExportRequest {
  filter: LogFilter;
  format: 'json' | 'csv' | 'txt';
  filePath: string;
}

export interface LogExportResult {
  filePath: string;
  entriesCount: number;
  fileSize: number;
  exportTime: string; // Duration as string
  format: string;
}

export interface ApplicationStatus {
  status: StatusLevel;
  message: string;
  activeOperations: number;
  lastUpdate: string;
  systemHealth: SystemHealth;
}

export type StatusLevel = 'idle' | 'working' | 'error' | 'warning';

export interface SystemHealth {
  memoryUsage: number; // MB
  cpuUsage: number; // Percentage
  diskSpace: number; // GB available
  databaseSize: number; // MB
  temporaryFiles: number; // Count
  activeConnections: number;
}

export interface ErrorReport {
  id: string;
  timestamp: string;
  type: ErrorType;
  severity: ErrorSeverity;
  component: string;
  operation: string;
  message: string;
  details?: string;
  stackTrace?: string;
  userContext: UserContext;
  systemInfo: SystemInfoReport;
  recovery: RecoveryInfo;
  frequency: number;
  firstSeen: string;
  lastSeen: string;
}

export type ErrorType = 'validation' | 'system' | 'network' | 'filesystem' | 'database' | 'generation' | 'permission' | 'configuration' | 'authentication';

export type ErrorSeverity = 'low' | 'medium' | 'high' | 'critical';

export interface UserContext {
  projectId?: string;
  projectName?: string;
  userAction?: string;
  uiState?: string;
  recentActions?: string[];
  settings?: Record<string, string>;
}

export interface SystemInfoReport {
  os: string;
  architecture: string;
  goVersion: string;
  appVersion: string;
  memoryMB: number;
  cpuUsage: number;
  diskSpaceGB: number;
  databaseSizeMB: number;
}

export interface RecoveryInfo {
  attempted: boolean;
  successful: boolean;
  method?: string;
  duration?: string; // Duration as string
  userInteraction: boolean;
  dataLoss: boolean;
}

export interface LogConfig {
  level: LogLevel;
  bufferSize: number;
  retentionDays: number;
  enableConsole: boolean;
  enableBuffer: boolean;
  flushInterval: string; // Duration as string
}