import { createSignal } from 'solid-js'
import { ExampleService } from "./core/client"
import { styled } from "@macaron-css/solid";
import { theme } from './ui/theme';
import { themeModeClass } from './ui/theme-mode';
import { Card, CardBody } from './ui/card';
import { ThemeSwitcher, ThemeSwitcherIcon } from './ui/theme-switcher';
import { utility } from './ui/utility';
import { style } from '@macaron-css/core';

const Root = styled("div", {
  base: {
    background: theme.color.Base,
    color: theme.color.Text,
    position: "fixed",
    inset: 0,
  },
});

const RootChild = styled("div", {
  base: {
    padding: theme.space[2],
  }
})

function App() {
  const [message, setMessage] = createSignal("Loading...");
  const api = new ExampleService(import.meta.env.VITE_BACKEND_URL, fetch)

  api.message().then(({ message }) => {
    setMessage(message.body + " " + message.time)
  }).catch((err) => {
    setMessage("Error: " + err)
  })

  return (
    <Root class={themeModeClass()}>
      <RootChild>
        <ThemeSwitcher>
          <ThemeSwitcherIcon class={style({ ...utility.size("8") })} />
        </ThemeSwitcher>
        <Card>
          <CardBody>
            {message()}
          </CardBody>
        </Card>
      </RootChild>
    </Root>
  )
}

export default App
