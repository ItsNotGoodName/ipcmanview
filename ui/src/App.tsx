import { styled } from "@macaron-css/solid";
import { theme } from '~/ui/theme';
import { themeModeClass } from '~/ui/theme-mode';
import { globalStyle } from '@macaron-css/core';
import { AuthProvider } from "./providers/auth";
import { Application } from "./views/Application";
import { Register } from "./views/Register";
import { Route, Routes } from "@solidjs/router";
import { Login } from "./views/Login";

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
      <AuthProvider
        anonymous={(
          <Routes >
            <Route path="/register" component={Register} />
            <Route path="/*" component={Login} />
          </Routes>
        )}
        authenticated={<Application />}
      />
    </Root>
  )
}

export default App
