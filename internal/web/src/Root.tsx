import { VariantProps, cva } from "class-variance-authority"
import { BiRegularCctv } from "solid-icons/bi";
import { As, DropdownMenu } from "@kobalte/core";
import { ErrorBoundary, Show, Suspense, batch, createSignal, } from "solid-js";
import { A, action, createAsync, revalidate, useAction, useLocation, useNavigate, useSubmission, RouteSectionProps } from "@solidjs/router";
import { RiDocumentFileLine, RiBuildingsHomeLine, RiDevelopmentBugLine, RiSystemLogoutBoxRFill, RiSystemMenuLine, RiUserFacesAdminLine, RiUserFacesUserLine, RiWeatherFlashlightLine, RiMediaLiveLine, RiBusinessMailLine, RiUserFacesGroupLine, RiSystemSettings2Line } from "solid-icons/ri";
import { Portal } from "solid-js/web";
import { makePersisted } from "@solid-primitives/storage";

import { DropdownMenuArrow, DropdownMenuContent, DropdownMenuPortal, DropdownMenuRoot, DropdownMenuTrigger } from "~/ui/DropdownMenu";
import { ThemeIcon } from "~/ui/ThemeIcon";
import { toggleTheme, useThemeTitle } from "~/ui/theme";
import { ToastList, ToastRegion } from "~/ui/Toast";
import { cn, catchAsToast } from "~/lib/utils";
import { getSession } from "~/providers/session";
import { PageError, PageLoading } from "~/ui/Page";
import { WSState, useWS } from "./providers/ws";
import { Shared } from "./components/Shared";
import { TooltipArrow, TooltipContent, TooltipRoot, TooltipTrigger } from "./ui/Tooltip";
import { SheetContent, SheetDescription, SheetHeader, SheetOverflow, SheetRoot, SheetTitle } from "./ui/Sheet";
import { useBus } from "./providers/bus";
import { getConfig } from "./pages/data";

const menuLinkVariants = cva("ui-disabled:pointer-events-none ui-disabled:opacity-50 relative flex cursor-pointer select-none items-center gap-1 rounded-sm px-2 py-1.5 text-sm outline-none transition-colors", {
  variants: {
    size: {
      default: "h-10 px-4 py-2",
      icon: "h-10 w-10",
    },
    variant: {
      default: "hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground",
      active: "bg-primary text-primary-foreground hover:bg-primary/90 focus:bg-primary/90",
    }
  },
  defaultVariants: {
    variant: "default"
  }
})

function MenuLinks(props: { onClick?: () => void }) {
  return (
    <div class="flex flex-col">
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()} onClick={props.onClick}
        href="/" noScroll end>
        <RiBuildingsHomeLine class="size-5" />Home
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()} onClick={props.onClick}
        href="/devices" noScroll>
        <BiRegularCctv class="size-5" />Devices
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()} onClick={props.onClick}
        href="/emails" noScroll>
        <RiBusinessMailLine class="size-5" />Emails
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()} onClick={props.onClick}
        href="/events" noScroll>
        <RiWeatherFlashlightLine class="size-5" />Events
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()} onClick={props.onClick}
        href="/files" noScroll>
        <RiDocumentFileLine class="size-5" />Files<div class="flex flex-1 justify-end"><div>🚧</div></div>
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()} onClick={props.onClick}
        href="/live" noScroll>
        <RiMediaLiveLine class="size-5" />Live<div class="flex flex-1 justify-end"><div>🚧</div></div>
      </A>
    </div>
  )
}

function AdminMenuLinks(props: { onClick?: () => void }) {
  return (
    <div class="flex flex-col">
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()} onClick={props.onClick}
        href="/admin" noScroll end>
        <RiBuildingsHomeLine class="size-5" />Home
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()} onClick={props.onClick}
        href="/admin/settings" noScroll>
        <RiSystemSettings2Line class="size-5" />Settings
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()} onClick={props.onClick}
        href="/admin/users" noScroll>
        <RiUserFacesUserLine class="size-5" />Users
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()} onClick={props.onClick}
        href="/admin/groups" noScroll>
        <RiUserFacesGroupLine class="size-5" />Groups
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()} onClick={props.onClick}
        href="/admin/devices" noScroll>
        <BiRegularCctv class="size-5" />Devices
      </A>
    </div>
  )
}

