import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath } from 'node:url'
export default defineConfig({
  plugins: [vue()],
  resolve: { alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) } },
  server: { port: 3000 },
  build: { emptyOutDir: false },
  test: { environment: 'happy-dom', globals: true, include: ['src/**/*.test.ts'] },
})
