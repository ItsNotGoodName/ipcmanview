import { cva } from "class-variance-authority"
import { As, DropdownMenu } from "@kobalte/core";
import { ErrorBoundary, JSX, ParentProps, Show, Suspense, createEffect, createSignal, splitProps } from "solid-js";
import { A, action, createAsync, revalidate, useAction, useLocation, Location, useNavigate, useSubmission } from "@solidjs/router";
import { RiDocumentFileLine, RiBuildingsHomeLine, RiDevelopmentBugLine, RiSystemLogoutBoxRFill, RiSystemMenuLine, RiUserFacesAdminLine, RiUserFacesUserLine, RiWeatherFlashlightLine, RiMediaLiveLine, RiBusinessMailLine, RiUserFacesGroupLine, RiSystemSettings2Line } from "solid-icons/ri";
import { Portal } from "solid-js/web";
import { makePersisted } from "@solid-primitives/storage";

import { DropdownMenuArrow, DropdownMenuContent, DropdownMenuPortal, DropdownMenuRoot, DropdownMenuTrigger } from "~/ui/DropdownMenu";
import { ThemeIcon } from "~/ui/ThemeIcon";
import { toggleTheme, useThemeTitle } from "~/ui/theme";
import { ToastList, ToastRegion } from "~/ui/Toast";
import { cn, catchAsToast } from "~/lib/utils";
import { getSession } from "~/providers/session";
import { PageError, PageLoading } from "./ui/Page";
import { BiRegularCctv } from "solid-icons/bi";

const menuLinkVariants = cva("ui-disabled:pointer-events-none ui-disabled:opacity-50 relative flex cursor-pointer select-none items-center gap-1 rounded-sm px-2 py-1.5 text-sm outline-none transition-colors", {
  variants: {
    size: {
      default: "h-10 px-4 py-2",
      icon: "h-10 w-10",
    },
    variant: {
      default: "hover:bg-accent hover:text-accent-foreground",
      active: "bg-primary text-primary-foreground hover:bg-primary/90",
    }
  },
  defaultVariants: {
    variant: "default"
  }
})

function useIsAdminPage<T>(location: Location<T>) {
  return () => location.pathname.startsWith("/admin")
}

function DropdownMenuLinks() {
  const navigate = useNavigate()

  return (
    <>
      <DropdownMenu.Item asChild onSelect={() => navigate("/")}>
        <As component={A} class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
          href="/" end>
          <RiBuildingsHomeLine class="h-5 w-5" />Home
        </As>
      </DropdownMenu.Item>
      <DropdownMenu.Item asChild onSelect={() => navigate("/devices")}>
        <As component={A} class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
          href="/devices">
          <BiRegularCctv class="h-5 w-5" />Devices
        </As>
      </DropdownMenu.Item>
      <DropdownMenu.Item asChild onSelect={() => navigate("/emails")}>
        <As component={A} class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
          href="/emails">
          <RiDocumentFileLine class="h-5 w-5" />Emails
        </As>
      </DropdownMenu.Item>
      <DropdownMenu.Item asChild onSelect={() => navigate("/events")}>
        <As component={A} class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
          href="/events">
          <RiWeatherFlashlightLine class="h-5 w-5" />Events
        </As>
      </DropdownMenu.Item>
      <DropdownMenu.Item asChild onSelect={() => navigate("/files")}>
        <As component={A} class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
          href="/files">
          <RiDocumentFileLine class="h-5 w-5" />Files
        </As>
      </DropdownMenu.Item>
      <DropdownMenu.Item asChild onSelect={() => navigate("/live")}>
        <As component={A} class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
          href="/live">
          <RiMediaLiveLine class="h-5 w-5" />Live
        </As>
      </DropdownMenu.Item>
    </>
  )
}

function MenuLinks() {
  return (
    <div class="flex flex-col p-2">
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
        href="/" noScroll end>
        <RiBuildingsHomeLine class="h-5 w-5" />Home
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
        href="/devices" noScroll>
        <BiRegularCctv class="h-5 w-5" />Devices
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
        href="/emails" noScroll>
        <RiBusinessMailLine class="h-5 w-5" />Emails
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
        href="/events" noScroll>
        <RiWeatherFlashlightLine class="h-5 w-5" />Events
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
        href="/files" noScroll>
        <RiDocumentFileLine class="h-5 w-5" />Files
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
        href="/live" noScroll>
        <RiMediaLiveLine class="h-5 w-5" />Live
      </A>
    </div>
  )
}

function AdminDropdownMenuLinks() {
  const navigate = useNavigate()

  return (
    <>
      <DropdownMenu.Item asChild onSelect={() => navigate("/admin")}>
        <As component={A} class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
          href="/admin" end>
          <RiSystemSettings2Line class="h-5 w-5" />Settings
        </As>
      </DropdownMenu.Item>
      <DropdownMenu.Item asChild onSelect={() => navigate("/admin/users")}>
        <As component={A} class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
          href="/admin/users" end>
          <RiUserFacesUserLine class="h-5 w-5" />Users
        </As>
      </DropdownMenu.Item>
      <DropdownMenu.Item asChild onSelect={() => navigate("/admin/groups")}>
        <As component={A} class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
          href="/admin/groups" end>
          <RiUserFacesGroupLine class="h-5 w-5" />Groups
        </As>
      </DropdownMenu.Item>
      <DropdownMenu.Item asChild onSelect={() => navigate("/admin/devices")}>
        <As component={A} class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
          href="/admin/devices" end>
          <BiRegularCctv class="h-5 w-5" />Devices
        </As>
      </DropdownMenu.Item>
    </>
  )
}

