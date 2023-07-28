import { style } from "@macaron-css/core";
import { styled } from "@macaron-css/solid";
import { theme } from "./theme";

export const Menu = styled("div", {
  base: {
    padding: theme.space[1],
    display: "flex",
    flexDirection: "column",
    gap: theme.space[1],
  },
});

export const menuChildClass = style({
  cursor: "pointer",
  background: theme.color.Surface0,
  textAlign: "left",
  border: "none",
  borderRadius: theme.borderRadius,
  color: theme.color.Text,
  padding: theme.space[1],
  ":hover": {
    backgroundColor: theme.color.Surface2,
  },
});

