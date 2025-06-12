import react from '@vitejs/plugin-react'
import path from 'path'
import { defineConfig } from 'vite'
import { TanStackRouterVite } from '@tanstack/router-plugin/vite'

export default defineConfig({
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src')
    }
  },
  plugins: [
    TanStackRouterVite({ target: 'react', autoCodeSplitting: true }),
    react(),
  ],
})