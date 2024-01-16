import { cva } from "class-variance-authority"
import { As, DropdownMenu } from "@kobalte/core";
import { JSX, ParentProps, createEffect, createSignal, splitProps } from "solid-js";
import { Button } from "~/ui/Button";
import { DropdownMenuArrow, DropdownMenuContent, DropdownMenuPortal, DropdownMenuRoot, DropdownMenuTrigger } from "~/ui/DropdownMenu";
import { A } from "@solidjs/router";
import { RiBuildingsHomeLine, RiSystemEyeLine, RiSystemMenuLine } from "solid-icons/ri";
import { ThemeIcon } from "~/ui/ThemeIcon";
import { Theme, currentTheme, toggleTheme } from "~/ui/theme";
import { makePersisted } from "@solid-primitives/storage";
import { Portal } from "solid-js/web";
import { ToastList, ToastRegion } from "~/ui/Toast";

const menuLinkVariants = cva("ui-disabled:pointer-events-none ui-disabled:opacity-50 relative flex select-none items-center gap-1 rounded-sm px-2 py-1.5 text-sm outline-none transition-colors", {
  variants: {
    variant: {
      active: "bg-primary text-primary-foreground",
      inactive: "hover:bg-primary hover:text-primary-foreground"
    }
  }
})

function DropdownMenuLinks() {
  return (
    <>
      <DropdownMenu.Item asChild>
        <As component={A} class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants({ variant: "inactive" })}
          href="/" end><RiBuildingsHomeLine class="h-4 w-4" />Home</As>
      </DropdownMenu.Item>
      <DropdownMenu.Item asChild>
        <As component={A} class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants({ variant: "inactive" })}
          href="/view"><RiSystemEyeLine class="h-4 w-4" />View</As>
      </DropdownMenu.Item>
    </>
  )
}

function MenuLinks() {
  return (
    <div class="flex flex-col p-2">
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants({ variant: "inactive" })}
        href="/" noScroll end><RiBuildingsHomeLine class="h-4 w-4" />Home</A>
      <A class={menuLinkVariants()} activeClass={menuLinkVariants({ variant: "active" })} inactiveClass={menuLinkVariants({ variant: "inactive" })}
        href="/view" noScroll><RiSystemEyeLine class="h-4 w-4" />View</A>
    </div>
  )
}

function Header(props: { onMenuClick: () => void }) {
  const themeTitle = () => {
    switch (currentTheme()) {
      case Theme.System:
        return "System Theme"
      case Theme.Light:
        return "Light Theme"
      case Theme.Dark:
        return "Dark Theme"
    }
  }

  return (
    <div
      class="bg-background text-foreground border-b-border z-10 h-14 w-full border-b">
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
        <Button onClick={props.onMenuClick} size='icon' variant='ghost' title="Menu" class="hidden md:flex">
          <RiSystemMenuLine class="h-6 w-6" />
        </Button>
        <div class="flex flex-1 items-center truncate text-xl">
          IPCManView
        </div>
        <Button size='icon' variant='ghost' onClick={toggleTheme} title={themeTitle()}>
          <ThemeIcon class="h-6 w-6" />
        </Button>
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

export function Layout(props: ParentProps) {
  const [menuOpen, setMenuOpen] = makePersisted(createSignal(true), { "name": "menu-open" })
  return (
    <>
      <Portal>
        <ToastRegion>
          <ToastList />
        </ToastRegion>
      </Portal>
      <Header onMenuClick={() => setMenuOpen((prev) => !prev)} />
      <div class="flex">
        <Menu menuOpen={menuOpen()}>
          <MenuLinks />
        </Menu>
        <div class="flex-1">
          {props.children}
        </div>
      </div >
    </>
  )
}
