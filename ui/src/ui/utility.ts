import { CSSProperties, keyframes } from "@macaron-css/core";
import { theme } from "./theme";

const rotate = keyframes({
  from: { transform: "rotate(0deg)" },
  to: { transform: "rotate(360deg)" },
});

export const utility = {
  animateSpin: {
    animation: `${rotate} 1s linear infinite`,
  } as CSSProperties,

  shadow: {
    boxShadow: "0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1)",
  } as CSSProperties,

  shadowXl: {
    boxShadow: `0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1)`,
  } as CSSProperties,

  textLine(): CSSProperties {
    return {
      overflow: "hidden",
      textOverflow: "ellipsis",
      whiteSpace: "nowrap",
    };
  },

  stack(space: keyof typeof theme["space"]): CSSProperties {
    return {
      display: "flex",
      flexDirection: "column",
      gap: theme.space[space],
    };
  },

  row(space: keyof typeof theme["space"]): CSSProperties {
    return {
      display: "flex",
      gap: theme.space[space],
    };
  },

  size(space: keyof typeof theme["space"]): CSSProperties {
    return {
      width: theme.space[space],
      height: theme.space[space],
    };
  },
};

