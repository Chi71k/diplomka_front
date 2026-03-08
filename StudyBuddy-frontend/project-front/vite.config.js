import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss(),],
  server: {
    host: true,
    proxy: {
      '/api/v1/courses': {
        target: 'http://localhost:8082',
        changeOrigin: true,
        secure: false,
      },
      '/api/v1/users': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        secure: false,
      },
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      },
    },
  },
})
