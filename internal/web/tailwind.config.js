/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'media',
  content: [
    "./views/**/*.{html,js,ts}",
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

