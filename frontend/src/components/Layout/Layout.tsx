import React, { useState, useEffect } from 'react';
import MenuBar from './MenuBar';
import Toolbar from './Toolbar';
import Sidebar from './Sidebar';
import StatusBar from './StatusBar';
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
        break;
      case 'toggleStatusBar':
        setStatusBarVisible(!statusBarVisible);
        break;
      case 'newProject':
        console.log('New project action');
        break;
      case 'openProject':
        console.log('Open project action');
        setProjectLoaded(true);
        break;
      case 'importSpec':
        console.log('Import spec action');
        break;
      case 'validateSpec':
        console.log('Validate spec action');
        setProjectValidated(true);
        break;
      case 'generateServer':
        console.log('Generate server action');
        setGenerationComplete(true);
        break;
      case 'exportServer':
        console.log('Export server action');
        break;
      case 'openSettings':
        console.log('Open settings action');
        break;
      case 'exit':
        console.log('Exit action');
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

  return (
    <div className="app-layout">
      <MenuBar onMenuAction={handleMenuAction} />
      
      <Toolbar 
        onAction={handleToolbarAction}
        projectLoaded={projectLoaded}
        projectValidated={projectValidated}
        generationComplete={generationComplete}
      />
      
      <div className="layout-body">
        <Sidebar 
          isOpen={sidebarOpen} 
          onToggle={() => setSidebarOpen(!sidebarOpen)} 
        />
        
        <main className={`main-content ${sidebarOpen ? 'sidebar-open' : 'sidebar-closed'}`}>
          {children}
        </main>
      </div>
      
      {statusBarVisible && (
        <StatusBar 
          status="ready"
          activeOperations={0}
          systemHealth={systemHealth}
          onStatusClick={handleStatusClick}
        />
      )}
    </div>
  );
};

export default Layout;