import React, { useState, useEffect } from 'react';
import { 
  AppSettings, 
  WindowSettings, 
  EditorSettings, 
  GenerationSettings, 
  NotificationSettings, 
  AppearanceSettings 
} from '../../types';
import { GetSettings, UpdateSettings, ResetSettings, ExportSettings, ImportSettings } from '../../../wailsjs/go/app/App';
import './Settings.scss';

interface SettingsTabProps {
  activeTab: string;
  setActiveTab: (tab: string) => void;
}

const SettingsTabs: React.FC<SettingsTabProps> = ({ activeTab, setActiveTab }) => {
  const tabs = [
    { id: 'general', label: 'General', icon: '⚙️' },
    { id: 'appearance', label: 'Appearance', icon: '🎨' },
    { id: 'editor', label: 'Editor', icon: '📝' },
    { id: 'generation', label: 'Generation', icon: '🔧' },
    { id: 'notifications', label: 'Notifications', icon: '🔔' },
    { id: 'advanced', label: 'Advanced', icon: '⚡' }
  ];

  return (
    <div className="settings-tabs">
      {tabs.map(tab => (
        <button
          key={tab.id}
          className={`settings-tab ${activeTab === tab.id ? 'active' : ''}`}
          onClick={() => setActiveTab(tab.id)}
        >
          <span className="tab-icon">{tab.icon}</span>
          <span className="tab-label">{tab.label}</span>
        </button>
      ))}
    </div>
  );
};

interface GeneralSettingsProps {
  settings: AppSettings;
  onUpdate: (settings: Partial<AppSettings>) => void;
}

const GeneralSettings: React.FC<GeneralSettingsProps> = ({ settings, onUpdate }) => {
  return (
    <div className="settings-section">
      <h3>General Settings</h3>
      
      <div className="settings-group">
        <label>Theme</label>
        <select 
          value={settings.theme} 
          onChange={(e) => onUpdate({ theme: e.target.value })}
        >
          <option value="light">Light</option>
          <option value="dark">Dark</option>
          <option value="auto">System</option>
        </select>
      </div>

      <div className="settings-group">
        <label>Language</label>
        <select 
          value={settings.language} 
          onChange={(e) => onUpdate({ language: e.target.value })}
        >
          <option value="en">English</option>
          <option value="es">Español</option>
          <option value="fr">Français</option>
          <option value="de">Deutsch</option>
        </select>
      </div>

      <div className="settings-group">
        <label>Default Output Path</label>
        <input 
          type="text" 
          value={settings.defaultOutputPath} 
          onChange={(e) => onUpdate({ defaultOutputPath: e.target.value })}
          placeholder="./output"
        />
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.autoSave} 
            onChange={(e) => onUpdate({ autoSave: e.target.checked })}
          />
          Enable auto-save
        </label>
      </div>
    </div>
  );
};

interface AppearanceSettingsProps {
  settings: AppearanceSettings;
  onUpdate: (settings: Partial<AppearanceSettings>) => void;
}

