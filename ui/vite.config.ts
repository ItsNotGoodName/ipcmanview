import { defineConfig } from 'vite'
import solid from 'vite-plugin-solid'
import path from "path";

export default defineConfig({
  plugins: [solid()],
  server: {
    port: 3000,
  },
  resolve: {
    alias: {
      "~": path.resolve(__dirname, "./src"),
    },
  },
})
