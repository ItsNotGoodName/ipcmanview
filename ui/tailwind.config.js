const { createThemes } = require('tw-colors');

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./components/**/*.{js,vue,ts}",
    "./layouts/**/*.vue",
    "./pages/**/*.vue",
    "./plugins/**/*.{js,ts}",
    "./nuxt.config.{js,ts}",
  ],
  theme: {
    colors: {
      inherit: "inherit",
      current: "currentColor",
      transparent: "transparent"
    }
  },
  plugins: [
    createThemes({
      light: {
        Rosewater: "#dc8a78", // hsl(11, 59%, 67%)
        Flamingo: "#dd7878", // hsl(0, 60%, 67%)
        Pink: "#ea76cb", // hsl(316, 73%, 69%)
        Mauve: {
          DEFAULT: "#8839ef", // hsl(266, 85%, 58%)
          200: "hsl(266, 85%, 53%)",
        },
        Red: {
          DEFAULT: "#d20f39", // hsl(347, 87%, 44%)
          200: "hsl(347, 87%, 39%)",
        },
        Maroon: "#e64553", // hsl(355, 76%, 59%)
        Peach: "#fe640b", // hsl(22, 99%, 52%)
        Yellow: {
          DEFAULT: "#df8e1d", // hsl(35, 77%, 49%)
          200: "hsl(35, 77%, 44%)",
        },
        Green: {
          DEFAULT: "#40a02b", // hsl(109, 58%, 40%)
          200: "hsl(109, 58%, 35%)",
        },
        Teal: "#179299", // hsl(183, 74%, 35%)
        Sky: "#04a5e5", // hsl(197, 97%, 46%)
        Sapphire: "#209fb5", // hsl(189, 70%, 42%)
        Blue: {
          DEFAULT: "#1e66f5", // hsl(220, 91%, 54%)
          200: "hsl(220, 91%, 49%)",
        },
        Lavender: "#7287fd", // hsl(231, 97%, 72%)
        Text: "#4c4f69", // hsl(234, 16%, 35%)
        Subtext: {
          100: "#5c5f77", // hsl(233, 13%, 41%)
          DEFAULT: "#6c6f85", // hsl(233, 10%, 47%)
        },
        Overlay: {
          200: "#7c7f93", // hsl(232, 10%, 53%)
          100: "#8c8fa1", // hsl(231, 10%, 59%)
          DEFAULT: "#9ca0b0", // hsl(228, 11%, 65%)
        },
        Surface: {
          200: "#acb0be", // hsl(227, 12%, 71%)
          100: "#bcc0cc", // hsl(225, 14%, 77%)
          DEFAULT: "#ccd0da", // hsl(223, 16%, 83%)
        },
        Base: "#eff1f5", // hsl(220, 23%, 95%)
        Mantle: "#e6e9ef", // hsl(220, 22%, 92%)
        Crust: "#dce0e8", // hsl(220, 21%, 89%)
      },
      dark: {
        Rosewater: "#f5e0dc", // hsl(10, 56%, 91%)
        Flamingo: "#f2cdcd", // hsl(0, 59%, 88%)
        Pink: "#f5c2e7", // hsl(316, 72%, 86%)
        Mauve: {
          DEFAULT: "#cba6f7", // hsl(267, 84%, 81%)
          200: "hsl(267, 84%, 86%)",
        },
        Red: {
          DEFAULT: "#f38ba8", // hsl(343, 81%, 75%)
          200: "hsl(343, 81%, 80%)",
        },
        Maroon: "#eba0ac", // hsl(350, 65%, 77%)
        Peach: "#fab387", // hsl(23, 92%, 75%)
        Yellow: {
          DEFAULT: "#f9e2af", // hsl(41, 86%, 83%)
          200: "hsl(41, 86%, 88%)",
        },
        Green: {
          DEFAULT: "#a6e3a1", // hsl(115, 54%, 76%)
          200: "hsl(115, 54%, 81%)",
        },
        Teal: "#94e2d5", // hsl(170, 57%, 73%)
        Sky: "#89dceb", // hsl(189, 71%, 73%)
        Sapphire: "#74c7ec", // hsl(199, 76%, 69%)
        Blue: {
          DEFAULT: "#89b4fa", // hsl(217, 92%, 76%)
          200: "hsl(217, 92%, 81%)",
        },
        Lavender: "#b4befe", // hsl(232, 97%, 85%)
        Text: "#cdd6f4", // hsl(226, 64%, 88%)
        Subtext: {
          100: "#bac2de", // hsl(227, 35%, 80%)
          DEFAULT: "#a6adc8", // hsl(228, 24%, 72%)
        },
        Overlay: {
          200: "#9399b2", // hsl(228, 17%, 64%)
          100: "#7f849c", // hsl(230, 13%, 55%)
          DEFAULT: "#6c7086", // hsl(231, 11%, 47%)
        },
        Surface: {
          200: "#585b70", // hsl(233, 12%, 39%)
          100: "#45475a", // hsl(234, 13%, 31%)
          DEFAULT: "#313244", // hsl(237, 16%, 23%)
        },
        Base: "#1e1e2e", // hsl(240, 21%, 15%)
        Mantle: "#181825", // hsl(240, 21%, 12%)
        Crust: "#11111b", // hsl(240, 23%, 9%)
      }
    }),
  ],
}