const AppearanceSettingsTab: React.FC<AppearanceSettingsProps> = ({ settings, onUpdate }) => {
  return (
    <div className="settings-section">
      <h3>Appearance Settings</h3>
      
      <div className="settings-group">
        <label>UI Theme</label>
        <select 
          value={settings.uiTheme} 
          onChange={(e) => onUpdate({ uiTheme: e.target.value })}
        >
          <option value="system">System</option>
          <option value="light">Light</option>
          <option value="dark">Dark</option>
          <option value="high-contrast">High Contrast</option>
        </select>
      </div>

      <div className="settings-group">
        <label>Accent Color</label>
        <input 
          type="color" 
          value={settings.accentColor} 
          onChange={(e) => onUpdate({ accentColor: e.target.value })}
        />
      </div>

      <div className="settings-group">
        <label>Window Opacity</label>
        <input 
          type="range" 
          min="0.1" 
          max="1" 
          step="0.1" 
          value={settings.windowOpacity} 
          onChange={(e) => onUpdate({ windowOpacity: parseFloat(e.target.value) })}
        />
        <span>{Math.round(settings.windowOpacity * 100)}%</span>
      </div>

      <div className="settings-group">
        <label>Font Scale</label>
        <input 
          type="range" 
          min="0.5" 
          max="2.0" 
          step="0.1" 
          value={settings.fontScale} 
          onChange={(e) => onUpdate({ fontScale: parseFloat(e.target.value) })}
        />
        <span>{settings.fontScale}x</span>
      </div>

      <div className="settings-group">
        <label>Sidebar Position</label>
        <select 
          value={settings.sidebarPosition} 
          onChange={(e) => onUpdate({ sidebarPosition: e.target.value })}
        >
          <option value="left">Left</option>
          <option value="right">Right</option>
        </select>
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.showAnimation} 
            onChange={(e) => onUpdate({ showAnimation: e.target.checked })}
          />
          Enable animations
        </label>
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.reducedMotion} 
            onChange={(e) => onUpdate({ reducedMotion: e.target.checked })}
          />
          Reduced motion (accessibility)
        </label>
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.compactMode} 
            onChange={(e) => onUpdate({ compactMode: e.target.checked })}
          />
          Compact mode
        </label>
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.showSidebar} 
            onChange={(e) => onUpdate({ showSidebar: e.target.checked })}
          />
          Show sidebar
        </label>
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.showStatusBar} 
            onChange={(e) => onUpdate({ showStatusBar: e.target.checked })}
          />
          Show status bar
        </label>
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.showToolbar} 
            onChange={(e) => onUpdate({ showToolbar: e.target.checked })}
          />
          Show toolbar
        </label>
      </div>
    </div>
  );
};

interface EditorSettingsProps {
  settings: EditorSettings;
  onUpdate: (settings: Partial<EditorSettings>) => void;
}

const EditorSettingsTab: React.FC<EditorSettingsProps> = ({ settings, onUpdate }) => {
  return (
    <div className="settings-section">
      <h3>Editor Settings</h3>
      
      <div className="settings-group">
        <label>Font Family</label>
        <select 
          value={settings.fontFamily} 
          onChange={(e) => onUpdate({ fontFamily: e.target.value })}
        >
          <option value="Monaco">Monaco</option>
          <option value="Fira Code">Fira Code</option>
          <option value="Source Code Pro">Source Code Pro</option>
          <option value="Consolas">Consolas</option>
          <option value="monospace">System Monospace</option>
        </select>
      </div>

      <div className="settings-group">
        <label>Font Size</label>
        <input 
          type="number" 
          min="8" 
          max="24" 
          value={settings.fontSize} 
          onChange={(e) => onUpdate({ fontSize: parseInt(e.target.value) })}
        />
        <span>px</span>
      </div>

      <div className="settings-group">
        <label>Tab Size</label>
        <input 
          type="number" 
          min="2" 
          max="8" 
          value={settings.tabSize} 
          onChange={(e) => onUpdate({ tabSize: parseInt(e.target.value) })}
        />
        <span>spaces</span>
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.wordWrap} 
            onChange={(e) => onUpdate({ wordWrap: e.target.checked })}
          />
          Word wrap
        </label>
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.lineNumbers} 
            onChange={(e) => onUpdate({ lineNumbers: e.target.checked })}
          />
          Show line numbers
        </label>
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.syntaxHighlight} 
            onChange={(e) => onUpdate({ syntaxHighlight: e.target.checked })}
          />
          Syntax highlighting
        </label>
      </div>
    </div>
  );
};

interface GenerationSettingsProps {
  settings: GenerationSettings;
  onUpdate: (settings: Partial<GenerationSettings>) => void;
}

