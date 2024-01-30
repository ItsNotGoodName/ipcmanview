import "./index.css";

import { Navigate, Route, Router, } from '@solidjs/router'
import { Show, lazy } from "solid-js";

import { provideTheme } from "./ui/theme";
import { NotFound } from './pages/404'
import { Home } from "./pages/Home";
import loadHome from "./pages/Home.data";
import { View } from "./pages/View";
import { SignIn, Signup, Forgot } from "./pages/Public";
import { Profile } from "./pages/Profile";
import loadProfile from "./pages/Profile.data";
import { Root } from "./Root";
import { ClientProvider } from "./providers/client";
import { sessionCache } from "./providers/session";
import { AdminGroups } from "./pages/admin/Groups";
import loadAdminGroups from "./pages/admin/Groups.data";
import { AdminHome } from "./pages/admin/Home";
import { AdminGroupsID } from "./pages/admin/GroupsID";
import loadAdminGroupsID from "./pages/admin/GroupsID.data";
import { AdminUsers } from "./pages/admin/Users";
import loadAdminUsers from "./pages/admin/Users.data";
import { AdminDevices } from "./pages/admin/Devices";
import loadAdminDevices from "./pages/admin/Devices.data";
import loadAdminDevicesID from "./pages/admin/DevicesID.data";
import { AdminDevicesID } from "./pages/admin/DevicesID";

const Debug = lazy(() => import("./pages/debug"));

function App() {
  provideTheme()

  return (
    <ClientProvider>
      <Router root={Root}>
        <Show when={import.meta.env.DEV}>
          <Route path="/debug">
            <Debug />
          </Route>
        </Show>
        <Show when={sessionCache.valid && !sessionCache.disabled} fallback={<>
          <Route path="/signin" component={SignIn} />
          <Route path="/signup" component={Signup} />
          <Route path="/forgot" component={Forgot} />
          <Route path="*404" component={() => <Navigate href="/signin" />} />
        </>}>
          <Route path="/" component={Home} load={loadHome} />
          <Route path="/profile" component={Profile} load={loadProfile} />
          <Route path="/view" component={View} />
          <Show when={sessionCache.admin} fallback={<Route path="/admin/*" component={() => <>You are not an admin.</>}></Route>}>
            <Route path="/admin" component={AdminHome} />
            <Route path="/admin/users" component={AdminUsers} load={loadAdminUsers} />
            <Route path="/admin/groups" component={AdminGroups} load={loadAdminGroups} />
            <Route path="/admin/groups/:id" component={AdminGroupsID} load={loadAdminGroupsID} />
            <Route path="/admin/devices" component={AdminDevices} load={loadAdminDevices} />
            <Route path="/admin/devices/:id" component={AdminDevicesID} load={loadAdminDevicesID} />
          </Show>
          <Route path={["/signin", "/signup", "/forgot"]} component={() => <Navigate href="/" />} />
          <Route path="*404" component={NotFound} />
        </Show>
      </Router>
    </ClientProvider>
  )
}

export default App
