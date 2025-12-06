import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
        '@': './src'
    }
  },
  build: {
    rolldownOptions: {
        output: {
            dir: '../service/app/dist',
            entryFileNames: 'static/js/[name]-[hash].js',
            chunkFileNames: 'static/js/[name]-[hash].js',
            assetFileNames: 'static/[ext]/[name]-[hash].[ext]'
        }
    }
  }
})