const GenerationSettingsTab: React.FC<GenerationSettingsProps> = ({ settings, onUpdate }) => {
  return (
    <div className="settings-section">
      <h3>Generation Settings</h3>
      
      <div className="settings-group">
        <label>Default Template</label>
        <select 
          value={settings.defaultTemplate} 
          onChange={(e) => onUpdate({ defaultTemplate: e.target.value })}
        >
          <option value="default">Default</option>
          <option value="minimal">Minimal</option>
          <option value="complete">Complete</option>
          <option value="custom">Custom</option>
        </select>
      </div>

      <div className="settings-group">
        <label>Max Workers</label>
        <input 
          type="number" 
          min="1" 
          max="16" 
          value={settings.maxWorkers} 
          onChange={(e) => onUpdate({ maxWorkers: parseInt(e.target.value) })}
        />
        <span>threads</span>
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.enableValidation} 
            onChange={(e) => onUpdate({ enableValidation: e.target.checked })}
          />
          Enable validation
        </label>
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.autoOpenOutput} 
            onChange={(e) => onUpdate({ autoOpenOutput: e.target.checked })}
          />
          Auto-open output folder
        </label>
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.showAdvancedOptions} 
            onChange={(e) => onUpdate({ showAdvancedOptions: e.target.checked })}
          />
          Show advanced options
        </label>
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.backupOnGenerate} 
            onChange={(e) => onUpdate({ backupOnGenerate: e.target.checked })}
          />
          Backup on generate
        </label>
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.performanceMode} 
            onChange={(e) => onUpdate({ performanceMode: e.target.checked })}
          />
          Performance mode
        </label>
      </div>
    </div>
  );
};

interface NotificationSettingsProps {
  settings: NotificationSettings;
  onUpdate: (settings: Partial<NotificationSettings>) => void;
}

const NotificationSettingsTab: React.FC<NotificationSettingsProps> = ({ settings, onUpdate }) => {
  return (
    <div className="settings-section">
      <h3>Notification Settings</h3>
      
      <div className="settings-group">
        <label>Notification Position</label>
        <select 
          value={settings.notificationPosition} 
          onChange={(e) => onUpdate({ notificationPosition: e.target.value })}
        >
          <option value="top-right">Top Right</option>
          <option value="top-left">Top Left</option>
          <option value="bottom-right">Bottom Right</option>
          <option value="bottom-left">Bottom Left</option>
          <option value="center">Center</option>
        </select>
      </div>

      <div className="settings-group">
        <label>Notification Duration</label>
        <input 
          type="number" 
          min="1000" 
          max="10000" 
          step="500" 
          value={settings.notificationDuration} 
          onChange={(e) => onUpdate({ notificationDuration: parseInt(e.target.value) })}
        />
        <span>ms</span>
      </div>

      <div className="settings-group">
        <label>Sound Volume</label>
        <input 
          type="range" 
          min="0" 
          max="1" 
          step="0.1" 
          value={settings.soundVolume} 
          onChange={(e) => onUpdate({ soundVolume: parseFloat(e.target.value) })}
        />
        <span>{Math.round(settings.soundVolume * 100)}%</span>
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.enableDesktopNotifications} 
            onChange={(e) => onUpdate({ enableDesktopNotifications: e.target.checked })}
          />
          Enable desktop notifications
        </label>
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.enableSoundNotifications} 
            onChange={(e) => onUpdate({ enableSoundNotifications: e.target.checked })}
          />
          Enable sound notifications
        </label>
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.showGenerationProgress} 
            onChange={(e) => onUpdate({ showGenerationProgress: e.target.checked })}
          />
          Show generation progress
        </label>
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.showErrorNotifications} 
            onChange={(e) => onUpdate({ showErrorNotifications: e.target.checked })}
          />
          Show error notifications
        </label>
      </div>

      <div className="settings-group">
        <label className="checkbox-label">
          <input 
            type="checkbox" 
            checked={settings.showSuccessNotifications} 
            onChange={(e) => onUpdate({ showSuccessNotifications: e.target.checked })}
          />
          Show success notifications
        </label>
      </div>
    </div>
  );
};

interface AdvancedSettingsProps {
  settings: AppSettings;
  onReset: () => void;
  onExport: () => void;
  onImport: () => void;
}

