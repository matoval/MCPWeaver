import React, { useState, useEffect, useRef } from 'react';
import MenuBar from './MenuBar';
import Toolbar from './Toolbar';
import Sidebar from './Sidebar';
import StatusBar from './StatusBar';
import useAccessibility from '../../hooks/useAccessibility';
import useResponsive from '../../hooks/useResponsive';
import './Layout.scss';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [statusBarVisible, setStatusBarVisible] = useState(true);
  const [projectLoaded, setProjectLoaded] = useState(false);
  const [projectValidated, setProjectValidated] = useState(false);
  const [generationComplete, setGenerationComplete] = useState(false);
  
  const mainContentRef = useRef<HTMLElement>(null);
  const sidebarRef = useRef<HTMLElement>(null);
  
  const accessibility = useAccessibility({
    announcePageChanges: true,
    manageFocus: true,
    enableKeyboardNavigation: true,
    skipLinks: true
  });

  const responsive = useResponsive();

  // Auto-collapse sidebar on mobile
  useEffect(() => {
    if (responsive.shouldCollapseSidebar()) {
      setSidebarOpen(false);
    }
  }, [responsive.shouldCollapseSidebar]);

  // System health simulation
  const [systemHealth, setSystemHealth] = useState({
    memoryUsage: 45.2,
    cpuUsage: 12.8,
    diskUsage: 68.5,
    status: 'healthy' as const
  });

  const handleMenuAction = (action: string) => {
    switch (action) {
      case 'toggleSidebar':
        setSidebarOpen(!sidebarOpen);
        accessibility.announce(
          sidebarOpen ? 'Sidebar closed' : 'Sidebar opened',
          'polite'
        );
        break;
      case 'toggleStatusBar':
        setStatusBarVisible(!statusBarVisible);
        accessibility.announce(
          statusBarVisible ? 'Status bar hidden' : 'Status bar shown',
          'polite'
        );
        break;
      case 'newProject':
        console.log('New project action');
        accessibility.announce('Creating new project', 'polite');
        break;
      case 'openProject':
        console.log('Open project action');
        setProjectLoaded(true);
        accessibility.announce('Project loaded successfully', 'polite');
        break;
      case 'importSpec':
        console.log('Import spec action');
        accessibility.announce('Importing OpenAPI specification', 'polite');
        break;
      case 'validateSpec':
        console.log('Validate spec action');
        setProjectValidated(true);
        accessibility.announce('Specification validated successfully', 'polite');
        break;
      case 'generateServer':
        console.log('Generate server action');
        setGenerationComplete(true);
        accessibility.announce('MCP server generated successfully', 'polite');
        break;
      case 'exportServer':
        console.log('Export server action');
        accessibility.announce('Exporting server files', 'polite');
        break;
      case 'openSettings':
        console.log('Open settings action');
        accessibility.announce('Opening settings', 'polite');
        break;
      case 'exit':
        console.log('Exit action');
        accessibility.announce('Exiting application', 'polite');
        break;
      default:
        console.log(`Unhandled action: ${action}`);
    }
  };

  const handleToolbarAction = (action: string) => {
    handleMenuAction(action);
  };

  const handleStatusClick = () => {
    console.log('Status clicked');
  };

  // Simulate system health updates
  useEffect(() => {
    const interval = setInterval(() => {
      setSystemHealth(prev => ({
        ...prev,
        memoryUsage: Math.max(20, Math.min(80, prev.memoryUsage + (Math.random() - 0.5) * 5)),
        cpuUsage: Math.max(5, Math.min(95, prev.cpuUsage + (Math.random() - 0.5) * 10)),
        status: prev.memoryUsage > 75 || prev.cpuUsage > 80 ? 'warning' : 'healthy'
      }));
    }, 5000);

    return () => clearInterval(interval);
  }, []);

  // Handle keyboard shortcuts
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.ctrlKey) {
        switch (event.key) {
          case 'n':
            event.preventDefault();
            handleMenuAction('newProject');
            break;
          case 'o':
            event.preventDefault();
            handleMenuAction('openProject');
            break;
          case 'i':
            event.preventDefault();
            handleMenuAction('importSpec');
            break;
          case 'e':
            event.preventDefault();
            handleMenuAction('exportServer');
            break;
          case 'b':
            event.preventDefault();
            handleMenuAction('toggleSidebar');
            break;
          case ',':
            event.preventDefault();
            handleMenuAction('openSettings');
            break;
          case 'q':
            event.preventDefault();
            handleMenuAction('exit');
            break;
        }
      } else {
        switch (event.key) {
          case 'F5':
            event.preventDefault();
            handleMenuAction('validateSpec');
            break;
          case 'F6':
            event.preventDefault();
            handleMenuAction('generateServer');
            break;
          case 'F7':
            event.preventDefault();
            handleMenuAction('testServer');
            break;
        }
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, []);

  // Set up focus management and landmarks on mount
  useEffect(() => {
    if (mainContentRef.current) {
      accessibility.createLandmark(mainContentRef.current, 'main', 'Main content area');
      mainContentRef.current.id = 'main-content';
    }
    
    // Create skip link
    const skipLink = accessibility.createSkipLink('main-content', 'Skip to main content');
    document.body.insertBefore(skipLink, document.body.firstChild);
    
    // Register focus groups
    const unregisterMain = accessibility.registerFocusGroup('main', mainContentRef.current || undefined);
    const unregisterSidebar = accessibility.registerFocusGroup('sidebar', sidebarRef.current || undefined);
    
    return () => {
      unregisterMain();
      unregisterSidebar();
      if (skipLink.parentNode) {
        skipLink.parentNode.removeChild(skipLink);
      }
    };
  }, [accessibility]);

  const layoutClasses = responsive.getResponsiveClasses({
    default: 'app-layout',
    mobile: 'mobile-layout',
    tablet: 'tablet-layout',
    desktop: 'desktop-layout'
  });

  const toolbarVisible = responsive.useBreakpointValue({
    mobile: false,
    tablet: true,
    desktop: true,
    large: true,
    default: true
  });

  return (
    <div className={layoutClasses} role="application" aria-label="MCPWeaver Desktop Application">
      <header role="banner">
        <MenuBar onMenuAction={handleMenuAction} />
        
        {toolbarVisible && (
          <Toolbar 
            onAction={handleToolbarAction}
            projectLoaded={projectLoaded}
            projectValidated={projectValidated}
            generationComplete={generationComplete}
          />
        )}
      </header>
      
      <div className="layout-body">
        <aside 
          ref={sidebarRef}
          role="complementary"
          aria-label="Project navigation"
          aria-expanded={sidebarOpen}
          aria-hidden={!sidebarOpen}
          className={`sidebar ${sidebarOpen ? 'open' : 'closed'} ${responsive.isMobile ? 'mobile-sidebar' : ''}`}
        >
          <Sidebar 
            isOpen={sidebarOpen} 
            onToggle={() => setSidebarOpen(!sidebarOpen)} 
          />
        </aside>
        
        <main 
          ref={mainContentRef}
          className={`main-content ${sidebarOpen ? 'sidebar-open' : 'sidebar-closed'}`}
          tabIndex={-1}
          aria-label="Main application content"
        >
          {children}
        </main>
      </div>
      
      {statusBarVisible && (
        <footer role="contentinfo" aria-label="Application status">
          <StatusBar 
            status="ready"
            activeOperations={0}
            systemHealth={systemHealth}
            onStatusClick={handleStatusClick}
          />
        </footer>
      )}
      
      {/* Mobile backdrop for sidebar */}
      {responsive.isMobile && sidebarOpen && (
        <div 
          className="mobile-backdrop"
          onClick={() => setSidebarOpen(false)}
          aria-hidden="true"
        />
      )}
    </div>
  );
};

export default Layout;