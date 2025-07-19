import React, { useState, useEffect } from 'react'
import { ThemeProvider } from './contexts/ThemeContext'
import { ContextMenuProvider } from './contexts/ContextMenuContext'
import Layout from './components/Layout/Layout'
import Router from './components/Router/Router'
import KeyboardShortcutsDialog from './components/ui/KeyboardShortcutsDialog'
import useKeyboardShortcuts from './hooks/useKeyboardShortcuts'
import './styles/variables.scss'
import './styles/themes.scss'
import './App.css'

const App: React.FC = () => {
  const [showShortcuts, setShowShortcuts] = useState(false)

  // Initialize keyboard shortcuts
  useKeyboardShortcuts()

  useEffect(() => {
    const handleShowShortcuts = () => setShowShortcuts(true)
    const handleShowHelp = () => setShowShortcuts(true)

    window.addEventListener('keyboard:show-shortcuts', handleShowShortcuts)
    window.addEventListener('keyboard:show-help', handleShowHelp)

    return () => {
      window.removeEventListener('keyboard:show-shortcuts', handleShowShortcuts)
      window.removeEventListener('keyboard:show-help', handleShowHelp)
    }
  }, [])

  return (
    <ThemeProvider>
      <ContextMenuProvider>
        <Layout>
          <Router />
        </Layout>
        
        <KeyboardShortcutsDialog
          isOpen={showShortcuts}
          onClose={() => setShowShortcuts(false)}
        />
      </ContextMenuProvider>
    </ThemeProvider>
  )
}

export default App