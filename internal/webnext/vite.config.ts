import { defineConfig } from 'vite'
import solid from 'vite-plugin-solid'
import path from "path";

export default defineConfig({
  plugins: [solid()],
  // TODO: remove when this UI becomes the default
  base: "/next",
  server: {
    port: 3000,
    proxy: {
      "/": "http://localhost:8080/",
    }
  },
  resolve: {
    alias: {
      "~": path.resolve(__dirname, "./src"),
    },
  },
})
