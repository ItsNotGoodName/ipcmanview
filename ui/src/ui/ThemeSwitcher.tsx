import { styled } from "@macaron-css/solid";
import {
  RiDesignContrastLine,
  RiWeatherMoonLine,
  RiWeatherSunLine,
} from "solid-icons/ri";
import { Component, JSX, Match, ParentComponent, Switch } from "solid-js";
import { IconProps } from "solid-icons";
import {
  DARK_MODE,
  LIGHT_MODE,
  themeMode,
  toggleThemeMode,
} from "./theme-mode";
import { theme } from "./theme";

const TheButton = styled("button", {
  base: {
    padding: 0,
    border: "none",
    background: "none",
    color: theme.color.Text,
    cursor: "pointer",
  },
});

export const ThemeSwitcher: ParentComponent<
  Omit<
    JSX.ButtonHTMLAttributes<HTMLButtonElement>,
    "onClick" | "aria-label" | "title"
  >
> = (props) => (
  <TheButton
    {...props}
    onClick={toggleThemeMode}
    aria-label="Toggle Theme"
    title="Toggle Theme"
  />
);

export const ThemeSwitcherIcon: Component<IconProps> = (props) => {
  return (
    <Switch fallback={<RiDesignContrastLine {...props} />}>
      <Match when={themeMode() == DARK_MODE}>
        <RiWeatherMoonLine {...props} />
      </Match>
      <Match when={themeMode() == LIGHT_MODE}>
        <RiWeatherSunLine {...props} />
      </Match>
    </Switch>
  );
};

