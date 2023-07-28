import { styled } from "@macaron-css/solid";
import { ParentComponent } from "solid-js";
import { theme } from "./theme";
import { utility } from "./utility";

const Center = styled("div", {
  base: {
    display: "flex",
    justifyContent: "center",
    padding: `${theme.space[16]} ${theme.space[4]} 0 ${theme.space[4]}`,
  },
});

const CenterChild = styled("div", {
  base: {
    ...utility.stack("4"),
    flex: "1",
    maxWidth: theme.space[96],
  },
});

export const LayoutCenter: ParentComponent = (props) => (
  <Center>
    <CenterChild>{props.children}</CenterChild>
  </Center>
);

export const LayoutDefault = styled("div", {
  base: {
    padding: theme.space[4],
  },
});