const actionSignOut = action(() =>
  fetch("/v1/session", {
    credentials: "include",
    headers: [['Content-Type', 'application/json'], ['Accept', 'application/json']],
    method: "DELETE"
  }).then(async (resp) => {
    if (!resp.ok) {
      const json = await resp.json()
      throw new Error(json.message)
    }
    return revalidate(getSession.key)
  }).catch(catchAsToast))

type HeaderProps = {
  onMenuClick: () => void
  onMobileMenuClick: () => void
  isAdminPage: boolean
  siteName?: string
}

function Header(props: HeaderProps) {
  const location = useLocation()
  const navigate = useNavigate()

  const session = createAsync(() => getSession())

  const signOutSubmission = useSubmission(actionSignOut)
  const signOut = useAction(actionSignOut)

  const ws = useWS()
  const wsState = (): VariantProps<typeof Shared.connectionIndicatorVariants>["state"] => {
    switch (ws.state()) {
      case WSState.Connecting:
        return "connecting"
      case WSState.Connected:
        return "connected"
      case WSState.Disconnecting:
      case WSState.Disconnected:
        return "disconnected"
    }
  }

  return (
    <div class="bg-background text-foreground border-b-border z-10 h-12 w-full overflow-x-hidden border-b">
      <div class="flex h-full items-center gap-1 px-1">
        <div onClick={props.onMobileMenuClick} title="Menu" class={cn(menuLinkVariants(), "md:hidden")}>
          <RiSystemMenuLine class="size-6" />
        </div>
        <button onClick={props.onMenuClick} title="Menu" class={cn(menuLinkVariants(), "hidden md:inline-flex")}>
          <RiSystemMenuLine class="size-6" />
        </button>
        <div class="flex flex-1 items-baseline gap-2 truncate">
          <A href="/" class="flex items-center text-xl">
            IPCManView
          </A>
          <p class="text-muted-foreground text-sm">{props.siteName}</p>
        </div>
        <div class="flex gap-1">
          <TooltipRoot>
            <TooltipTrigger class="px-2">
              <div class={Shared.connectionIndicatorVariants({ state: wsState() })} />
            </TooltipTrigger>
            <TooltipContent>
              <TooltipArrow />
              <p>WebSocket {wsState()}</p>
            </TooltipContent>
          </TooltipRoot>
          <Show when={import.meta.env.DEV}>
            <A class={menuLinkVariants({ size: "icon" })} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants({ size: "icon" })}
              href="/debug" title="Debug" end>
              <RiDevelopmentBugLine class="size-6" />
            </A>
          </Show>
          <Show when={session()?.admin}>
            <A class={menuLinkVariants({ size: "icon" })} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants({ size: "icon" })}
              href={props.isAdminPage ? "/" : "/admin"} title="Toggle admin">
              <RiUserFacesAdminLine class="size-6" />
            </A>
          </Show>
          <DropdownMenuRoot>
            <DropdownMenuTrigger class={menuLinkVariants({ size: "icon", variant: location.pathname.startsWith("/profile") ? "active" : "default" })} title="User">
              <RiUserFacesUserLine class="size-6" />
            </DropdownMenuTrigger>
            <DropdownMenuPortal>
              <DropdownMenuContent class="z-[200]">
                <DropdownMenuArrow />
                <DropdownMenu.Item class="truncate px-2 pb-1.5 text-lg font-semibold" closeOnSelect={false}>
                  {session()?.username}
                </DropdownMenu.Item>
                <DropdownMenu.Item asChild onSelect={() => navigate("/profile")}>
                  <As component={A} inactiveClass={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })}
                    href="/profile" end>
                    <RiUserFacesUserLine class="size-5" />Profile
                  </As>
                </DropdownMenu.Item>
                <DropdownMenu.Item class={menuLinkVariants()} onSelect={signOut} disabled={signOutSubmission.pending}>
                  <RiSystemLogoutBoxRFill class="size-5" />Sign out
                </DropdownMenu.Item>
              </DropdownMenuContent>
            </DropdownMenuPortal>
          </DropdownMenuRoot>
          <button class={menuLinkVariants({ size: "icon" })} onClick={toggleTheme} title={useThemeTitle()}>
            <ThemeIcon class="size-6" />
          </button>
        </div>
      </div>
    </div>
  )
}

