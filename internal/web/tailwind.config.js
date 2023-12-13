/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'media',
  content: [
    "./views/**/*.html",
    "./src/**/*.{ts,js,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('daisyui'),
  ],
  daisyui: {
    themes: ["light", "dark"],
  }
}