const AdvancedSettings: React.FC<AdvancedSettingsProps> = ({ settings, onReset, onExport, onImport }) => {
  return (
    <div className="settings-section">
      <h3>Advanced Settings</h3>
      
      <div className="settings-group">
        <label>Recent Projects</label>
        <div className="settings-info">
          <span>{settings.recentProjects.length} projects in history</span>
        </div>
      </div>

      <div className="settings-group">
        <label>Recent Files</label>
        <div className="settings-info">
          <span>{settings.recentFiles.length} files in history</span>
        </div>
      </div>

      <div className="settings-actions">
        <button className="btn btn-secondary" onClick={onExport}>
          Export Settings
        </button>
        <button className="btn btn-secondary" onClick={onImport}>
          Import Settings
        </button>
        <button className="btn btn-danger" onClick={onReset}>
          Reset to Defaults
        </button>
      </div>
    </div>
  );
};

const Settings: React.FC = () => {
  const [activeTab, setActiveTab] = useState('general');
  const [settings, setSettings] = useState<AppSettings | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadSettings();
  }, []);

  const loadSettings = async () => {
    try {
      setLoading(true);
      const currentSettings = await GetSettings();
      setSettings(currentSettings);
      setError(null);
    } catch (err) {
      setError('Failed to load settings');
      console.error('Error loading settings:', err);
    } finally {
      setLoading(false);
    }
  };

  const updateSettings = async (updates: Partial<AppSettings>) => {
    if (!settings) return;

    try {
      const updatedSettings = { ...settings, ...updates };
      await UpdateSettings(updatedSettings);
      setSettings(updatedSettings);
    } catch (err) {
      setError('Failed to update settings');
      console.error('Error updating settings:', err);
    }
  };

  const resetSettings = async () => {
    if (!confirm('Are you sure you want to reset all settings to defaults?')) {
      return;
    }

    try {
      await ResetSettings();
      await loadSettings();
    } catch (err) {
      setError('Failed to reset settings');
      console.error('Error resetting settings:', err);
    }
  };

  const exportSettings = async () => {
    try {
      await ExportSettings();
    } catch (err) {
      setError('Failed to export settings');
      console.error('Error exporting settings:', err);
    }
  };

  const importSettings = async () => {
    try {
      await ImportSettings();
      await loadSettings();
    } catch (err) {
      setError('Failed to import settings');
      console.error('Error importing settings:', err);
    }
  };

  if (loading) {
    return <div className="settings-loading">Loading settings...</div>;
  }

  if (error) {
    return <div className="settings-error">Error: {error}</div>;
  }

  if (!settings) {
    return <div className="settings-error">Settings not available</div>;
  }

  const renderTabContent = () => {
    switch (activeTab) {
      case 'general':
        return <GeneralSettings settings={settings} onUpdate={updateSettings} />;
      case 'appearance':
        return <AppearanceSettingsTab settings={settings.appearanceSettings} onUpdate={(updates) => updateSettings({ appearanceSettings: { ...settings.appearanceSettings, ...updates } })} />;
      case 'editor':
        return <EditorSettingsTab settings={settings.editorSettings} onUpdate={(updates) => updateSettings({ editorSettings: { ...settings.editorSettings, ...updates } })} />;
      case 'generation':
        return <GenerationSettingsTab settings={settings.generationSettings} onUpdate={(updates) => updateSettings({ generationSettings: { ...settings.generationSettings, ...updates } })} />;
      case 'notifications':
        return <NotificationSettingsTab settings={settings.notificationSettings} onUpdate={(updates) => updateSettings({ notificationSettings: { ...settings.notificationSettings, ...updates } })} />;
      case 'advanced':
        return <AdvancedSettings settings={settings} onReset={resetSettings} onExport={exportSettings} onImport={importSettings} />;
      default:
        return <GeneralSettings settings={settings} onUpdate={updateSettings} />;
    }
  };

  return (
    <div className="settings-container">
      <div className="settings-header">
        <h1>Settings</h1>
        <p>Configure MCPWeaver to suit your preferences</p>
      </div>
      
      <div className="settings-content">
        <SettingsTabs activeTab={activeTab} setActiveTab={setActiveTab} />
        <div className="settings-panel">
          {renderTabContent()}
        </div>
      </div>
    </div>
  );
};

export default Settings;