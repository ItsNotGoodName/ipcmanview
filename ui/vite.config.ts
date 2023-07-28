import { defineConfig } from 'vite'
import { macaronVitePlugin } from '@macaron-css/vite';
import solid from 'vite-plugin-solid'
import path from "path";

export default defineConfig({
  plugins: [
    macaronVitePlugin(),
    solid(),
  ],
  server: {
    port: 3000,
  },
  resolve: {
    alias: {
      "~": path.resolve(__dirname, "./src"),
    },
  },
})
