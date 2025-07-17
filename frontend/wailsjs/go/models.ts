export namespace app {
	
	export class GenerationSettings {
	    defaultTemplate: string;
	    enableValidation: boolean;
	    autoOpenOutput: boolean;
	    showAdvancedOptions: boolean;
	    backupOnGenerate: boolean;
	    customTemplates: string[];
	
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
	    windowSettings: WindowSettings;
	    editorSettings: EditorSettings;
	    generationSettings: GenerationSettings;
	
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
	        this.windowSettings = this.convertValues(source["windowSettings"], WindowSettings);
	        this.editorSettings = this.convertValues(source["editorSettings"], EditorSettings);
	        this.generationSettings = this.convertValues(source["generationSettings"], GenerationSettings);
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

