import { VariantProps, cva } from "class-variance-authority"
import { BiRegularCctv } from "solid-icons/bi";
import { As, DropdownMenu } from "@kobalte/core";
import { ErrorBoundary, JSX, Show, Suspense, batch, createSignal, splitProps } from "solid-js";
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
        <RiDocumentFileLine class="size-5" />Files
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()} onClick={props.onClick}
        href="/live" noScroll>
        <RiMediaLiveLine class="size-5" />Live
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
  const signOutAction = useAction(actionSignOut)
  const signOut = () => signOutAction().catch(catchAsToast)

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

function Menu(props: JSX.HTMLAttributes<HTMLDivElement> & { open?: boolean }) {
  const [_, rest] = splitProps(props, ["children", "open", "class"])
  return (
    <div data-open={props.open} class="border-border border-r-0 transition-all duration-200 md:data-[open=true]:border-r" {...rest}>
      <div data-open={props.open} class="sticky top-0 w-0 transition-all duration-200 md:data-[open=true]:w-48">
        <div class="h-screen overflow-y-auto overflow-x-hidden">
          <div class={props.class}>
            {props.children}
          </div>
        </div>
      </div>
    </div>
  )
}

export function Root(props: RouteSectionProps) {
  const bus = useBus()
  const config = createAsync(() => getConfig())
  const session = createAsync(() => getSession())

  bus.event.listen((e) => {
    if (e.action == "user-security:updated" && e.data == session()?.user_id)
      revalidate(getSession.key)
  })

  const isAuthenticatedLayout = () => session()?.valid && !session()?.disabled
  const isAdminPage = () => props.location.pathname.startsWith("/admin")

  const [mobileMenuOpen, setMobileMenuOpen] = createSignal(false)
  const toggleMobileMenuOpen = () => setMobileMenuOpen(!mobileMenuOpen())
  const closeMobileMenu = () => setMobileMenuOpen(false)

  const [menuOpen, setMenuOpen] = makePersisted(createSignal(true), { "name": "menu-open" })
  const toggleMenuOpen = () => {
    if (menuOpen()) {
      batch(() => {
        setMenuOpen(false)
        setMobileMenuOpen(false)
      })
    } else {
      setMenuOpen(true)
    }
  }

  return (
    <ErrorBoundary fallback={(e) =>
      <div class="p-4">
        <PageError error={e} />
      </div>
    }>
      <Suspense fallback={<PageLoading class="pt-10" />}>
        <Portal>
          <ToastRegion class={isAuthenticatedLayout() ? "top-12 sm:top-12" : ""}>
            <ToastList class={isAuthenticatedLayout() ? "top-12 sm:top-12" : ""} />
          </ToastRegion>
        </Portal>
        <Show when={isAuthenticatedLayout()} fallback={<>{props.children}</>}>
          <SheetRoot open={mobileMenuOpen()} onOpenChange={toggleMobileMenuOpen}>
            <SheetContent side="left">
              <SheetHeader>
                <SheetTitle>IPCManView</SheetTitle>
                <SheetDescription>{config()?.siteName}</SheetDescription>
              </SheetHeader>
              <SheetOverflow>
                <Show when={!isAdminPage()} fallback={<AdminMenuLinks onClick={closeMobileMenu} />}>
                  <MenuLinks onClick={closeMobileMenu} />
                </Show>
              </SheetOverflow>
            </SheetContent>
          </SheetRoot>
          <Header
            onMenuClick={toggleMenuOpen}
            onMobileMenuClick={toggleMobileMenuOpen}
            isAdminPage={isAdminPage()}
            siteName={config()?.siteName}
          />
          <div class="flex">
            <Menu open={menuOpen()} class="flex h-full flex-col justify-between p-2">
              <Show when={!isAdminPage()} fallback={<AdminMenuLinks />}>
                <MenuLinks />
              </Show>
              <div class="flex flex-col">
                <button class={menuLinkVariants()} onClick={toggleMenuOpen}>
                  <RiSystemMenuLine class="size-5" />Menu
                </button>
              </div>
            </Menu>
            <div class="w-full overflow-x-auto">
              {props.children}
            </div>
          </div>
        </Show>
      </Suspense>
    </ErrorBoundary>
  )
}
