# Plugin Management UI Components

This document outlines the React components needed for plugin management in the MCPWeaver frontend.

## Component Structure

### 1. PluginManager Component
Main container component for plugin management functionality.

```typescript
interface PluginManagerProps {
  onPluginChange?: (plugins: PluginInstance[]) => void;
}

export const PluginManager: React.FC<PluginManagerProps> = ({ onPluginChange }) => {
  // Implementation would include tabs for installed, marketplace, and settings
};
```

### 2. InstalledPlugins Component
Displays currently installed plugins with management actions.

```typescript
interface InstalledPluginsProps {
  plugins: PluginInstance[];
  onEnable: (pluginId: string) => void;
  onDisable: (pluginId: string) => void;
  onUninstall: (pluginId: string) => void;
  onConfigure: (pluginId: string) => void;
}

export const InstalledPlugins: React.FC<InstalledPluginsProps> = ({
  plugins,
  onEnable,
  onDisable,
  onUninstall,
  onConfigure
}) => {
  // Implementation would render plugin cards with action buttons
};
```

### 3. PluginMarketplace Component
Browse and install plugins from the marketplace.

```typescript
interface PluginMarketplaceProps {
  onInstall: (pluginId: string) => void;
  onSearch: (query: string, category?: string) => void;
}

export const PluginMarketplace: React.FC<PluginMarketplaceProps> = ({
  onInstall,
  onSearch
}) => {
  // Implementation would include search, categories, and plugin listings
};
```

### 4. PluginCard Component
Individual plugin display component.

```typescript
interface PluginCardProps {
  plugin: PluginInstance | MarketplacePlugin;
  actions: PluginAction[];
  isInstalled?: boolean;
  status?: string;
}

export const PluginCard: React.FC<PluginCardProps> = ({
  plugin,
  actions,
  isInstalled,
  status
}) => {
  // Implementation would show plugin info, status, and available actions
};
```

### 5. PluginConfigDialog Component
Configuration dialog for plugin-specific settings.

```typescript
interface PluginConfigDialogProps {
  plugin: PluginInstance;
  isOpen: boolean;
  onClose: () => void;
  onSave: (config: any) => void;
}

export const PluginConfigDialog: React.FC<PluginConfigDialogProps> = ({
  plugin,
  isOpen,
  onClose,
  onSave
}) => {
  // Implementation would render dynamic form based on plugin config schema
};
```

## Integration with Wails Backend

### API Service
```typescript
import { GetPlugins, LoadPlugin, UnloadPlugin, EnablePlugin, DisablePlugin, InstallPlugin, SearchPlugins } from '../../wailsjs/go/app/App';

export class PluginService {
  static async getPlugins(): Promise<Record<string, PluginInstance>> {
    return await GetPlugins();
  }

  static async loadPlugin(path: string): Promise<void> {
    return await LoadPlugin(path);
  }

  static async unloadPlugin(pluginId: string): Promise<void> {
    return await UnloadPlugin(pluginId);
  }

  static async enablePlugin(pluginId: string): Promise<void> {
    return await EnablePlugin(pluginId);
  }

  static async disablePlugin(pluginId: string): Promise<void> {
    return await DisablePlugin(pluginId);
  }

  static async installPlugin(pluginId: string): Promise<void> {
    return await InstallPlugin(pluginId);
  }

  static async searchPlugins(query: string, category?: string, tags?: string[], limit: number = 20): Promise<SearchResponse> {
    return await SearchPlugins(query, category || '', tags || [], limit);
  }
}
```

### Event Handling
```typescript
import { Events } from '@wailsapp/runtime';

export const usePluginEvents = () => {
  useEffect(() => {
    const unsubscribe = Events.on('plugin:loaded', (plugin: PluginInstance) => {
      // Handle plugin loaded event
    });

    Events.on('plugin:unloaded', (pluginId: string) => {
      // Handle plugin unloaded event
    });

    Events.on('plugin:error', (error: PluginError) => {
      // Handle plugin error event
    });

    return unsubscribe;
  }, []);
};
```

## Plugin Management Features

### 1. Plugin Installation Flow
1. Browse marketplace or upload local plugin
2. Validate plugin manifest and permissions
3. Show permission confirmation dialog
4. Download and install plugin
5. Load and initialize plugin
6. Show success/error feedback

### 2. Plugin Configuration
1. Dynamic form generation based on plugin config schema
2. Validation of configuration values
3. Live preview of configuration changes
4. Save/cancel configuration changes

### 3. Plugin Updates
1. Check for available updates
2. Show update notifications
3. Backup current plugin before update
4. Download and install updates
5. Migrate configuration if needed

### 4. Plugin Permissions
1. Display required permissions clearly
2. Allow users to review and approve permissions
3. Show security warnings for dangerous permissions
4. Option to run plugins in restricted mode

## Styling and UX Guidelines

### Visual Design
- Use consistent styling with main MCPWeaver theme
- Clear status indicators (enabled/disabled/error)
- Progress indicators for long operations
- Toast notifications for actions

### User Experience
- Clear categorization and search functionality
- Bulk operations for multiple plugins
- Keyboard shortcuts for power users
- Responsive design for different screen sizes

### Error Handling
- Clear error messages with actionable advice
- Recovery options for failed operations
- Debug information for developers
- Offline mode support

## Implementation Notes

1. **State Management**: Use React Context or Redux for plugin state management
2. **Caching**: Cache marketplace data with appropriate TTL
3. **Security**: Validate all plugin operations on backend
4. **Performance**: Lazy load plugin details and icons
5. **Accessibility**: Ensure all components are accessible
6. **Testing**: Unit tests for all components and integration tests for workflows

## Future Enhancements

1. **Plugin Development Tools**: Built-in plugin development environment
2. **Plugin Analytics**: Usage statistics and performance metrics
3. **Plugin Reviews**: User ratings and reviews system
4. **Plugin Dependencies**: Automatic dependency resolution
5. **Plugin Themes**: Allow plugins to customize UI appearance