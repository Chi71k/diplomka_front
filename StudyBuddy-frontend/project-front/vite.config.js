import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss(),],
  server: {
    host: true,
    proxy: {
      // Route users API to users service, other API calls to auth service
      '/api/v1/users': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        secure: false,
      },
      // Proxy remaining API calls (auth) to auth service
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      },
    },
  },
})
