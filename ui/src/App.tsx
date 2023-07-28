import { styled } from "@macaron-css/solid";
import { theme } from '~/ui/theme';
import { themeModeClass } from '~/ui/theme-mode';
import { globalStyle } from '@macaron-css/core';
import { Login } from "~/views/Login";
import { AuthProvider } from "./providers/auth";
import { Loading } from "./views/Loading";
import { Application } from "./views/Application";

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
      <AuthProvider login={<Login />} loading={<Loading />}>
        <Application />
      </AuthProvider>
    </Root>
  )
}

export default App
