import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath } from 'node:url'

export default defineConfig(({ command }) => ({
  plugins: [vue()],
  resolve: { alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) } },
  server: { port: 8080 },
  // dev 用根路径（http://localhost:8080/）；build 产物挂 /admin/ 子路径（配合 serve-same-origin.js 同源部署）
  base: command === 'build' ? '/admin/' : '/',
  build: { chunkSizeWarningLimit: 1200, emptyOutDir: false },
}))
