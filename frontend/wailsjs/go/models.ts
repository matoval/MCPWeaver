export namespace app {
	
	export class ErrorContext {
	    operation: string;
	    component: string;
	    projectId?: string;
	    userId?: string;
	    sessionId?: string;
	    requestId?: string;
	    stackTrace?: string;
	    metadata?: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new ErrorContext(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.operation = source["operation"];
	        this.component = source["component"];
	        this.projectId = source["projectId"];
	        this.userId = source["userId"];
	        this.sessionId = source["sessionId"];
	        this.requestId = source["requestId"];
	        this.stackTrace = source["stackTrace"];
	        this.metadata = source["metadata"];
	    }
	}
	export class APIError {
	    type: string;
	    code: string;
	    message: string;
	    details?: Record<string, string>;
	    // Go type: time
	    timestamp: any;
	    suggestions?: string[];
	    correlationId?: string;
	    severity: string;
	    recoverable: boolean;
	    retryAfter?: number;
	    context?: ErrorContext;
	
	    static createFrom(source: any = {}) {
	        return new APIError(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.code = source["code"];
	        this.message = source["message"];
	        this.details = source["details"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.suggestions = source["suggestions"];
	        this.correlationId = source["correlationId"];
	        this.severity = source["severity"];
	        this.recoverable = source["recoverable"];
	        this.retryAfter = source["retryAfter"];
	        this.context = this.convertValues(source["context"], ErrorContext);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ActivityLogEntry {
	    id: string;
	    // Go type: time
	    timestamp: any;
	    level: number;
	    component: string;
	    operation: string;
	    message: string;
	    details?: string;
	    duration?: number;
	    projectId?: string;
	    userAction: boolean;
	    metadata?: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new ActivityLogEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.level = source["level"];
	        this.component = source["component"];
	        this.operation = source["operation"];
	        this.message = source["message"];
	        this.details = source["details"];
	        this.duration = source["duration"];
	        this.projectId = source["projectId"];
	        this.userAction = source["userAction"];
	        this.metadata = source["metadata"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AggregatedMetrics {
	    operation: string;
	    totalExecutions: number;
	    successCount: number;
	    failureCount: number;
	    successRate: number;
	    averageDuration: number;
	    minDuration: number;
	    maxDuration: number;
	    p50Duration: number;
	    p95Duration: number;
	    p99Duration: number;
	    averageMemory: number;
	    averageCpu: number;
	    cacheHitRate: number;
	    // Go type: time
	    lastUpdated: any;
	
	    static createFrom(source: any = {}) {
	        return new AggregatedMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.operation = source["operation"];
	        this.totalExecutions = source["totalExecutions"];
	        this.successCount = source["successCount"];
	        this.failureCount = source["failureCount"];
	        this.successRate = source["successRate"];
	        this.averageDuration = source["averageDuration"];
	        this.minDuration = source["minDuration"];
	        this.maxDuration = source["maxDuration"];
	        this.p50Duration = source["p50Duration"];
	        this.p95Duration = source["p95Duration"];
	        this.p99Duration = source["p99Duration"];
	        this.averageMemory = source["averageMemory"];
	        this.averageCpu = source["averageCpu"];
	        this.cacheHitRate = source["cacheHitRate"];
	        this.lastUpdated = this.convertValues(source["lastUpdated"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UpdateRetryPolicy {
	    maxRetries: number;
	    initialDelay: number;
	    maxDelay: number;
	    backoffMultiplier: number;
	    retryOnNetworkError: boolean;
	    retryOnVerificationError: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UpdateRetryPolicy(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.maxRetries = source["maxRetries"];
	        this.initialDelay = source["initialDelay"];
	        this.maxDelay = source["maxDelay"];
	        this.backoffMultiplier = source["backoffMultiplier"];
	        this.retryOnNetworkError = source["retryOnNetworkError"];
	        this.retryOnVerificationError = source["retryOnVerificationError"];
	    }
	}
	export class UpdateSchedule {
	    type: string;
	    time?: string;
	    dayOfWeek?: number;
	    dayOfMonth?: number;
	    // Go type: time
	    nextCheck: any;
	
	    static createFrom(source: any = {}) {
	        return new UpdateSchedule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.time = source["time"];
	        this.dayOfWeek = source["dayOfWeek"];
	        this.dayOfMonth = source["dayOfMonth"];
	        this.nextCheck = this.convertValues(source["nextCheck"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UpdateSettings {
	    enabled: boolean;
	    autoCheck: boolean;
	    checkInterval: number;
	    autoDownload: boolean;
	    autoInstall: boolean;
	    promptUser: boolean;
	    schedule?: UpdateSchedule;
	    updateChannel: string;
	    preReleaseEnabled: boolean;
	    bandwidthLimit: number;
	    retryPolicy: UpdateRetryPolicy;
	
	    static createFrom(source: any = {}) {
	        return new UpdateSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.autoCheck = source["autoCheck"];
	        this.checkInterval = source["checkInterval"];
	        this.autoDownload = source["autoDownload"];
	        this.autoInstall = source["autoInstall"];
	        this.promptUser = source["promptUser"];
	        this.schedule = this.convertValues(source["schedule"], UpdateSchedule);
	        this.updateChannel = source["updateChannel"];
	        this.preReleaseEnabled = source["preReleaseEnabled"];
	        this.bandwidthLimit = source["bandwidthLimit"];
	        this.retryPolicy = this.convertValues(source["retryPolicy"], UpdateRetryPolicy);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AppearanceSettings {
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
	
	    static createFrom(source: any = {}) {
	        return new AppearanceSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uiTheme = source["uiTheme"];
	        this.accentColor = source["accentColor"];
	        this.windowOpacity = source["windowOpacity"];
	        this.showAnimation = source["showAnimation"];
	        this.reducedMotion = source["reducedMotion"];
	        this.fontScale = source["fontScale"];
	        this.compactMode = source["compactMode"];
	        this.showSidebar = source["showSidebar"];
	        this.sidebarPosition = source["sidebarPosition"];
	        this.showStatusBar = source["showStatusBar"];
	        this.showToolbar = source["showToolbar"];
	    }
	}
	export class NotificationSettings {
	    enableDesktopNotifications: boolean;
	    enableSoundNotifications: boolean;
	    notificationPosition: string;
	    notificationDuration: number;
	    soundVolume: number;
	    showGenerationProgress: boolean;
	    showErrorNotifications: boolean;
	    showSuccessNotifications: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NotificationSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enableDesktopNotifications = source["enableDesktopNotifications"];
	        this.enableSoundNotifications = source["enableSoundNotifications"];
	        this.notificationPosition = source["notificationPosition"];
	        this.notificationDuration = source["notificationDuration"];
	        this.soundVolume = source["soundVolume"];
	        this.showGenerationProgress = source["showGenerationProgress"];
	        this.showErrorNotifications = source["showErrorNotifications"];
	        this.showSuccessNotifications = source["showSuccessNotifications"];
	    }
	}
	export class GenerationSettings {
	    defaultTemplate: string;
	    enableValidation: boolean;
	    autoOpenOutput: boolean;
	    showAdvancedOptions: boolean;
	    backupOnGenerate: boolean;
	    customTemplates: string[];
	    performanceMode: boolean;
	    maxWorkers: number;
	
	    static createFrom(source: any = {}) {
	        return new GenerationSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.defaultTemplate = source["defaultTemplate"];
	        this.enableValidation = source["enableValidation"];
	        this.autoOpenOutput = source["autoOpenOutput"];
	        this.showAdvancedOptions = source["showAdvancedOptions"];
	        this.backupOnGenerate = source["backupOnGenerate"];
	        this.customTemplates = source["customTemplates"];
	        this.performanceMode = source["performanceMode"];
	        this.maxWorkers = source["maxWorkers"];
	    }
	}
	export class EditorSettings {
	    fontSize: number;
	    fontFamily: string;
	    tabSize: number;
	    wordWrap: boolean;
	    lineNumbers: boolean;
	    syntaxHighlight: boolean;
	
	    static createFrom(source: any = {}) {
	        return new EditorSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fontSize = source["fontSize"];
	        this.fontFamily = source["fontFamily"];
	        this.tabSize = source["tabSize"];
	        this.wordWrap = source["wordWrap"];
	        this.lineNumbers = source["lineNumbers"];
	        this.syntaxHighlight = source["syntaxHighlight"];
	    }
	}
	export class WindowSettings {
	    width: number;
	    height: number;
	    maximized: boolean;
	    x: number;
	    y: number;
	
	    static createFrom(source: any = {}) {
	        return new WindowSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.width = source["width"];
	        this.height = source["height"];
	        this.maximized = source["maximized"];
	        this.x = source["x"];
	        this.y = source["y"];
	    }
	}
	export class AppSettings {
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
	
	    static createFrom(source: any = {}) {
	        return new AppSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.theme = source["theme"];
	        this.language = source["language"];
	        this.autoSave = source["autoSave"];
	        this.defaultOutputPath = source["defaultOutputPath"];
	        this.recentProjects = source["recentProjects"];
	        this.recentFiles = source["recentFiles"];
	        this.windowSettings = this.convertValues(source["windowSettings"], WindowSettings);
	        this.editorSettings = this.convertValues(source["editorSettings"], EditorSettings);
	        this.generationSettings = this.convertValues(source["generationSettings"], GenerationSettings);
	        this.notificationSettings = this.convertValues(source["notificationSettings"], NotificationSettings);
	        this.appearanceSettings = this.convertValues(source["appearanceSettings"], AppearanceSettings);
	        this.updateSettings = this.convertValues(source["updateSettings"], UpdateSettings);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class SystemHealth {
	    memoryUsage: number;
	    cpuUsage: number;
	    diskSpace: number;
	    databaseSize: number;
	    temporaryFiles: number;
	    activeConnections: number;
	
	    static createFrom(source: any = {}) {
	        return new SystemHealth(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.memoryUsage = source["memoryUsage"];
	        this.cpuUsage = source["cpuUsage"];
	        this.diskSpace = source["diskSpace"];
	        this.databaseSize = source["databaseSize"];
	        this.temporaryFiles = source["temporaryFiles"];
	        this.activeConnections = source["activeConnections"];
	    }
	}
	export class ApplicationStatus {
	    status: string;
	    message: string;
	    activeOperations: number;
	    // Go type: time
	    lastUpdate: any;
	    systemHealth: SystemHealth;
	
	    static createFrom(source: any = {}) {
	        return new ApplicationStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.message = source["message"];
	        this.activeOperations = source["activeOperations"];
	        this.lastUpdate = this.convertValues(source["lastUpdate"], null);
	        this.systemHealth = this.convertValues(source["systemHealth"], SystemHealth);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BackupInfo {
	    path: string;
	    name: string;
	    version: string;
	    size: number;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new BackupInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.name = source["name"];
	        this.version = source["version"];
	        this.size = source["size"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BackupValidation {
	    path: string;
	    valid: boolean;
	    size: number;
	    errors: string[];
	    warnings: string[];
	    // Go type: time
	    validatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new BackupValidation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.valid = source["valid"];
	        this.size = source["size"];
	        this.errors = source["errors"];
	        this.warnings = source["warnings"];
	        this.validatedAt = this.convertValues(source["validatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CategoryPreference {
	    enabled: boolean;
	    toastEnabled: boolean;
	    systemEnabled: boolean;
	    soundEnabled: boolean;
	    minPriority: string;
	    maxPerHour: number;
	
	    static createFrom(source: any = {}) {
	        return new CategoryPreference(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.toastEnabled = source["toastEnabled"];
	        this.systemEnabled = source["systemEnabled"];
	        this.soundEnabled = source["soundEnabled"];
	        this.minPriority = source["minPriority"];
	        this.maxPerHour = source["maxPerHour"];
	    }
	}
	export class CategoryStats {
	    sent: number;
	    read: number;
	    dismissed: number;
	    interacted: number;
	    readRate: number;
	
	    static createFrom(source: any = {}) {
	        return new CategoryStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sent = source["sent"];
	        this.read = source["read"];
	        this.dismissed = source["dismissed"];
	        this.interacted = source["interacted"];
	        this.readRate = source["readRate"];
	    }
	}
	export class ClientInfo {
	    version: string;
	    platform: string;
	    architecture: string;
	    os: string;
	    osVersion: string;
	
	    static createFrom(source: any = {}) {
	        return new ClientInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.platform = source["platform"];
	        this.architecture = source["architecture"];
	        this.os = source["os"];
	        this.osVersion = source["osVersion"];
	    }
	}
	export class ProjectSettings {
	    packageName: string;
	    serverPort: number;
	    enableLogging: boolean;
	    logLevel: string;
	    customTemplates?: string[];
	
	    static createFrom(source: any = {}) {
	        return new ProjectSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.packageName = source["packageName"];
	        this.serverPort = source["serverPort"];
	        this.enableLogging = source["enableLogging"];
	        this.logLevel = source["logLevel"];
	        this.customTemplates = source["customTemplates"];
	    }
	}
	export class CreateProjectRequest {
	    name: string;
	    specPath?: string;
	    specUrl?: string;
	    outputPath: string;
	    settings: ProjectSettings;
	
	    static createFrom(source: any = {}) {
	        return new CreateProjectRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.specPath = source["specPath"];
	        this.specUrl = source["specUrl"];
	        this.outputPath = source["outputPath"];
	        this.settings = this.convertValues(source["settings"], ProjectSettings);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DoNotDisturbSchedule {
	    enabled: boolean;
	    startTime: string;
	    endTime: string;
	    days: number[];
	    exceptions: time.Time[];
	    allowUrgent: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DoNotDisturbSchedule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.startTime = source["startTime"];
	        this.endTime = source["endTime"];
	        this.days = source["days"];
	        this.exceptions = this.convertValues(source["exceptions"], time.Time);
	        this.allowUrgent = source["allowUrgent"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class ErrorLocation {
	    file: string;
	    line: number;
	    column: number;
	
	    static createFrom(source: any = {}) {
	        return new ErrorLocation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.file = source["file"];
	        this.line = source["line"];
	        this.column = source["column"];
	    }
	}
	export class RecoveryInfo {
	    attempted: boolean;
	    successful: boolean;
	    method?: string;
	    duration?: number;
	    userInteraction: boolean;
	    dataLoss: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RecoveryInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.attempted = source["attempted"];
	        this.successful = source["successful"];
	        this.method = source["method"];
	        this.duration = source["duration"];
	        this.userInteraction = source["userInteraction"];
	        this.dataLoss = source["dataLoss"];
	    }
	}
	export class SystemInfo {
	    os: string;
	    architecture: string;
	    goVersion: string;
	    appVersion: string;
	    memoryMB: number;
	    cpuUsage: number;
	    diskSpaceGB: number;
	    databaseSizeMB: number;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.os = source["os"];
	        this.architecture = source["architecture"];
	        this.goVersion = source["goVersion"];
	        this.appVersion = source["appVersion"];
	        this.memoryMB = source["memoryMB"];
	        this.cpuUsage = source["cpuUsage"];
	        this.diskSpaceGB = source["diskSpaceGB"];
	        this.databaseSizeMB = source["databaseSizeMB"];
	    }
	}
	export class UserContext {
	    projectId?: string;
	    projectName?: string;
	    userAction?: string;
	    uiState?: string;
	    recentActions?: string[];
	    settings?: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new UserContext(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.projectId = source["projectId"];
	        this.projectName = source["projectName"];
	        this.userAction = source["userAction"];
	        this.uiState = source["uiState"];
	        this.recentActions = source["recentActions"];
	        this.settings = source["settings"];
	    }
	}
	export class ErrorReport {
	    id: string;
	    // Go type: time
	    timestamp: any;
	    type: string;
	    severity: string;
	    component: string;
	    operation: string;
	    message: string;
	    details?: string;
	    stackTrace?: string;
	    userContext: UserContext;
	    systemInfo: SystemInfo;
	    recovery: RecoveryInfo;
	    frequency: number;
	    // Go type: time
	    firstSeen: any;
	    // Go type: time
	    lastSeen: any;
	
	    static createFrom(source: any = {}) {
	        return new ErrorReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.type = source["type"];
	        this.severity = source["severity"];
	        this.component = source["component"];
	        this.operation = source["operation"];
	        this.message = source["message"];
	        this.details = source["details"];
	        this.stackTrace = source["stackTrace"];
	        this.userContext = this.convertValues(source["userContext"], UserContext);
	        this.systemInfo = this.convertValues(source["systemInfo"], SystemInfo);
	        this.recovery = this.convertValues(source["recovery"], RecoveryInfo);
	        this.frequency = source["frequency"];
	        this.firstSeen = this.convertValues(source["firstSeen"], null);
	        this.lastSeen = this.convertValues(source["lastSeen"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ExportedFile {
	    name: string;
	    path: string;
	    size: number;
	    // Go type: time
	    modifiedTime: any;
	
	    static createFrom(source: any = {}) {
	        return new ExportedFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.size = source["size"];
	        this.modifiedTime = this.convertValues(source["modifiedTime"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ExportResult {
	    projectId: string;
	    projectName: string;
	    targetDir: string;
	    exportedFiles: ExportedFile[];
	    totalFiles: number;
	    totalSize: number;
	    // Go type: time
	    exportedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new ExportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.projectId = source["projectId"];
	        this.projectName = source["projectName"];
	        this.targetDir = source["targetDir"];
	        this.exportedFiles = this.convertValues(source["exportedFiles"], ExportedFile);
	        this.totalFiles = source["totalFiles"];
	        this.totalSize = source["totalSize"];
	        this.exportedAt = this.convertValues(source["exportedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class FileFilter {
	    displayName: string;
	    pattern: string;
	    extensions: string[];
	
	    static createFrom(source: any = {}) {
	        return new FileFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.displayName = source["displayName"];
	        this.pattern = source["pattern"];
	        this.extensions = source["extensions"];
	    }
	}
	export class GeneratedFile {
	    path: string;
	    type: string;
	    size: number;
	    linesOfCode: number;
	
	    static createFrom(source: any = {}) {
	        return new GeneratedFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.type = source["type"];
	        this.size = source["size"];
	        this.linesOfCode = source["linesOfCode"];
	    }
	}
	export class GenerationError {
	    type: string;
	    message: string;
	    details?: string;
	    suggestions?: string[];
	    location?: ErrorLocation;
	
	    static createFrom(source: any = {}) {
	        return new GenerationError(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.message = source["message"];
	        this.details = source["details"];
	        this.suggestions = source["suggestions"];
	        this.location = this.convertValues(source["location"], ErrorLocation);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GenerationStats {
	    totalEndpoints: number;
	    generatedTools: number;
	    processingTime: number;
	    specComplexity: string;
	    templateVersion: string;
	
	    static createFrom(source: any = {}) {
	        return new GenerationStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalEndpoints = source["totalEndpoints"];
	        this.generatedTools = source["generatedTools"];
	        this.processingTime = source["processingTime"];
	        this.specComplexity = source["specComplexity"];
	        this.templateVersion = source["templateVersion"];
	    }
	}
	export class GenerationResults {
	    serverPath: string;
	    generatedFiles: GeneratedFile[];
	    mcpTools: mapping.MCPTool[];
	    statistics: GenerationStats;
	
	    static createFrom(source: any = {}) {
	        return new GenerationResults(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.serverPath = source["serverPath"];
	        this.generatedFiles = this.convertValues(source["generatedFiles"], GeneratedFile);
	        this.mcpTools = this.convertValues(source["mcpTools"], mapping.MCPTool);
	        this.statistics = this.convertValues(source["statistics"], GenerationStats);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GenerationJob {
	    id: string;
	    projectId: string;
	    status: string;
	    progress: number;
	    currentStep: string;
	    // Go type: time
	    startTime: any;
	    // Go type: time
	    endTime?: any;
	    results?: GenerationResults;
	    errors?: GenerationError[];
	    warnings?: string[];
	
	    static createFrom(source: any = {}) {
	        return new GenerationJob(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.projectId = source["projectId"];
	        this.status = source["status"];
	        this.progress = source["progress"];
	        this.currentStep = source["currentStep"];
	        this.startTime = this.convertValues(source["startTime"], null);
	        this.endTime = this.convertValues(source["endTime"], null);
	        this.results = this.convertValues(source["results"], GenerationResults);
	        this.errors = this.convertValues(source["errors"], GenerationError);
	        this.warnings = source["warnings"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	export class ServerInfo {
	    url: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new ServerInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.description = source["description"];
	    }
	}
	export class SecurityScheme {
	    type: string;
	    name: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new SecurityScheme(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.name = source["name"];
	        this.description = source["description"];
	    }
	}
	export class SpecInfo {
	    version: string;
	    title: string;
	    description: string;
	    operationCount: number;
	    schemaCount: number;
	    securitySchemes: SecurityScheme[];
	    servers: ServerInfo[];
	
	    static createFrom(source: any = {}) {
	        return new SpecInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.operationCount = source["operationCount"];
	        this.schemaCount = source["schemaCount"];
	        this.securitySchemes = this.convertValues(source["securitySchemes"], SecurityScheme);
	        this.servers = this.convertValues(source["servers"], ServerInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ImportResult {
	    content: string;
	    valid: boolean;
	    specInfo?: SpecInfo;
	    errors?: string[];
	    warnings?: string[];
	    importedFrom: string;
	    filePath?: string;
	    sourceUrl?: string;
	    fileSize: number;
	    // Go type: time
	    importedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new ImportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.content = source["content"];
	        this.valid = source["valid"];
	        this.specInfo = this.convertValues(source["specInfo"], SpecInfo);
	        this.errors = source["errors"];
	        this.warnings = source["warnings"];
	        this.importedFrom = source["importedFrom"];
	        this.filePath = source["filePath"];
	        this.sourceUrl = source["sourceUrl"];
	        this.fileSize = source["fileSize"];
	        this.importedAt = this.convertValues(source["importedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LogConfig {
	    level: number;
	    bufferSize: number;
	    retentionDays: number;
	    enableConsole: boolean;
	    enableBuffer: boolean;
	    flushInterval: number;
	
	    static createFrom(source: any = {}) {
	        return new LogConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.level = source["level"];
	        this.bufferSize = source["bufferSize"];
	        this.retentionDays = source["retentionDays"];
	        this.enableConsole = source["enableConsole"];
	        this.enableBuffer = source["enableBuffer"];
	        this.flushInterval = source["flushInterval"];
	    }
	}
	export class LogFilter {
	    level?: number;
	    component?: string;
	    operation?: string;
	    projectId?: string;
	    userAction?: boolean;
	    // Go type: time
	    startTime?: any;
	    // Go type: time
	    endTime?: any;
	    search?: string;
	    limit?: number;
	
	    static createFrom(source: any = {}) {
	        return new LogFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.level = source["level"];
	        this.component = source["component"];
	        this.operation = source["operation"];
	        this.projectId = source["projectId"];
	        this.userAction = source["userAction"];
	        this.startTime = this.convertValues(source["startTime"], null);
	        this.endTime = this.convertValues(source["endTime"], null);
	        this.search = source["search"];
	        this.limit = source["limit"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LogExportRequest {
	    filter: LogFilter;
	    format: string;
	    filePath: string;
	
	    static createFrom(source: any = {}) {
	        return new LogExportRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filter = this.convertValues(source["filter"], LogFilter);
	        this.format = source["format"];
	        this.filePath = source["filePath"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LogExportResult {
	    filePath: string;
	    entriesCount: number;
	    fileSize: number;
	    exportTime: number;
	    format: string;
	
	    static createFrom(source: any = {}) {
	        return new LogExportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filePath = source["filePath"];
	        this.entriesCount = source["entriesCount"];
	        this.fileSize = source["fileSize"];
	        this.exportTime = source["exportTime"];
	        this.format = source["format"];
	    }
	}
	
	export class LogSearchRequest {
	    query: string;
	    filter: LogFilter;
	    limit: number;
	    offset: number;
	
	    static createFrom(source: any = {}) {
	        return new LogSearchRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.query = source["query"];
	        this.filter = this.convertValues(source["filter"], LogFilter);
	        this.limit = source["limit"];
	        this.offset = source["offset"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LogSearchResult {
	    entries: ActivityLogEntry[];
	    total: number;
	    hasMore: boolean;
	    searchTime: number;
	
	    static createFrom(source: any = {}) {
	        return new LogSearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entries = this.convertValues(source["entries"], ActivityLogEntry);
	        this.total = source["total"];
	        this.hasMore = source["hasMore"];
	        this.searchTime = source["searchTime"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class NotificationActionBtn {
	    id: string;
	    label: string;
	    type: string;
	    icon?: string;
	    style: string;
	    callback?: string;
	    data?: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new NotificationActionBtn(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.label = source["label"];
	        this.type = source["type"];
	        this.icon = source["icon"];
	        this.style = source["style"];
	        this.callback = source["callback"];
	        this.data = source["data"];
	    }
	}
	export class NotificationFilter {
	    id: string;
	    name: string;
	    enabled: boolean;
	    condition: string;
	    action: string;
	    keywords?: string[];
	    category?: string;
	    priority?: string;
	    source?: string;
	
	    static createFrom(source: any = {}) {
	        return new NotificationFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.enabled = source["enabled"];
	        this.condition = source["condition"];
	        this.action = source["action"];
	        this.keywords = source["keywords"];
	        this.category = source["category"];
	        this.priority = source["priority"];
	        this.source = source["source"];
	    }
	}
	export class NotificationFilterTestResult {
	    filterId: string;
	    filterName: string;
	    passes: boolean;
	    action: string;
	    testResult: string;
	    // Go type: time
	    testedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new NotificationFilterTestResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filterId = source["filterId"];
	        this.filterName = source["filterName"];
	        this.passes = source["passes"];
	        this.action = source["action"];
	        this.testResult = source["testResult"];
	        this.testedAt = this.convertValues(source["testedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class NotificationHistory {
	    id: string;
	    type: string;
	    title: string;
	    message: string;
	    icon?: string;
	    actions?: NotificationActionBtn[];
	    category: string;
	    priority: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    readAt?: any;
	    // Go type: time
	    dismissedAt?: any;
	    // Go type: time
	    interactedAt?: any;
	    actionTaken?: string;
	    source: string;
	    metadata?: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new NotificationHistory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.title = source["title"];
	        this.message = source["message"];
	        this.icon = source["icon"];
	        this.actions = this.convertValues(source["actions"], NotificationActionBtn);
	        this.category = source["category"];
	        this.priority = source["priority"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.readAt = this.convertValues(source["readAt"], null);
	        this.dismissedAt = this.convertValues(source["dismissedAt"], null);
	        this.interactedAt = this.convertValues(source["interactedAt"], null);
	        this.actionTaken = source["actionTaken"];
	        this.source = source["source"];
	        this.metadata = source["metadata"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class NotificationPreferences {
	    categories: Record<string, CategoryPreference>;
	    filters: NotificationFilter[];
	    sounds: Record<string, string>;
	    volumes: Record<string, number>;
	
	    static createFrom(source: any = {}) {
	        return new NotificationPreferences(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.categories = this.convertValues(source["categories"], CategoryPreference, true);
	        this.filters = this.convertValues(source["filters"], NotificationFilter);
	        this.sounds = source["sounds"];
	        this.volumes = source["volumes"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class NotificationProgress {
	    current: number;
	    total: number;
	    percent: number;
	    label?: string;
	
	    static createFrom(source: any = {}) {
	        return new NotificationProgress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.current = source["current"];
	        this.total = source["total"];
	        this.percent = source["percent"];
	        this.label = source["label"];
	    }
	}
	export class NotificationQueueStatus {
	    size: number;
	    maxSize: number;
	    paused: boolean;
	    drainRate: number;
	    queuedItems: number;
	
	    static createFrom(source: any = {}) {
	        return new NotificationQueueStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.size = source["size"];
	        this.maxSize = source["maxSize"];
	        this.paused = source["paused"];
	        this.drainRate = source["drainRate"];
	        this.queuedItems = source["queuedItems"];
	    }
	}
	
	export class PriorityStats {
	    sent: number;
	    read: number;
	    dismissed: number;
	    interacted: number;
	    readRate: number;
	
	    static createFrom(source: any = {}) {
	        return new PriorityStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sent = source["sent"];
	        this.read = source["read"];
	        this.dismissed = source["dismissed"];
	        this.interacted = source["interacted"];
	        this.readRate = source["readRate"];
	    }
	}
	export class NotificationStats {
	    totalSent: number;
	    totalToast: number;
	    totalSystem: number;
	    totalRead: number;
	    totalDismissed: number;
	    totalInteracted: number;
	    byCategory: Record<string, CategoryStats>;
	    byPriority: Record<string, PriorityStats>;
	    byHour: Record<number, number>;
	    byDay: Record<number, number>;
	    // Go type: time
	    periodStart: any;
	    // Go type: time
	    periodEnd: any;
	
	    static createFrom(source: any = {}) {
	        return new NotificationStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalSent = source["totalSent"];
	        this.totalToast = source["totalToast"];
	        this.totalSystem = source["totalSystem"];
	        this.totalRead = source["totalRead"];
	        this.totalDismissed = source["totalDismissed"];
	        this.totalInteracted = source["totalInteracted"];
	        this.byCategory = this.convertValues(source["byCategory"], CategoryStats, true);
	        this.byPriority = this.convertValues(source["byPriority"], PriorityStats, true);
	        this.byHour = source["byHour"];
	        this.byDay = source["byDay"];
	        this.periodStart = this.convertValues(source["periodStart"], null);
	        this.periodEnd = this.convertValues(source["periodEnd"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ThrottleRule {
	    maxPerMinute: number;
	    maxPerHour: number;
	    burstAllowance: number;
	    cooldownPeriod: number;
	
	    static createFrom(source: any = {}) {
	        return new ThrottleRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.maxPerMinute = source["maxPerMinute"];
	        this.maxPerHour = source["maxPerHour"];
	        this.burstAllowance = source["burstAllowance"];
	        this.cooldownPeriod = source["cooldownPeriod"];
	    }
	}
	export class NotificationThrottle {
	    enabled: boolean;
	    maxPerMinute: number;
	    maxPerHour: number;
	    burstAllowance: number;
	    cooldownPeriod: number;
	    byCategory: Record<string, ThrottleRule>;
	    byPriority: Record<string, ThrottleRule>;
	
	    static createFrom(source: any = {}) {
	        return new NotificationThrottle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.maxPerMinute = source["maxPerMinute"];
	        this.maxPerHour = source["maxPerHour"];
	        this.burstAllowance = source["burstAllowance"];
	        this.cooldownPeriod = source["cooldownPeriod"];
	        this.byCategory = this.convertValues(source["byCategory"], ThrottleRule, true);
	        this.byPriority = this.convertValues(source["byPriority"], ThrottleRule, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class NotificationSystem {
	    enabled: boolean;
	    toastEnabled: boolean;
	    systemEnabled: boolean;
	    soundEnabled: boolean;
	    doNotDisturbMode: boolean;
	    doNotDisturbSchedule?: DoNotDisturbSchedule;
	    maxToastNotifications: number;
	    toastDuration: number;
	    historyRetention: number;
	    throttleSettings?: NotificationThrottle;
	    preferences?: NotificationPreferences;
	
	    static createFrom(source: any = {}) {
	        return new NotificationSystem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.toastEnabled = source["toastEnabled"];
	        this.systemEnabled = source["systemEnabled"];
	        this.soundEnabled = source["soundEnabled"];
	        this.doNotDisturbMode = source["doNotDisturbMode"];
	        this.doNotDisturbSchedule = this.convertValues(source["doNotDisturbSchedule"], DoNotDisturbSchedule);
	        this.maxToastNotifications = source["maxToastNotifications"];
	        this.toastDuration = source["toastDuration"];
	        this.historyRetention = source["historyRetention"];
	        this.throttleSettings = this.convertValues(source["throttleSettings"], NotificationThrottle);
	        this.preferences = this.convertValues(source["preferences"], NotificationPreferences);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class PerformanceAlert {
	    id: string;
	    type: string;
	    severity: string;
	    message: string;
	    threshold: Record<string, any>;
	    actualValue: Record<string, any>;
	    templateId?: string;
	    operation?: string;
	    // Go type: time
	    timestamp: any;
	    resolved: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PerformanceAlert(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.severity = source["severity"];
	        this.message = source["message"];
	        this.threshold = source["threshold"];
	        this.actualValue = source["actualValue"];
	        this.templateId = source["templateId"];
	        this.operation = source["operation"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.resolved = source["resolved"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PerformanceMetrics {
	    startup_time: number;
	    memory_usage: number;
	
	    static createFrom(source: any = {}) {
	        return new PerformanceMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.startup_time = source["startup_time"];
	        this.memory_usage = source["memory_usage"];
	    }
	}
	export class PerformanceMonitor {
	
	
	    static createFrom(source: any = {}) {
	        return new PerformanceMonitor(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	
	export class Project {
	    id: string;
	    name: string;
	    specPath: string;
	    specUrl: string;
	    outputPath: string;
	    settings: ProjectSettings;
	    status: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    // Go type: time
	    lastGenerated?: any;
	    generationCount: number;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.specPath = source["specPath"];
	        this.specUrl = source["specUrl"];
	        this.outputPath = source["outputPath"];
	        this.settings = this.convertValues(source["settings"], ProjectSettings);
	        this.status = source["status"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.lastGenerated = this.convertValues(source["lastGenerated"], null);
	        this.generationCount = source["generationCount"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class ProjectUpdateRequest {
	    name?: string;
	    specPath?: string;
	    specUrl?: string;
	    outputPath?: string;
	    settings?: ProjectSettings;
	
	    static createFrom(source: any = {}) {
	        return new ProjectUpdateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.specPath = source["specPath"];
	        this.specUrl = source["specUrl"];
	        this.outputPath = source["outputPath"];
	        this.settings = this.convertValues(source["settings"], ProjectSettings);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RecentFile {
	    path: string;
	    name: string;
	    size: number;
	    lastAccessed: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new RecentFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.name = source["name"];
	        this.size = source["size"];
	        this.lastAccessed = source["lastAccessed"];
	        this.type = source["type"];
	    }
	}
	
	export class RollbackCapabilities {
	    available: boolean;
	    backupCount: number;
	    maxBackups: number;
	    backupDir: string;
	    features: string[];
	
	    static createFrom(source: any = {}) {
	        return new RollbackCapabilities(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.available = source["available"];
	        this.backupCount = source["backupCount"];
	        this.maxBackups = source["maxBackups"];
	        this.backupDir = source["backupDir"];
	        this.features = source["features"];
	    }
	}
	export class RollbackInfo {
	    available: boolean;
	    backupPath: string;
	    backupVersion: string;
	    // Go type: time
	    backupCreatedAt: any;
	    backupSize: number;
	
	    static createFrom(source: any = {}) {
	        return new RollbackInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.available = source["available"];
	        this.backupPath = source["backupPath"];
	        this.backupVersion = source["backupVersion"];
	        this.backupCreatedAt = this.convertValues(source["backupCreatedAt"], null);
	        this.backupSize = source["backupSize"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UpdateInfo {
	    version: string;
	    releaseNotes: string;
	    downloadUrl: string;
	    checksumUrl: string;
	    signatureUrl: string;
	    size: number;
	    // Go type: time
	    publishedAt: any;
	    critical: boolean;
	    metadata?: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.releaseNotes = source["releaseNotes"];
	        this.downloadUrl = source["downloadUrl"];
	        this.checksumUrl = source["checksumUrl"];
	        this.signatureUrl = source["signatureUrl"];
	        this.size = source["size"];
	        this.publishedAt = this.convertValues(source["publishedAt"], null);
	        this.critical = source["critical"];
	        this.metadata = source["metadata"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ScheduledJob {
	    id: string;
	    type: string;
	    schedule?: UpdateSchedule;
	    updateInfo?: UpdateInfo;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    lastRun?: any;
	    // Go type: time
	    nextRun: any;
	    runCount: number;
	    status: string;
	    error?: APIError;
	
	    static createFrom(source: any = {}) {
	        return new ScheduledJob(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.schedule = this.convertValues(source["schedule"], UpdateSchedule);
	        this.updateInfo = this.convertValues(source["updateInfo"], UpdateInfo);
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.lastRun = this.convertValues(source["lastRun"], null);
	        this.nextRun = this.convertValues(source["nextRun"], null);
	        this.runCount = source["runCount"];
	        this.status = source["status"];
	        this.error = this.convertValues(source["error"], APIError);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	export class SystemMetrics {
	    // Go type: time
	    timestamp: any;
	    memoryTotal: number;
	    memoryUsed: number;
	    memoryAvailable: number;
	    cpuUsage: number;
	    goRoutines: number;
	    heapAlloc: number;
	    heapSys: number;
	    numGC: number;
	    activeTemplates: number;
	    cacheSize: number;
	    requestsPerSecond: number;
	
	    static createFrom(source: any = {}) {
	        return new SystemMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.memoryTotal = source["memoryTotal"];
	        this.memoryUsed = source["memoryUsed"];
	        this.memoryAvailable = source["memoryAvailable"];
	        this.cpuUsage = source["cpuUsage"];
	        this.goRoutines = source["goRoutines"];
	        this.heapAlloc = source["heapAlloc"];
	        this.heapSys = source["heapSys"];
	        this.numGC = source["numGC"];
	        this.activeTemplates = source["activeTemplates"];
	        this.cacheSize = source["cacheSize"];
	        this.requestsPerSecond = source["requestsPerSecond"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SystemNotification {
	    id: string;
	    title: string;
	    body: string;
	    icon?: string;
	    sound?: string;
	    actions?: NotificationActionBtn[];
	    urgency: string;
	    tag?: string;
	    // Go type: time
	    createdAt: any;
	    category: string;
	    timeout: number;
	    silent: boolean;
	    metadata?: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new SystemNotification(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.body = source["body"];
	        this.icon = source["icon"];
	        this.sound = source["sound"];
	        this.actions = this.convertValues(source["actions"], NotificationActionBtn);
	        this.urgency = source["urgency"];
	        this.tag = source["tag"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.category = source["category"];
	        this.timeout = source["timeout"];
	        this.silent = source["silent"];
	        this.metadata = source["metadata"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TemplatePerformanceMetrics {
	    templateId: string;
	    operation: string;
	    // Go type: time
	    startTime: any;
	    // Go type: time
	    endTime: any;
	    duration: number;
	    memoryUsage: number;
	    cpuUsage: number;
	    success: boolean;
	    errorMessage?: string;
	    inputSize?: number;
	    outputSize?: number;
	    cacheHit: boolean;
	    complexity: string;
	    variableCount: number;
	    functionCount: number;
	    // Go type: time
	    timestamp: any;
	
	    static createFrom(source: any = {}) {
	        return new TemplatePerformanceMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.templateId = source["templateId"];
	        this.operation = source["operation"];
	        this.startTime = this.convertValues(source["startTime"], null);
	        this.endTime = this.convertValues(source["endTime"], null);
	        this.duration = source["duration"];
	        this.memoryUsage = source["memoryUsage"];
	        this.cpuUsage = source["cpuUsage"];
	        this.success = source["success"];
	        this.errorMessage = source["errorMessage"];
	        this.inputSize = source["inputSize"];
	        this.outputSize = source["outputSize"];
	        this.cacheHit = source["cacheHit"];
	        this.complexity = source["complexity"];
	        this.variableCount = source["variableCount"];
	        this.functionCount = source["functionCount"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class ToastNotification {
	    id: string;
	    type: string;
	    title: string;
	    message: string;
	    icon?: string;
	    duration: number;
	    position: string;
	    actions?: NotificationActionBtn[];
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    expiresAt: any;
	    persistent: boolean;
	    autoDismiss: boolean;
	    priority: string;
	    category: string;
	    metadata?: Record<string, any>;
	    progress?: NotificationProgress;
	
	    static createFrom(source: any = {}) {
	        return new ToastNotification(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.title = source["title"];
	        this.message = source["message"];
	        this.icon = source["icon"];
	        this.duration = source["duration"];
	        this.position = source["position"];
	        this.actions = this.convertValues(source["actions"], NotificationActionBtn);
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.expiresAt = this.convertValues(source["expiresAt"], null);
	        this.persistent = source["persistent"];
	        this.autoDismiss = source["autoDismiss"];
	        this.priority = source["priority"];
	        this.category = source["category"];
	        this.metadata = source["metadata"];
	        this.progress = this.convertValues(source["progress"], NotificationProgress);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UpdateAnalytics {
	    userId?: string;
	    sessionId: string;
	    eventType: string;
	    version: string;
	    previousVersion?: string;
	    updateChannel: string;
	    duration?: number;
	    success: boolean;
	    error?: string;
	    clientInfo: ClientInfo;
	    // Go type: time
	    timestamp: any;
	    size?: number;
	    downloadSpeed?: number;
	    userAction?: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateAnalytics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.userId = source["userId"];
	        this.sessionId = source["sessionId"];
	        this.eventType = source["eventType"];
	        this.version = source["version"];
	        this.previousVersion = source["previousVersion"];
	        this.updateChannel = source["updateChannel"];
	        this.duration = source["duration"];
	        this.success = source["success"];
	        this.error = source["error"];
	        this.clientInfo = this.convertValues(source["clientInfo"], ClientInfo);
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.size = source["size"];
	        this.downloadSpeed = source["downloadSpeed"];
	        this.userAction = source["userAction"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UpdateCheck {
	    id: string;
	    // Go type: time
	    checkedAt: any;
	    success: boolean;
	    updateInfo?: UpdateInfo;
	    error?: APIError;
	    source: string;
	    userAgent: string;
	    clientInfo: ClientInfo;
	
	    static createFrom(source: any = {}) {
	        return new UpdateCheck(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.checkedAt = this.convertValues(source["checkedAt"], null);
	        this.success = source["success"];
	        this.updateInfo = this.convertValues(source["updateInfo"], UpdateInfo);
	        this.error = this.convertValues(source["error"], APIError);
	        this.source = source["source"];
	        this.userAgent = source["userAgent"];
	        this.clientInfo = this.convertValues(source["clientInfo"], ClientInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UpdateConfiguration {
	    updateUrl: string;
	    publicKey: string;
	    certificatePath: string;
	    backupDirectory: string;
	    tempDirectory: string;
	    userAgent: string;
	    timeout: number;
	    verificationMode: string;
	    hashAlgorithm: number;
	    deltaUpdates: boolean;
	    compressionLevel: number;
	    customHeaders?: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new UpdateConfiguration(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.updateUrl = source["updateUrl"];
	        this.publicKey = source["publicKey"];
	        this.certificatePath = source["certificatePath"];
	        this.backupDirectory = source["backupDirectory"];
	        this.tempDirectory = source["tempDirectory"];
	        this.userAgent = source["userAgent"];
	        this.timeout = source["timeout"];
	        this.verificationMode = source["verificationMode"];
	        this.hashAlgorithm = source["hashAlgorithm"];
	        this.deltaUpdates = source["deltaUpdates"];
	        this.compressionLevel = source["compressionLevel"];
	        this.customHeaders = source["customHeaders"];
	    }
	}
	export class UpdateConfigurationValidation {
	    valid: boolean;
	    errors: string[];
	    warnings: string[];
	
	    static createFrom(source: any = {}) {
	        return new UpdateConfigurationValidation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.valid = source["valid"];
	        this.errors = source["errors"];
	        this.warnings = source["warnings"];
	    }
	}
	export class UpdateServerInfo {
	    version: string;
	    status: string;
	    lastUpdated?: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateServerInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.status = source["status"];
	        this.lastUpdated = source["lastUpdated"];
	    }
	}
	export class UpdateConnectionTest {
	    success: boolean;
	    responseTime: number;
	    error?: APIError;
	    serverInfo?: UpdateServerInfo;
	    testedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateConnectionTest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.responseTime = source["responseTime"];
	        this.error = this.convertValues(source["error"], APIError);
	        this.serverInfo = this.convertValues(source["serverInfo"], UpdateServerInfo);
	        this.testedAt = source["testedAt"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class UpdateProgress {
	    status: string;
	    progress: number;
	    currentStep: string;
	    bytesTotal: number;
	    bytesReceived: number;
	    speed: number;
	    estimatedTime?: number;
	    error?: APIError;
	    // Go type: time
	    lastUpdate: any;
	
	    static createFrom(source: any = {}) {
	        return new UpdateProgress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.progress = source["progress"];
	        this.currentStep = source["currentStep"];
	        this.bytesTotal = source["bytesTotal"];
	        this.bytesReceived = source["bytesReceived"];
	        this.speed = source["speed"];
	        this.estimatedTime = source["estimatedTime"];
	        this.error = this.convertValues(source["error"], APIError);
	        this.lastUpdate = this.convertValues(source["lastUpdate"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class VerificationResult {
	    checksumValid: boolean;
	    signatureValid: boolean;
	    certificateValid: boolean;
	    algorithm: string;
	    // Go type: time
	    verifiedAt: any;
	    trustedSource: boolean;
	
	    static createFrom(source: any = {}) {
	        return new VerificationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.checksumValid = source["checksumValid"];
	        this.signatureValid = source["signatureValid"];
	        this.certificateValid = source["certificateValid"];
	        this.algorithm = source["algorithm"];
	        this.verifiedAt = this.convertValues(source["verifiedAt"], null);
	        this.trustedSource = source["trustedSource"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UpdateResult {
	    success: boolean;
	    version: string;
	    previousVersion: string;
	    // Go type: time
	    updatedAt: any;
	    duration: number;
	    error?: APIError;
	    rollbackInfo?: RollbackInfo;
	    verificationResult?: VerificationResult;
	
	    static createFrom(source: any = {}) {
	        return new UpdateResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.version = source["version"];
	        this.previousVersion = source["previousVersion"];
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.duration = source["duration"];
	        this.error = this.convertValues(source["error"], APIError);
	        this.rollbackInfo = this.convertValues(source["rollbackInfo"], RollbackInfo);
	        this.verificationResult = this.convertValues(source["verificationResult"], VerificationResult);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	export class ValidationError {
	    type: string;
	    message: string;
	    path: string;
	    line?: number;
	    column?: number;
	    severity: string;
	    code: string;
	    location?: ErrorLocation;
	
	    static createFrom(source: any = {}) {
	        return new ValidationError(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.message = source["message"];
	        this.path = source["path"];
	        this.line = source["line"];
	        this.column = source["column"];
	        this.severity = source["severity"];
	        this.code = source["code"];
	        this.location = this.convertValues(source["location"], ErrorLocation);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ValidationWarning {
	    type: string;
	    message: string;
	    path: string;
	    suggestion: string;
	
	    static createFrom(source: any = {}) {
	        return new ValidationWarning(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.message = source["message"];
	        this.path = source["path"];
	        this.suggestion = source["suggestion"];
	    }
	}
	export class ValidationResult {
	    valid: boolean;
	    errors: ValidationError[];
	    warnings: ValidationWarning[];
	    suggestions: string[];
	    specInfo?: SpecInfo;
	    validationTime: number;
	    cacheHit: boolean;
	    // Go type: time
	    validatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new ValidationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.valid = source["valid"];
	        this.errors = this.convertValues(source["errors"], ValidationError);
	        this.warnings = this.convertValues(source["warnings"], ValidationWarning);
	        this.suggestions = source["suggestions"];
	        this.specInfo = this.convertValues(source["specInfo"], SpecInfo);
	        this.validationTime = source["validationTime"];
	        this.cacheHit = source["cacheHit"];
	        this.validatedAt = this.convertValues(source["validatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	

}

export namespace database {
	
	export class ValidationCacheStats {
	    totalEntries: number;
	    activeEntries: number;
	    expiredEntries: number;
	
	    static createFrom(source: any = {}) {
	        return new ValidationCacheStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalEntries = source["totalEntries"];
	        this.activeEntries = source["activeEntries"];
	        this.expiredEntries = source["expiredEntries"];
	    }
	}

}

export namespace mapping {
	
	export class Property {
	    type: string;
	    description?: string;
	    example?: any;
	    enum?: string[];
	    format?: string;
	    items?: Property;
	
	    static createFrom(source: any = {}) {
	        return new Property(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.description = source["description"];
	        this.example = source["example"];
	        this.enum = source["enum"];
	        this.format = source["format"];
	        this.items = this.convertValues(source["items"], Property);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class InputSchema {
	    type: string;
	    properties: Record<string, Property>;
	    required: string[];
	
	    static createFrom(source: any = {}) {
	        return new InputSchema(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.properties = this.convertValues(source["properties"], Property, true);
	        this.required = source["required"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MCPTool {
	    name: string;
	    description: string;
	    inputSchema: InputSchema;
	
	    static createFrom(source: any = {}) {
	        return new MCPTool(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.inputSchema = this.convertValues(source["inputSchema"], InputSchema);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace plugin {
	
	export class DependencyAPI {
	    name: string;
	    version: string;
	    type: string;
	    optional: boolean;
	    repository?: string;
	
	    static createFrom(source: any = {}) {
	        return new DependencyAPI(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	        this.type = source["type"];
	        this.optional = source["optional"];
	        this.repository = source["repository"];
	    }
	}
	export class Trial {
	    duration: number;
	    features?: string[];
	
	    static createFrom(source: any = {}) {
	        return new Trial(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.duration = source["duration"];
	        this.features = source["features"];
	    }
	}
	export class Price {
	    amount: number;
	    currency: string;
	    type: string;
	    trial?: Trial;
	
	    static createFrom(source: any = {}) {
	        return new Price(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.amount = source["amount"];
	        this.currency = source["currency"];
	        this.type = source["type"];
	        this.trial = this.convertValues(source["trial"], Trial);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MarketplaceStats {
	    downloads: number;
	    rating: number;
	    reviewCount: number;
	    // Go type: time
	    lastUpdated: any;
	    compatibility: string[];
	
	    static createFrom(source: any = {}) {
	        return new MarketplaceStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.downloads = source["downloads"];
	        this.rating = source["rating"];
	        this.reviewCount = source["reviewCount"];
	        this.lastUpdated = this.convertValues(source["lastUpdated"], null);
	        this.compatibility = source["compatibility"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Review {
	    userId: string;
	    userName: string;
	    rating: number;
	    comment: string;
	    // Go type: time
	    createdAt: any;
	    helpful: number;
	
	    static createFrom(source: any = {}) {
	        return new Review(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.userId = source["userId"];
	        this.userName = source["userName"];
	        this.rating = source["rating"];
	        this.comment = source["comment"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.helpful = source["helpful"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Dependency {
	    name: string;
	    version: string;
	    type: string;
	    optional: boolean;
	    repository?: string;
	
	    static createFrom(source: any = {}) {
	        return new Dependency(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	        this.type = source["type"];
	        this.optional = source["optional"];
	        this.repository = source["repository"];
	    }
	}
	export class PluginConfig {
	    schema: number[];
	    default: number[];
	    required: string[];
	    examples: number[][];
	
	    static createFrom(source: any = {}) {
	        return new PluginConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.schema = source["schema"];
	        this.default = source["default"];
	        this.required = source["required"];
	        this.examples = source["examples"];
	    }
	}
	export class MarketplacePlugin {
	    id: string;
	    name: string;
	    version: string;
	    description: string;
	    author: string;
	    homepage?: string;
	    repository?: string;
	    license: string;
	    tags?: string[];
	    minVersion: string;
	    maxVersion: string;
	    // Go type: PluginConfig
	    config?: any;
	    permissions?: string[];
	    dependencies?: Dependency[];
	    metadata?: Record<string, string>;
	    downloadUrl: string;
	    screenshots?: string[];
	    documentation?: string;
	    reviews?: Review[];
	    stats?: MarketplaceStats;
	    // Go type: time
	    updatedAt: any;
	    featured: boolean;
	    verified: boolean;
	    category: string;
	    price?: Price;
	
	    static createFrom(source: any = {}) {
	        return new MarketplacePlugin(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.version = source["version"];
	        this.description = source["description"];
	        this.author = source["author"];
	        this.homepage = source["homepage"];
	        this.repository = source["repository"];
	        this.license = source["license"];
	        this.tags = source["tags"];
	        this.minVersion = source["minVersion"];
	        this.maxVersion = source["maxVersion"];
	        this.config = this.convertValues(source["config"], null);
	        this.permissions = source["permissions"];
	        this.dependencies = this.convertValues(source["dependencies"], Dependency);
	        this.metadata = source["metadata"];
	        this.downloadUrl = source["downloadUrl"];
	        this.screenshots = source["screenshots"];
	        this.documentation = source["documentation"];
	        this.reviews = this.convertValues(source["reviews"], Review);
	        this.stats = this.convertValues(source["stats"], MarketplaceStats);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.featured = source["featured"];
	        this.verified = source["verified"];
	        this.category = source["category"];
	        this.price = this.convertValues(source["price"], Price);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class PluginConfigAPI {
	    schema: number[];
	    default: number[];
	    required: string[];
	    examples: number[][];
	
	    static createFrom(source: any = {}) {
	        return new PluginConfigAPI(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.schema = source["schema"];
	        this.default = source["default"];
	        this.required = source["required"];
	        this.examples = source["examples"];
	    }
	}
	export class PluginFileAPI {
	    path: string;
	    size: number;
	    checksum: string;
	    type: string;
	    platform?: string;
	    arch?: string;
	
	    static createFrom(source: any = {}) {
	        return new PluginFileAPI(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.size = source["size"];
	        this.checksum = source["checksum"];
	        this.type = source["type"];
	        this.platform = source["platform"];
	        this.arch = source["arch"];
	    }
	}
	export class PluginInfoAPI {
	    id: string;
	    name: string;
	    version: string;
	    description: string;
	    author: string;
	    homepage?: string;
	    repository?: string;
	    license: string;
	    tags?: string[];
	    minVersion: string;
	    maxVersion: string;
	    config?: PluginConfigAPI;
	    permissions?: string[];
	    dependencies?: DependencyAPI[];
	    metadata?: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new PluginInfoAPI(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.version = source["version"];
	        this.description = source["description"];
	        this.author = source["author"];
	        this.homepage = source["homepage"];
	        this.repository = source["repository"];
	        this.license = source["license"];
	        this.tags = source["tags"];
	        this.minVersion = source["minVersion"];
	        this.maxVersion = source["maxVersion"];
	        this.config = this.convertValues(source["config"], PluginConfigAPI);
	        this.permissions = source["permissions"];
	        this.dependencies = this.convertValues(source["dependencies"], DependencyAPI);
	        this.metadata = source["metadata"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PluginManifestAPI {
	    id: string;
	    name: string;
	    version: string;
	    description: string;
	    author: string;
	    homepage?: string;
	    repository?: string;
	    license: string;
	    tags?: string[];
	    minVersion: string;
	    maxVersion: string;
	    config?: PluginConfigAPI;
	    permissions?: string[];
	    dependencies?: DependencyAPI[];
	    metadata?: Record<string, string>;
	    files: PluginFileAPI[];
	    checksum: string;
	    size: number;
	    installPath?: string;
	    verified: boolean;
	    signature?: string;
	
	    static createFrom(source: any = {}) {
	        return new PluginManifestAPI(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.version = source["version"];
	        this.description = source["description"];
	        this.author = source["author"];
	        this.homepage = source["homepage"];
	        this.repository = source["repository"];
	        this.license = source["license"];
	        this.tags = source["tags"];
	        this.minVersion = source["minVersion"];
	        this.maxVersion = source["maxVersion"];
	        this.config = this.convertValues(source["config"], PluginConfigAPI);
	        this.permissions = source["permissions"];
	        this.dependencies = this.convertValues(source["dependencies"], DependencyAPI);
	        this.metadata = source["metadata"];
	        this.files = this.convertValues(source["files"], PluginFileAPI);
	        this.checksum = source["checksum"];
	        this.size = source["size"];
	        this.installPath = source["installPath"];
	        this.verified = source["verified"];
	        this.signature = source["signature"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PluginStatsAPI {
	    callCount: number;
	    totalDuration: number;
	    averageDuration: number;
	    errorCount: number;
	    lastUsed: string;
	    memoryUsage: number;
	
	    static createFrom(source: any = {}) {
	        return new PluginStatsAPI(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.callCount = source["callCount"];
	        this.totalDuration = source["totalDuration"];
	        this.averageDuration = source["averageDuration"];
	        this.errorCount = source["errorCount"];
	        this.lastUsed = source["lastUsed"];
	        this.memoryUsage = source["memoryUsage"];
	    }
	}
	export class PluginInstanceAPI {
	    info?: PluginInfoAPI;
	    status: string;
	    config?: number[];
	    loadedAt: string;
	    lastError?: string;
	    stats?: PluginStatsAPI;
	    manifest?: PluginManifestAPI;
	
	    static createFrom(source: any = {}) {
	        return new PluginInstanceAPI(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.info = this.convertValues(source["info"], PluginInfoAPI);
	        this.status = source["status"];
	        this.config = source["config"];
	        this.loadedAt = source["loadedAt"];
	        this.lastError = source["lastError"];
	        this.stats = this.convertValues(source["stats"], PluginStatsAPI);
	        this.manifest = this.convertValues(source["manifest"], PluginManifestAPI);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	export class SearchResponse {
	    plugins: MarketplacePlugin[];
	    total: number;
	    categories: string[];
	    tags: string[];
	    page: number;
	    perPage: number;
	
	    static createFrom(source: any = {}) {
	        return new SearchResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.plugins = this.convertValues(source["plugins"], MarketplacePlugin);
	        this.total = source["total"];
	        this.categories = source["categories"];
	        this.tags = source["tags"];
	        this.page = source["page"];
	        this.perPage = source["perPage"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

