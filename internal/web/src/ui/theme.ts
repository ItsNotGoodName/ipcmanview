import { createEffect, createSignal } from "solid-js";

const THEME_KEY = "theme";
const DARK_THEME_CLASS = "dark"
const LIGHT_THEME_CLASS = "light"

export enum Theme {
  System = "system",
  Light = "light",
  Dark = "dark"
}

const query = window.matchMedia("(prefers-color-scheme: dark)");

const [currentSystemTheme, setCurrentSystemTheme] = createSignal(
  query.matches ? Theme.Dark : Theme.Light
);

query.addEventListener("change", (e: MediaQueryListEvent) => {
  setCurrentSystemTheme(e.matches ? Theme.Dark : Theme.Light);
});

const currentThemeSignal = createSignal(localStorage.getItem(THEME_KEY) ?? Theme.System);
export const useCurrentTheme = currentThemeSignal[0]
const setCurrentTheme = currentThemeSignal[1]

export function setTheme(theme: Theme) {
  setCurrentTheme(theme);
  theme === Theme.System
    ? localStorage.removeItem(THEME_KEY)
    : localStorage.setItem(THEME_KEY, theme);
}

export const toggleTheme = () => {
  switch (useCurrentTheme()) {
    case Theme.Light:
      setTheme(Theme.Dark);
      break
    case Theme.System:
      setTheme(Theme.Light);
      break
    default:
      setTheme(Theme.System);
  }
};

const themeClass = () => {
  if (useCurrentTheme() == Theme.System) {
    return currentSystemTheme() == Theme.Dark ? DARK_THEME_CLASS : LIGHT_THEME_CLASS;
  }
  return useCurrentTheme() == Theme.Dark ? DARK_THEME_CLASS : LIGHT_THEME_CLASS;
};

export const provideTheme = () => {
  createEffect(() => {
    document.getElementsByTagName("body")![0].className = themeClass()
  })
}

export const useThemeTitle = () => {
  switch (useCurrentTheme()) {
    case Theme.System:
      return "System theme"
    case Theme.Light:
      return "Light theme"
    case Theme.Dark:
      return "Dark theme"
  }
}
