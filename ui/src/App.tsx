import { styled } from "@macaron-css/solid";
import { theme } from './ui/theme';
import { themeModeClass } from './ui/theme-mode';
import { Application } from './views/Application';
import { globalStyle } from '@macaron-css/core';

globalStyle("a", {
  textDecoration: "none",
  color: theme.color.Blue,
});

const Root = styled("div", {
  base: {
    background: theme.color.Base,
    color: theme.color.Text,
    position: "fixed",
    inset: 0,
  },
});


function App() {
  return (
    <Root class={themeModeClass()}>
      <Application />
    </Root>
  )
}

export default App
