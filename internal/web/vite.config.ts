import { defineConfig } from 'vite'
import solid from 'vite-plugin-solid'
import path from "path";

export default defineConfig({
  plugins: [solid()],
  resolve: {
    alias: {
      "~": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    port: 5173,
    strictPort: true
  }
})