function AdminMenuLinks() {
  return (
    <div class="flex flex-col p-2">
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
        href="/admin" noScroll end>
        <RiSystemSettings2Line class="h-5 w-5" />Settings
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
        href="/admin/users" noScroll>
        <RiUserFacesUserLine class="h-5 w-5" />Users
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
        href="/admin/groups" noScroll>
        <RiUserFacesGroupLine class="h-5 w-5" />Groups
      </A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
        href="/admin/devices" noScroll>
        <BiRegularCctv class="h-5 w-5" />Devices
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
  }).catch(catchAsToast)
)

function Header(props: ParentProps<{ onMenuClick: () => void }>) {
  const signOutSubmission = useSubmission(actionSignOut)
  const signOutAction = useAction(actionSignOut)
  const signOut = () => signOutAction().catch(catchAsToast)
  const session = createAsync(() => getSession())
  const location = useLocation()
  const navigate = useNavigate()
  const isAdminPage = useIsAdminPage(location)

  return (
    <div class="bg-background text-foreground border-b-border z-10 h-12 w-full overflow-x-hidden border-b">
      <div class="flex h-full items-center gap-1 px-1">
        <DropdownMenuRoot>
          <DropdownMenuTrigger title="Menu" class={cn(menuLinkVariants(), "md:hidden")}>
            <RiSystemMenuLine class="h-6 w-6" />
          </DropdownMenuTrigger>
          <DropdownMenuPortal>
            <DropdownMenuContent class="md:hidden">
              <DropdownMenuArrow />
              {props.children}
            </DropdownMenuContent>
          </DropdownMenuPortal>
        </DropdownMenuRoot>
        <button onClick={props.onMenuClick} title="Menu" class={cn(menuLinkVariants(), "hidden md:inline-flex")}>
          <RiSystemMenuLine class="h-6 w-6" />
        </button>
        <A href="/" class="flex flex-1 items-center truncate text-xl">
          IPCManView
        </A>
        <div class="flex gap-1">
          <Show when={import.meta.env.DEV}>
            <A class={menuLinkVariants({ size: "icon" })} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants({ size: "icon" })}
              href="/debug" title="Debug" end>
              <RiDevelopmentBugLine class="h-6 w-6" />
            </A>
          </Show>
          <Show when={session()?.admin}>
            <A class={menuLinkVariants({ size: "icon" })} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants({ size: "icon" })}
              href={isAdminPage() ? "/" : "/admin"} title="Toggle admin">
              <RiUserFacesAdminLine class="h-6 w-6" />
            </A>
          </Show>
          <DropdownMenuRoot>
            <DropdownMenuTrigger class={menuLinkVariants({ size: "icon", variant: location.pathname.startsWith("/profile") ? "active" : "default" })} title="User">
              <RiUserFacesUserLine class="h-6 w-6" />
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
                    <RiUserFacesUserLine class="h-5 w-5" />Profile
                  </As>
                </DropdownMenu.Item>
                <DropdownMenu.Item class={menuLinkVariants()} onSelect={signOut} disabled={signOutSubmission.pending}>
                  <RiSystemLogoutBoxRFill class="h-5 w-5" />Sign out
                </DropdownMenu.Item>
              </DropdownMenuContent>
            </DropdownMenuPortal>
          </DropdownMenuRoot>
          <button class={menuLinkVariants({ size: "icon" })} onClick={toggleTheme} title={useThemeTitle()}>
            <ThemeIcon class="h-6 w-6" />
          </button>
        </div>
      </div>
    </div>
  )
}

function Menu(props: Omit<JSX.HTMLAttributes<HTMLDivElement>, "class"> & { menuOpen?: boolean }) {
  const [_, rest] = splitProps(props, ["children"])

  let refs: HTMLDivElement[] = []
  createEffect(() => {
    if (props.menuOpen) {
      refs.forEach(r => r.dataset.open = "")
    } else {
      refs.forEach(r => delete r.dataset.open)
    }
  })

  return (
    <div ref={refs[0]} class="border-border border-r-0 transition-all duration-200 md:data-[open=]:border-r" {...rest}>
      <div ref={refs[1]} class="sticky top-0 w-0 transition-all duration-200 md:data-[open=]:w-48">
        <div class="h-screen overflow-y-auto">
          {props.children}
        </div>
      </div>
    </div>
  )
}

export function Root(props: any) {
  const session = createAsync(() => getSession())
  const [menuOpen, setMenuOpen] = makePersisted(createSignal(true), { "name": "menu-open" })
  const layoutActive = () => session()?.valid && !session()?.disabled
  const isAdminPage = useIsAdminPage(props.location)

  return (
    <ErrorBoundary fallback={(e) =>
      <div class="p-4">
        <PageError error={e} />
      </div>
    }>
      <Suspense fallback={<PageLoading class="pt-10" />}>
        <Portal>
          <ToastRegion class={layoutActive() ? "top-12 sm:top-12" : ""}>
            <ToastList class={layoutActive() ? "top-12 sm:top-12" : ""} />
          </ToastRegion>
        </Portal>
        <Show when={layoutActive()} fallback={<>{props.children}</>}>
          <Header onMenuClick={() => setMenuOpen((prev) => !prev)}>
            <Show when={!isAdminPage()} fallback={<AdminDropdownMenuLinks />}>
              <DropdownMenuLinks />
            </Show>
          </Header>
          <div class="flex">
            <Menu menuOpen={menuOpen()}>
              <Show when={!isAdminPage()} fallback={<AdminMenuLinks />}>
                <MenuLinks />
              </Show>
            </Menu>
            <div class="w-full overflow-x-auto"> {/* FIXME: overflow-x-auto is needed to fix overflowing tables BUT it also breaks something and I forgot what it was ¯\_(ツ)_/¯ */}
              {props.children}
            </div>
          </div >
        </Show>
      </Suspense >
    </ErrorBoundary>
  )
}
