import React from 'react'
import './App.css'

const App: React.FC = () => {
  return (
    <div className="app">
      <header className="app-header">
        <h1>MCPWeaver</h1>
        <p>OpenAPI to MCP Server Generator</p>
      </header>
      <main className="app-main">
        <div className="welcome-message">
          <h2>Welcome to MCPWeaver</h2>
          <p>Transform your OpenAPI specifications into Model Context Protocol (MCP) servers.</p>
        </div>
      </main>
    </div>
  )
}

export default App