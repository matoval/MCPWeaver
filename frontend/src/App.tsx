import React from 'react'
import { ThemeProvider } from './contexts/ThemeContext'
import Layout from './components/Layout/Layout'
import Router from './components/Router/Router'
import './styles/variables.scss'
import './styles/themes.scss'
import './App.css'

const App: React.FC = () => {
  return (
    <ThemeProvider>
      <Layout>
        <Router />
      </Layout>
    </ThemeProvider>
  )
}

export default App