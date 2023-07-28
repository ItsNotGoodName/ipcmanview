import { styled } from "@macaron-css/solid";
import { theme } from '~/ui/theme';
import { themeModeClass } from '~/ui/theme-mode';
import { globalStyle } from '@macaron-css/core';
import { Login } from "~/views/Login";
import { AuthService, UserService } from "~/core/client.gen";

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
  const auth = new AuthService(import.meta.env.VITE_BACKEND_URL, fetch)
  auth.register({
    user: {
      username: "fancy",
      email: "admin123@example.com",
      password: "12345678",
      passwordConfirm: "12345678",
    }
  })
  auth.login({
    usernameOrEmail: "fancy",
    password: "12345678"
  }).then((res) => {
    const user = new UserService(import.meta.env.VITE_BACKEND_URL, fetch)
    user.me({ "Authorization": `BEARER ${res.token}` })
  })

  return (
    <Root class={themeModeClass()}>
      <Login />
    </Root>
  )
}

export default App
