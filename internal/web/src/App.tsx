import "./index.css";

import { Navigate, Route, Router, } from '@solidjs/router'

import { useTheme } from "./ui/theme";
import { Debug } from './pages/debug'
import { NotFound } from './pages/404'
import { Home } from "./pages/Home";
import { View } from "./pages/View";
import { SignIn, Signup, Forgot } from "./pages/Landing";
import { Profile } from "./pages/Profile";
import { loadProfile } from "./pages/Profile.data";
import { Layout } from "./Layout";
import { ClientProvider } from "./providers/client";
import { Match, Switch } from "solid-js";
import { session } from "./providers/session";

function App() {
  useTheme()

  return (
    <ClientProvider>
      <Router root={Layout}> {/* FIXME: solid-router explicitLinks={true} is broken, https://github.com/solidjs/solid-router/issues/356 */}
        <Debug />
        <Switch>
          <Match when={session()}>
            <Route path="/" component={Home} />
            <Route path="/profile" component={Profile} load={loadProfile} />
            <Route path="/view" component={View} />
            <Route path={["/signin", "/signup", "/forgot"]} component={() => <Navigate href="/" />} />
            <Route path="*404" component={NotFound} />
          </Match>
          <Match when={!session()}>
            <Route path="/signin" component={SignIn} />
            <Route path="/signup" component={Signup} />
            <Route path="/forgot" component={Forgot} />
            <Route path="*404" component={() => <Navigate href="/signin" />} />
          </Match>
        </Switch>
      </Router>
    </ClientProvider>
  )
}

export default App
