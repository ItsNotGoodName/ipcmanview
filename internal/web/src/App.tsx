import "./styles/index.css";
import "./styles/hljs.css";

import { Navigate, Route, Router, } from '@solidjs/router'
import { Show, lazy } from "solid-js";

import { provideTheme } from "~/ui/theme";
import { lastConfig } from "./pages/data";
import { lastSession } from "~/providers/session";

import { Root } from "~/Root";
import loadRoot from "~/Root.data";
import { NotFound } from '~/pages/404'
import { SignIn, SignUp, Forgot } from "~/pages/Public";

const Debug = lazy(() => import("./pages/debug"));
const Home = lazy(() => import("~/pages/Home"));
import loadHome from "~/pages/Home.data";
const Profile = lazy(() => import("~/pages/Profile"));
import loadProfile from "~/pages/Profile.data";
const Live = lazy(() => import("~/pages/Live"));
const Devices = lazy(() => import("~/pages/Devices"));
import loadDevices from "~/pages/Devices.data";
const Emails = lazy(() => import("~/pages/Emails"));
import loadEmails from "~/pages/Emails.data";
const EmailsID = lazy(() => import("~/pages/EmailsID"));
import loadEmailsID from "~/pages/EmailsID.data";
const Files = lazy(() => import("~/pages/Files"));
import loadFiles from "~/pages/Files.data";
const Events = lazy(() => import("~/pages/Events"));
import loadEvents from "~/pages/Events.data";
const EventsLive = lazy(() => import("~/pages/EventsLive"));
import loadEventsLive from "~/pages/EventsLive.data";
const AdminSettings = lazy(() => import("./pages/admin/Settings"));
import loadAdminSettings from "~/pages/admin/Settings.data";
const AdminHome = lazy(() => import("./pages/admin/Home"));
import loadAdminHome from "./pages/admin/Home.data";
const AdminGroups = lazy(() => import("~/pages/admin/Groups"));
import loadAdminGroups from "~/pages/admin/Groups.data";
const AdminGroupsID = lazy(() => import("~/pages/admin/GroupsID"));
import loadAdminGroupsID from "~/pages/admin/GroupsID.data";
const AdminUsers = lazy(() => import("~/pages/admin/Users"));
import loadAdminUsers from "~/pages/admin/Users.data";
const AdminDevices = lazy(() => import("~/pages/admin/Devices"));
import loadAdminDevices from "~/pages/admin/Devices.data";
const AdminDevicesID = lazy(() => import("~/pages/admin/DevicesID"));
import loadAdminDevicesID from "~/pages/admin/DevicesID.data";

function NavigateHome() {
  return <Navigate href="/" />
}

function App() {
  provideTheme()
  const isAuthenticated = () => lastSession.valid && !lastSession.disabled
  const isAdmin = () => lastSession.admin

  return (
    <Router root={Root} rootLoad={loadRoot} explicitLinks>
      <Show when={import.meta.env.DEV}>
        <Route path="/debug">
          <Debug />
        </Route>
      </Show>
      <Show when={isAuthenticated()} fallback={<>
        <Route path="/signin" component={SignIn} />
        <Show when={lastConfig.enableSignUp} fallback={
          <Route path="/signup" component={() => <Navigate href="/signin" />} />
        }>
          <Route path="/signup" component={SignUp} />
        </Show>
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
        <Route path="/events/live" component={EventsLive} load={loadEventsLive} />
        <Route path="/files" component={Files} load={loadFiles} />
        <Show when={isAdmin()} fallback={<Route path="/admin/*" component={NavigateHome} />}>
          <Route path="/admin" component={AdminHome} load={loadAdminHome} />
          <Route path="/admin/settings" component={AdminSettings} load={loadAdminSettings} />
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
  )
}

export default App
