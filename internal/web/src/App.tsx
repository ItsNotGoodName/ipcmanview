import "./index.css";

import { Navigate, Route, Router, } from '@solidjs/router'
import { Match, Show, Switch, lazy } from "solid-js";

import { useTheme } from "./ui/theme";
import { NotFound } from './pages/404'
import { Home } from "./pages/Home";
import { View } from "./pages/View";
import { SignIn, Signup, Forgot } from "./pages/Landing";
import { Profile } from "./pages/Profile";
import { loadProfile } from "./pages/Profile.data";
import { Layout } from "./Layout";
import { ClientProvider } from "./providers/client";
import { session } from "./providers/session";
import { AdminHome } from "./pages/admin/Home";

const Debug = lazy(() => import("./pages/debug"));

function App() {
  useTheme()

  return (
    <ClientProvider>
      <Router root={Layout}> {/* FIXME: solid-router explicitLinks={true} is broken, https://github.com/solidjs/solid-router/issues/356 */}
        <Show when={import.meta.env.DEV}>
          <Route path="/debug">
            <Debug />
          </Route>
        </Show>
        <Switch>
          <Match when={session()}>
            <Route path="/" component={Home} />
            <Route path="/profile" component={Profile} load={loadProfile} />
            <Route path="/view" component={View} />
            <Route path="/admin" component={AdminHome} />
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
