import "./index.css";

import { Navigate, Route, Router, } from '@solidjs/router'
import { Show, lazy } from "solid-js";

import { provideTheme } from "~/ui/theme";
import { NotFound } from '~/pages/404'
import { Home } from "~/pages/Home";
import loadHome from "~/pages/Home.data";
import { SignIn, Signup, Forgot } from "~/pages/Public";
import { Profile } from "~/pages/Profile";
import loadProfile from "~/pages/Profile.data";
import { Root } from "~/Root";
import loadRoot from "~/Root.data";
import { ClientProvider } from "~/providers/client";
import { lastSession } from "~/providers/session";
import { AdminGroups } from "~/pages/admin/Groups";
import loadAdminGroups from "~/pages/admin/Groups.data";
import { AdminGroupsID } from "~/pages/admin/GroupsID";
import loadAdminGroupsID from "~/pages/admin/GroupsID.data";
import { AdminUsers } from "~/pages/admin/Users";
import loadAdminUsers from "~/pages/admin/Users.data";
import { AdminDevices } from "~/pages/admin/Devices";
import loadAdminDevices from "~/pages/admin/Devices.data";
import loadAdminDevicesID from "~/pages/admin/DevicesID.data";
import { AdminDevicesID } from "~/pages/admin/DevicesID";
import { Live } from "~/pages/Live";
import { Devices } from "~/pages/Devices";
import loadDevices from "~/pages/Devices.data";
import { Emails } from "~/pages/Emails";
import loadEmails from "~/pages/Emails.data";
import { EmailsID } from "~/pages/EmailsID";
import loadEmailsID from "~/pages/EmailsID.data";
import { AdminSettings } from "~/pages/admin/Settings";
import loadAdminSettings from "~/pages/admin/Settings.data";
import { Files } from "~/pages/Files";
import loadFiles from "~/pages/Files.data";
import { Events } from "~/pages/Events";
import loadEvents from "~/pages/Events.data";

const Debug = lazy(() => import("./pages/debug"));

function NavigateHome() {
  return <Navigate href="/" />
}

function App() {
  provideTheme()
  const isAuthenticated = () => lastSession.valid && !lastSession.disabled
  const isAdmin = () => lastSession.admin

  return (
    <ClientProvider>
      <Router root={Root} rootLoad={loadRoot} explicitLinks>
        <Show when={import.meta.env.DEV}>
          <Route path="/debug">
            <Debug />
          </Route>
        </Show>
        <Show when={isAuthenticated()} fallback={<>
          <Route path="/signin" component={SignIn} />
          <Route path="/signup" component={Signup} />
          <Route path="/forgot" component={Forgot} />
          <Route path="*404" component={SignIn} />
        </>}>
          <Route path="/" component={Home} load={loadHome} />
          <Route path="/profile" component={Profile} load={loadProfile} />
          <Route path="/live" component={Live} />
          <Route path="/devices" component={Devices} load={loadDevices} />
          <Route path="/emails" component={Emails} load={loadEmails} />
          <Route path="/emails/:id" component={EmailsID} load={loadEmailsID} />
          <Route path="/events" component={Events} load={loadEvents} />
          <Route path="/files" component={Files} load={loadFiles} />
          <Show when={isAdmin()} fallback={<Route path="/admin/*" component={NavigateHome} />}>
            <Route path="/admin" component={AdminSettings} load={loadAdminSettings} />
            <Route path="/admin/users" component={AdminUsers} load={loadAdminUsers} />
            <Route path="/admin/groups" component={AdminGroups} load={loadAdminGroups} />
            <Route path="/admin/groups/:id" component={AdminGroupsID} load={loadAdminGroupsID} />
            <Route path="/admin/devices" component={AdminDevices} load={loadAdminDevices} />
            <Route path="/admin/devices/:id" component={AdminDevicesID} load={loadAdminDevicesID} />
          </Show>
          <Route path={["/signin", "/signup", "/forgot"]} component={NavigateHome} />
          <Route path="*404" component={NotFound} />
        </Show>
      </Router>
    </ClientProvider>
  )
}

export default App
