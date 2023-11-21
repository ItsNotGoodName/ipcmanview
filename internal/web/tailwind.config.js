/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'media',
  content: [
    "./views/**/*.html",
    "./node_modules/flowbite/**/*.js"
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('flowbite/plugin')
  ],
}

