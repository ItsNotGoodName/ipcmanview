import { CSSProperties } from "@macaron-css/core";
import { styled } from "@macaron-css/solid";

import { theme } from "./theme";
import { utility } from "./utility";

export const Card = styled("div", {
  base: {
    ...utility.shadow,
    overflow: "hidden",
    background: theme.color.Surface0,
    borderRadius: theme.borderRadius,
  },
});

const border = {
  borderLeft: `${theme.space.px} solid ${theme.color.Overlay0}`,
  borderRight: `${theme.space.px} solid ${theme.color.Overlay0}`,
  ":first-child": {
    borderTopLeftRadius: theme.borderRadius,
    borderTopRightRadius: theme.borderRadius,
    borderTop: `${theme.space.px} solid ${theme.color.Overlay0}`,
  },
  ":last-child": {
    borderBottomLeftRadius: theme.borderRadius,
    borderBottomRightRadius: theme.borderRadius,
    borderBottom: `${theme.space.px} solid ${theme.color.Overlay0}`,
  },
} as CSSProperties;

export const CardHeader = styled("div", {
  base: {
    ...border,
    ...utility.row("2"),
    justifyContent: "space-between",
    alignItems: "center",
    background: theme.color.Surface1,
    height: theme.space[10],
    paddingLeft: theme.space[2],
    paddingRight: theme.space[2],
    selectors: {
      ["&:not(:last-child)"]: {
        borderBottom: `${theme.space.px} solid ${theme.color.Overlay0}`,
      },
    },
  },
});

export const CardHeaderTitle = styled("div", {
  base: {
    ...utility.textLine(),
  },
});

export const CardHeaderRow = styled("div", {
  base: {
    ...utility.row("2"),
  },
});

export const CardBody = styled("div", {
  base: {
    ...border,
    overflowX: "auto",
  },
  variants: {
    padding: {
      true: {
        padding: theme.space[4],
      },
    },
  },
  defaultVariants: {
    padding: true,
  },
});