function createMenu() {
  const [mobileOpen, setMobileOpen] = createSignal(false)
  const toggleMobileOpen = () => setMobileOpen(!mobileOpen())
  const closeMobile = () => setMobileOpen(false)

  const [open, setOpen] = makePersisted(createSignal(true), { "name": "menu-open" })
  const toggleOpen = () => {
    if (open()) {
      batch(() => {
        setOpen(false)
        setMobileOpen(false)
      })
    } else {
      setOpen(true)
    }
  }

  return {
    mobileOpen,
    toggleMobileOpen,
    closeMobile,
    open,
    toggleOpen,
  }
}

export function Root(props: RouteSectionProps) {
  const bus = useBus()

  const config = createAsync(() => getConfig())
  const session = createAsync(() => getSession())

  bus.event.listen((e) => {
    if (e.action == "user-security:updated" && e.data == session()?.user_id)
      revalidate(getSession.key)
  })

  const isAuthenticated = () => session()?.valid && !session()?.disabled
  const isAdminPage = () => props.location.pathname.startsWith("/admin")

  const menu = createMenu()

  return (
    <ErrorBoundary fallback={(e) =>
      <div class="p-4">
        <PageError error={e} />
      </div>
    }>
      <Suspense fallback={<PageLoading class="pt-10" />}>
        <Portal>
          <ToastRegion class={isAuthenticated() ? "top-12 sm:top-12" : ""}>
            <ToastList class={isAuthenticated() ? "top-12 sm:top-12" : ""} />
          </ToastRegion>
        </Portal>
        <Show when={isAuthenticated()} fallback={<>{props.children}</>}>
          <SheetRoot open={menu.mobileOpen()} onOpenChange={menu.toggleMobileOpen}>
            <SheetContent side="left" class="p-2">
              <SheetHeader class="px-2 sm:pt-2">
                <SheetTitle>IPCManView</SheetTitle>
                <SheetDescription>{config()?.siteName}</SheetDescription>
              </SheetHeader>
              <SheetOverflow class="pb-2">
                <Show when={!isAdminPage()} fallback={<AdminMenuLinks onClick={menu.closeMobile} />}>
                  <MenuLinks onClick={menu.closeMobile} />
                </Show>
              </SheetOverflow>
            </SheetContent>
          </SheetRoot>
          <Header
            onMenuClick={menu.toggleOpen}
            onMobileMenuClick={menu.toggleMobileOpen}
            isAdminPage={isAdminPage()}
            siteName={config()?.siteName}
          />
          <div class="flex">
            <div data-open={menu.open()} class="border-border w-0 shrink-0 border-r-0 transition-all duration-300 md:data-[open=true]:w-48 md:data-[open=true]:border-r">
              <div class="sticky top-0 max-h-screen overflow-y-auto overflow-x-clip">
                <div class="p-2">
                  <Show when={!isAdminPage()} fallback={
                    <AdminMenuLinks />
                  }>
                    <MenuLinks />
                  </Show>
                </div>
              </div>
            </div>
            <div class="min-h-[calc(100vh-3rem)] w-full overflow-x-clip">
              {props.children}
            </div>
          </div>
        </Show>
      </Suspense>
    </ErrorBoundary>
  )
}
