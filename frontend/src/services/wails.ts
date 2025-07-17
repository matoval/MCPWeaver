// Wails runtime integration
// This provides a bridge between the React frontend and the Go backend

// Type definitions for Wails runtime
declare global {
  interface Window {
    go?: any;
    runtime?: any;
  }
}

// Wails runtime API wrapper
export class WailsRuntime {
  private static instance: WailsRuntime;
  private isReady = false;

  private constructor() {
    this.init();
  }

  public static getInstance(): WailsRuntime {
    if (!WailsRuntime.instance) {
      WailsRuntime.instance = new WailsRuntime();
    }
    return WailsRuntime.instance;
  }

  private init() {
    // Check if running in Wails environment
    if (typeof window !== 'undefined' && window.go) {
      this.isReady = true;
    } else {
      // Mock for development
      console.log('Wails runtime not available - using mock implementation');
      this.setupMockRuntime();
    }
  }

  private setupMockRuntime() {
    // Mock implementation for development
    window.go = {
      main: {
        App: {
          CreateProject: (request: any) => Promise.resolve({ id: 'mock-project' }),
          GetProjects: () => Promise.resolve([]),
          GetProjectById: (id: string) => Promise.resolve(null),
          UpdateProject: (id: string, updates: any) => Promise.resolve(null),
          DeleteProject: (id: string) => Promise.resolve(null),
          ValidateProject: (id: string) => Promise.resolve({ valid: true }),
          GenerateServer: (id: string) => Promise.resolve({ success: true }),
          GetGenerationJobs: () => Promise.resolve([]),
          GetGenerationJobById: (id: string) => Promise.resolve(null),
          CancelGenerationJob: (id: string) => Promise.resolve(null),
          GetSettings: () => Promise.resolve({}),
          UpdateSettings: (settings: any) => Promise.resolve(null),
          GetSystemInfo: () => Promise.resolve({
            version: '1.0.0',
            os: 'mock',
            arch: 'mock'
          }),
          ExitApp: () => Promise.resolve(null)
        }
      }
    };

    window.runtime = {
      EventsEmit: (event: string, data?: any) => {
        console.log('Mock EventsEmit:', event, data);
      },
      EventsOn: (event: string, callback: (data: any) => void) => {
        console.log('Mock EventsOn:', event);
        // Return cleanup function
        return () => {};
      },
      EventsOff: (event: string) => {
        console.log('Mock EventsOff:', event);
      }
    };

    this.isReady = true;
  }

  public isWailsReady(): boolean {
    return this.isReady;
  }

  // Project management
  public async createProject(request: any) {
    if (!this.isReady) throw new Error('Wails runtime not ready');
    return await window.go.main.App.CreateProject(request);
  }

  public async getProjects() {
    if (!this.isReady) throw new Error('Wails runtime not ready');
    return await window.go.main.App.GetProjects();
  }

  public async getProjectById(id: string) {
    if (!this.isReady) throw new Error('Wails runtime not ready');
    return await window.go.main.App.GetProjectById(id);
  }

  public async updateProject(id: string, updates: any) {
    if (!this.isReady) throw new Error('Wails runtime not ready');
    return await window.go.main.App.UpdateProject(id, updates);
  }

  public async deleteProject(id: string) {
    if (!this.isReady) throw new Error('Wails runtime not ready');
    return await window.go.main.App.DeleteProject(id);
  }

  // Validation
  public async validateProject(id: string) {
    if (!this.isReady) throw new Error('Wails runtime not ready');
    return await window.go.main.App.ValidateProject(id);
  }

  // Generation
  public async generateServer(id: string) {
    if (!this.isReady) throw new Error('Wails runtime not ready');
    return await window.go.main.App.GenerateServer(id);
  }

  public async getGenerationJobs() {
    if (!this.isReady) throw new Error('Wails runtime not ready');
    return await window.go.main.App.GetGenerationJobs();
  }

  public async getGenerationJobById(id: string) {
    if (!this.isReady) throw new Error('Wails runtime not ready');
    return await window.go.main.App.GetGenerationJobById(id);
  }

  public async cancelGenerationJob(id: string) {
    if (!this.isReady) throw new Error('Wails runtime not ready');
    return await window.go.main.App.CancelGenerationJob(id);
  }

  // Settings
  public async getSettings() {
    if (!this.isReady) throw new Error('Wails runtime not ready');
    return await window.go.main.App.GetSettings();
  }

  public async updateSettings(settings: any) {
    if (!this.isReady) throw new Error('Wails runtime not ready');
    return await window.go.main.App.UpdateSettings(settings);
  }

  // System
  public async getSystemInfo() {
    if (!this.isReady) throw new Error('Wails runtime not ready');
    return await window.go.main.App.GetSystemInfo();
  }

  public async exitApp() {
    if (!this.isReady) throw new Error('Wails runtime not ready');
    return await window.go.main.App.ExitApp();
  }

  // Events
  public emitEvent(event: string, data?: any) {
    if (!this.isReady) throw new Error('Wails runtime not ready');
    window.runtime.EventsEmit(event, data);
  }

  public onEvent(event: string, callback: (data: any) => void): () => void {
    if (!this.isReady) throw new Error('Wails runtime not ready');
    return window.runtime.EventsOn(event, callback);
  }

  public offEvent(event: string) {
    if (!this.isReady) throw new Error('Wails runtime not ready');
    window.runtime.EventsOff(event);
  }
}

// Export singleton instance
export const wails = WailsRuntime.getInstance();