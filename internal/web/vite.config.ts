import { defineConfig } from "vite"
import FullReload from 'vite-plugin-full-reload'

export default defineConfig({
  plugins: [
    FullReload(['views/**/*'])
  ],
  build: {
    // generate manifest.json in outDir
    manifest: true,
    rollupOptions: {
      // overwrite default .html entry
      input: 'src/main.ts',
    },
  },
})
