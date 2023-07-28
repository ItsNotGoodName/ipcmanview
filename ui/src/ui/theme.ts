import { createTheme } from "@macaron-css/core"

const space = {
  "0": "0px",
  px: "1px",
  "0.5": "0.125rem",
  "1": "0.25rem",
  "1.5": "0.375rem",
  "2": "0.5rem",
  "2.5": "0.625rem",
  "3": "0.75rem",
  "3.5": "0.875rem",
  "4": "1rem",
  "5": "1.25rem",
  "6": "1.5rem",
  "7": "1.75rem",
  "8": "2rem",
  "9": "2.25rem",
  "10": "2.5rem",
  "11": "2.75rem",
  "12": "3rem",
  "14": "3.5rem",
  "16": "4rem",
  "20": "5rem",
  "24": "6rem",
  "28": "7rem",
  "32": "8rem",
  "36": "9rem",
  "40": "10rem",
  "44": "11rem",
  "48": "12rem",
  "52": "13rem",
  "56": "14rem",
  "60": "15rem",
  "64": "16rem",
  "72": "18rem",
  "80": "20rem",
  "96": "24rem",
};

const latte = {
  Rosewater: "#dc8a78", // hsl(11, 59%, 67%)
  Flamingo: "#dd7878", // hsl(0, 60%, 67%)
  Pink: "#ea76cb", // hsl(316, 73%, 69%)
  Mauve: "#8839ef", // hsl(266, 85%, 58%)
  Mauve2: "hsl(266, 85%, 53%)",
  Red: "#d20f39", // hsl(347, 87%, 44%)
  Red2: "hsl(347, 87%, 39%)",
  Maroon: "#e64553", // hsl(355, 76%, 59%)
  Peach: "#fe640b", // hsl(22, 99%, 52%)
  Yellow: "#df8e1d", // hsl(35, 77%, 49%)
  Yellow2: "hsl(35, 77%, 44%)",
  Green: "#40a02b", // hsl(109, 58%, 40%)
  Green2: "hsl(109, 58%, 35%)",
  Teal: "#179299", // hsl(183, 74%, 35%)
  Sky: "#04a5e5", // hsl(197, 97%, 46%)
  Sapphire: "#209fb5", // hsl(189, 70%, 42%)
  Blue: "#1e66f5", // hsl(220, 91%, 54%)
  Blue2: "hsl(220, 91%, 49%)",
  Lavender: "#7287fd", // hsl(231, 97%, 72%)
  Text: "#4c4f69", // hsl(234, 16%, 35%)
  Subtext1: "#5c5f77", // hsl(233, 13%, 41%)
  Subtext0: "#6c6f85", // hsl(233, 10%, 47%)
  Overlay2: "#7c7f93", // hsl(232, 10%, 53%)
  Overlay1: "#8c8fa1", // hsl(231, 10%, 59%)
  Overlay0: "#9ca0b0", // hsl(228, 11%, 65%)
  Surface2: "#acb0be", // hsl(227, 12%, 71%)
  Surface1: "#bcc0cc", // hsl(225, 14%, 77%)
  Surface0: "#ccd0da", // hsl(223, 16%, 83%)
  Base: "#eff1f5", // hsl(220, 23%, 95%)
  Mantle: "#e6e9ef", // hsl(220, 22%, 92%)
  Crust: "#dce0e8", // hsl(220, 21%, 89%)
};

const mocha = {
  Rosewater: "#f5e0dc", // hsl(10, 56%, 91%)
  Flamingo: "#f2cdcd", // hsl(0, 59%, 88%)
  Pink: "#f5c2e7", // hsl(316, 72%, 86%)
  Mauve: "#cba6f7", // hsl(267, 84%, 81%)
  Mauve2: "hsl(267, 84%, 86%)",
  Red: "#f38ba8", // hsl(343, 81%, 75%)
  Red2: "hsl(343, 81%, 80%)",
  Maroon: "#eba0ac", // hsl(350, 65%, 77%)
  Peach: "#fab387", // hsl(23, 92%, 75%)
  Yellow: "#f9e2af", // hsl(41, 86%, 83%)
  Yellow2: "hsl(41, 86%, 88%)",
  Green: "#a6e3a1", // hsl(115, 54%, 76%)
  Green2: "hsl(115, 54%, 81%)",
  Teal: "#94e2d5", // hsl(170, 57%, 73%)
  Sky: "#89dceb", // hsl(189, 71%, 73%)
  Sapphire: "#74c7ec", // hsl(199, 76%, 69%)
  Blue: "#89b4fa", // hsl(217, 92%, 76%)
  Blue2: "hsl(217, 92%, 81%)",
  Lavender: "#b4befe", // hsl(232, 97%, 85%)
  Text: "#cdd6f4", // hsl(226, 64%, 88%)
  Subtext1: "#bac2de", // hsl(227, 35%, 80%)
  Subtext0: "#a6adc8", // hsl(228, 24%, 72%)
  Overlay2: "#9399b2", // hsl(228, 17%, 64%)
  Overlay1: "#7f849c", // hsl(230, 13%, 55%)
  Overlay0: "#6c7086", // hsl(231, 11%, 47%)
  Surface2: "#585b70", // hsl(233, 12%, 39%)
  Surface1: "#45475a", // hsl(234, 13%, 31%)
  Surface0: "#313244", // hsl(237, 16%, 23%)
  Base: "#1e1e2e", // hsl(240, 21%, 15%)
  Mantle: "#181825", // hsl(240, 21%, 12%)
  Crust: "#11111b", // hsl(240, 23%, 9%)
};

const size = {
  sm: "640px",
  md: "768px",
  lg: "1024px",
  xl: "1280px",
  "2xl": "1536px",
};

export const minScreen = {
  sm: "screen and (min-width: 640px)",
  md: "screen and (min-width: 768px)",
  lg: "screen and (min-width: 1024px)",
  xl: "screen and (min-width: 1280px)",
  "2xl": "screen and (min-width: 1536px)",
};

export const maxScreen = {
  sm: "screen and (max-width: 639px)",
  md: "screen and (max-width: 767px)",
  lg: "screen and (max-width: 1023px)",
  xl: "screen and (max-width: 1279px)",
  "2xl": "screen and (max-width: 1535px)",
};

const themeDefault = {
  space,
  size,
  borderRadius: "4px",
  opacity: {
    disabled: "25%",
  },
};

export const [darkClass, theme] = createTheme({
  ...themeDefault,
  color: {
    ...mocha,
  },
});

export const lightClass = createTheme(theme, {
  ...themeDefault,
  color: {
    ...latte,
  },
});

