import "./index.css";

import { Navigate, Route, Router, } from '@solidjs/router'
import { Show } from "solid-js";
import { Portal } from "solid-js/web";

import { useTheme } from "./ui/theme";
import { AuthLayout } from "./layouts/Auth";
import { BaseLayout } from "./layouts/Base";
import { Debug } from './pages/debug'
import { NotFound } from './pages/404'
import { Home } from "./pages/Home";
import { View } from "./pages/View";
import { Signin, Signup, Forgot } from "./pages/Base";
import { useAuth } from "./providers/auth";
import { ToastList, ToastRegion } from "./ui/Toast";
import { Profile } from "./pages/Profile";
import { getProfile } from "./pages/Profile.data";

const base = "/next"

function App() {
  useTheme()
  const auth = useAuth()
  const toastClass = () => auth.valid() ? "top-12 sm:top-12" : ""

  return (
    <>
      <Portal>
        <ToastRegion class={toastClass()}>
          <ToastList class={toastClass()} />
        </ToastRegion>
      </Portal>
      <Show when={auth.valid()} fallback={
        <Router base={base} root={BaseLayout}>
          <Route path="/signin" component={Signin} />
          <Route path="/signup" component={Signup} />
          <Route path="/forgot" component={Forgot} />
          <Route path="*404" component={() => <Navigate href="/signin" />} />
        </Router>
      }>
        <Router base={base} root={AuthLayout}>
          <Debug />
          <Route path="/" component={Home} />
          <Route path="/profile" component={Profile} load={getProfile} />
          <Route path="/view" component={View} />
          <Route path={["/signin", "/signup", "/forgot"]} component={() => <Navigate href="/" />} />
          <Route path="*404" component={NotFound} />
        </Router>
      </Show>
    </>
  )
}

export default App
