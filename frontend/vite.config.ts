import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    outDir: 'dist',
    emptyOutDir: true,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom'],
          router: ['react-router-dom'],
          ui: ['lucide-react']
        }
      }
    },
    minify: 'terser',
    target: 'es2020',
    chunkSizeWarningLimit: 1000
  },
  server: {
    port: 3000,
  },
  optimizeDeps: {
    include: ['react', 'react-dom'],
    exclude: ['@wailsapp/runtime']
  },
  define: {
    __DEV__: JSON.stringify(false)
  }
})