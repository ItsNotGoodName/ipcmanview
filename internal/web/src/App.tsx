import "./index.css";

import { Navigate, Route, Router, } from '@solidjs/router'
import { Show, lazy } from "solid-js";

import { provideTheme } from "./ui/theme";
import { NotFound } from './pages/404'
import { Home } from "./pages/Home";
import { View } from "./pages/View";
import { SignIn, Signup, Forgot } from "./pages/Landing";
import { Profile } from "./pages/Profile";
import { loadProfile } from "./pages/Profile.data";
import { Layout } from "./Layout";
import { ClientProvider } from "./providers/client";
import { AdminHome } from "./pages/admin/Home";
import { sessionCache } from "./providers/session";

const Debug = lazy(() => import("./pages/debug"));

function App() {
  provideTheme()

  return (
    <ClientProvider>
      <Router root={Layout}> {/* FIXME: solid-router explicitLinks={true} is broken, https://github.com/solidjs/solid-router/issues/356 */}
        <Show when={import.meta.env.DEV}>
          <Route path="/debug">
            <Debug />
          </Route>
        </Show>
        <Show when={sessionCache.valid} fallback={
          <>
            <Route path="/signin" component={SignIn} />
            <Route path="/signup" component={Signup} />
            <Route path="/forgot" component={Forgot} />
            <Route path="*404" component={() => <Navigate href="/signin" />} />
          </>
        }>
          <Route path="/" component={Home} />
          <Route path="/profile" component={Profile} load={loadProfile} />
          <Route path="/view" component={View} />
          <Show when={sessionCache.admin} fallback={<Route path="/admin/*" component={() => <>You are not an admin.</>}></Route>}>
            <Route path="/admin" component={AdminHome} />
          </Show>
          <Route path={["/signin", "/signup", "/forgot"]} component={() => <Navigate href="/" />} />
          <Route path="*404" component={NotFound} />
        </Show>
      </Router>
    </ClientProvider >
  )
}

export default App
