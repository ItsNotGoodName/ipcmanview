import { cva } from "class-variance-authority"
import { As, DropdownMenu } from "@kobalte/core";
import { JSX, ParentProps, Show, Suspense, createEffect, createSignal, splitProps } from "solid-js";
import { A, action, createAsync, revalidate, useAction, useSubmission } from "@solidjs/router";
import { RiBuildingsHomeLine, RiDevelopmentBugLine, RiSystemEyeLine, RiSystemLoader4Line, RiSystemLogoutBoxRFill, RiSystemMenuLine, RiUserFacesUserLine } from "solid-icons/ri";
import { Portal } from "solid-js/web";

import { Button } from "~/ui/Button";
import { DropdownMenuArrow, DropdownMenuContent, DropdownMenuPortal, DropdownMenuRoot, DropdownMenuTrigger } from "~/ui/DropdownMenu";
import { ThemeIcon } from "~/ui/ThemeIcon";
import { toggleTheme, useThemeTitle } from "~/ui/theme";
import { makePersisted } from "@solid-primitives/storage";
import { ToastList, ToastRegion } from "~/ui/Toast";
import { cn, toastError } from "~/lib/utils";
import { getSession, useSession } from "~/providers/session";
import { Loading } from "./ui/Loading";

const menuLinkVariants = cva("ui-disabled:pointer-events-none hover:bg-primary hover:text-primary-foreground ui-disabled:opacity-50 relative flex cursor-pointer select-none items-center gap-1 rounded-sm px-2 py-1.5 text-sm outline-none transition-colors", {
  variants: {
    size: {
      default: "h-10 px-4 py-2",
      icon: "h-10 w-10",
    },
    variant: {
      active: "bg-primary text-primary-foreground",
    }
  }
})

// FIXME: dropdown menu item <A> links are broken on IOS
function DropdownMenuLinks() {
  return (
    <>
      <DropdownMenu.Item asChild>
        <As component={A} class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
          href="/" end><RiBuildingsHomeLine class="h-5 w-5" />Home</As>
      </DropdownMenu.Item>
      <DropdownMenu.Item asChild>
        <As component={A} class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
          href="/view"><RiSystemEyeLine class="h-5 w-5" />View</As>
      </DropdownMenu.Item>
    </>
  )
}

function MenuLinks() {
  return (
    <div class="flex flex-col p-2">
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
        href="/" noScroll end><RiBuildingsHomeLine class="h-5 w-5" />Home</A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants()}
        href="/view" noScroll><RiSystemEyeLine class="h-5 w-5" />View</A>
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
  }).catch(toastError)
)

function Header(props: { onMenuClick: () => void }) {
  const signOutSubmission = useSubmission(actionSignOut)
  const signOutAction = useAction(actionSignOut)
  const signOut = () => signOutAction().catch(toastError)

  return (
    <div
      class="bg-background text-foreground border-b-border z-10 h-12 w-full overflow-x-hidden border-b">
      <div
        class="flex h-full items-center gap-2 px-2"
      >
        <DropdownMenuRoot>
          <DropdownMenuTrigger asChild>
            <As component={Button} size='icon' variant='ghost' title="Menu" class="md:hidden">
              <RiSystemMenuLine class="h-6 w-6" />
            </As>
          </DropdownMenuTrigger>
          <DropdownMenuPortal>
            <DropdownMenuContent>
              <DropdownMenuArrow />
              <DropdownMenuLinks />
            </DropdownMenuContent>
          </DropdownMenuPortal>
        </DropdownMenuRoot>
        <Button onClick={props.onMenuClick} size='icon' variant='ghost' title="Menu" class="hidden md:inline-flex">
          <RiSystemMenuLine class="h-6 w-6" />
        </Button>
        <div class="flex flex-1 items-center truncate text-xl">
          IPCManView
        </div>
        <div>
        </div>
        <div class="flex">
          <A href="/debug" title="Debug" class={menuLinkVariants({ size: "icon" })} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants({ size: "icon" })} end>
            <RiDevelopmentBugLine class="h-6 w-6" title="Debug" />
          </A>
          <DropdownMenuRoot>
            <DropdownMenuTrigger asChild>
              <As component={Button} size='icon' variant='ghost' title="User">
                <RiUserFacesUserLine class="h-6 w-6" />
              </As>
            </DropdownMenuTrigger>
            <DropdownMenuPortal>
              <DropdownMenuContent class="z-[200]">
                <DropdownMenuArrow />
                <DropdownMenu.Item asChild>
                  <As component={A} class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })}
                    href="/profile" end>
                    <RiUserFacesUserLine class="h-5 w-5" />Profile
                  </As>
                </DropdownMenu.Item>
                <DropdownMenu.Item class={menuLinkVariants()} onClick={signOut} disabled={signOutSubmission.pending}>
                  <RiSystemLogoutBoxRFill class="h-5 w-5" />Sign out
                </DropdownMenu.Item>
              </DropdownMenuContent>
            </DropdownMenuPortal>
          </DropdownMenuRoot>
          <Button size='icon' variant='ghost' onClick={toggleTheme} title={useThemeTitle()}>
            <ThemeIcon class="h-6 w-6" />
          </Button>
        </div>
      </div>
    </div >
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

export function Layout(props: ParentProps) {
  useSession()
  const [menuOpen, setMenuOpen] = makePersisted(createSignal(true), { "name": "menu-open" })
  const session = createAsync(getSession)
  const toastClass = () => session() ? "top-12 sm:top-12" : ""

  return (
    <>
      <Portal>
        <ToastRegion class={toastClass()}>
          <ToastList class={toastClass()} />
        </ToastRegion>
      </Portal>
      <Suspense fallback={<Loading class="pt-10" />}>
        <Show when={session()} fallback={<>{props.children}</>}>
          <Header onMenuClick={() => setMenuOpen((prev) => !prev)} />
          <div class="flex">
            <Menu menuOpen={menuOpen()}>
              <MenuLinks />
            </Menu>
            <div class="w-full overflow-x-auto"> {/* FIXME: overflow-x-auto is needed to fix overflowing tables BUT it also breaks something and I forgot what it was ¯\_(ツ)_/¯ */}
              {props.children}
            </div>
          </div >
        </Show>
      </Suspense>
    </>
  )
}
