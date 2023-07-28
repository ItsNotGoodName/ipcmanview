import { style } from "@macaron-css/core";
import { styled } from "@macaron-css/solid";
import { Component } from "solid-js";

import { LayoutCenter } from "~/ui/Layout";
import { theme } from "~/ui/theme";
import { IconSpinner } from "~/ui/Icon";

const Title = styled("div", {
  base: {
    fontWeight: "bold",
    fontSize: "x-large",
  },
});

const Center = styled("div", {
  base: {
    display: "flex",
    justifyContent: "center",
    textAlign: "center",
  },
});

const iconClass = style({
  height: theme.space[10],
  width: theme.space[10],
});

export const Loading: Component = () => (
  <LayoutCenter>
    <Center>
      <Title>IPCManView</Title>
    </Center>
    <Center>
      <IconSpinner class={iconClass} />
    </Center>
  </LayoutCenter>
);

