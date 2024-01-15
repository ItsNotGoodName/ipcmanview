import {
  RiDesignContrastLine,
  RiWeatherMoonLine,
  RiWeatherSunLine,
} from "solid-icons/ri";
import { Component, Match, Switch } from "solid-js";
import { IconProps } from "solid-icons";
import {
  Theme,
  currentTheme,
} from "./theme";

export const ThemeIcon: Component<IconProps> = (props) => {
  return (
    <Switch fallback={<RiDesignContrastLine {...props} />}>
      <Match when={currentTheme() == Theme.Dark}>
        <RiWeatherMoonLine {...props} />
      </Match>
      <Match when={currentTheme() == Theme.Light}>
        <RiWeatherSunLine {...props} />
      </Match>
    </Switch>
  );
};
