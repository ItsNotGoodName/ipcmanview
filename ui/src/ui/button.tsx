import { styled } from "@macaron-css/solid";
import { theme } from "./theme";

export const Button = styled("button", {
  base: {
    appearance: "none",
    whiteSpace: "nowrap",
    border: "none",
    borderRadius: theme.borderRadius,
    cursor: "pointer",
    ":disabled": {
      cursor: "not-allowed",
      opacity: theme.opacity.disabled,
    },
  },
  variants: {
    size: {
      small: {
        padding: `${theme.space["0.5"]} ${theme.space[2]}`,
      },
      medium: {
        padding: `${theme.space[2]} ${theme.space[2]}`,
      },
      large: {
        padding: `${theme.space[4]} ${theme.space[4]}`,
      },
    },
    color: {
      primary: {
        background: theme.color.Mauve,
        selectors: {
          ["&:hover:enabled"]: {
            background: theme.color.Mauve2,
          },
        },
        color: theme.color.Crust,
      },
      secondary: {
        background: theme.color.Subtext0,
        selectors: {
          ["&:hover:enabled"]: {
            background: theme.color.Subtext1,
          },
        },
        color: theme.color.Crust,
      },
      success: {
        background: theme.color.Green,
        selectors: {
          ["&:hover:enabled"]: {
            background: theme.color.Green2,
          },
        },
        color: theme.color.Crust,
      },
      danger: {
        background: theme.color.Red,
        selectors: {
          ["&:hover:enabled"]: {
            background: theme.color.Red2,
          },
        },
        color: theme.color.Crust,
      },
    },
  },
  defaultVariants: {
    size: "medium",
    color: "primary",
  },
});

