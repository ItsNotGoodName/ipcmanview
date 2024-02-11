// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/dropdown-menu
// https://ui.shadcn.com/docs/components/dropdown-menu
import { DropdownMenu } from "@kobalte/core";
import { ComponentProps, JSX, splitProps } from "solid-js";
import { RiArrowsArrowRightSLine, RiSystemCheckLine, RiSystemCheckboxBlankCircleFill } from "solid-icons/ri"

import { cn } from "~/lib/utils"

export const DropdownMenuRoot = DropdownMenu.Root
export const DropdownMenuTrigger = DropdownMenu.Trigger
export const DropdownMenuIcon = DropdownMenu.Icon
export const DropdownMenuPortal = DropdownMenu.Portal
export const DropdownMenuArrow = DropdownMenu.Arrow
export const DropdownMenuGroup = DropdownMenu.Group
export const DropdownMenuSub = DropdownMenu.Sub
export const DropdownMenuItemDescription = DropdownMenu.ItemDescription
export const DropdownMenuRadioGroup = DropdownMenu.RadioGroup

export function DropdownMenuSubTrigger(props: ComponentProps<typeof DropdownMenu.SubTrigger> & { inset?: boolean }) {
  const [_, rest] = splitProps(props, ["class", "children", "inset"])
  return <DropdownMenu.SubTrigger
    class={cn(
      "focus:bg-accent ui-expanded:bg-accent flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none",
      props.inset && "pl-8",
      props.class
    )}
    {...rest}
  >
    {props.children}
  </DropdownMenu.SubTrigger>
}

export function DropdownMenuSubTriggerIndicator() {
  return <RiArrowsArrowRightSLine class="ml-auto h-4 w-4" />
}

export function DropdownMenuSubContent(props: Omit<DropdownMenu.DropdownMenuSubContentProps, "style">) {
  const [_, rest] = splitProps(props, ["class"])
  return <DropdownMenu.SubContent
    class={cn(
      "bg-popover text-popover-foreground ui-expanded:animate-in ui-not-expanded:animate-out ui-not-expanded:fade-out-0 ui-expanded:fade-in-0 ui-not-expanded:zoom-out-95 ui-expanded:zoom-in-95 z-50 min-w-[8rem] max-w-[var(--kb-popper-content-available-width)] origin-[var(--kb-menu-content-transform-origin)] overflow-hidden rounded-md border p-1 shadow-lg",
      props.class
    )}
    {...rest}
  />
}

export function DropdownMenuContent(props: Omit<DropdownMenu.DropdownMenuContentProps, "style">) {
  const [_, rest] = splitProps(props, ["class"])
  return <DropdownMenu.Content
    class={cn(
      "bg-popover text-popover-foreground ui-expanded:animate-in ui-not-expanded:animate-out ui-not-expanded:fade-out-0 ui-expanded:fade-in-0 ui-not-expanded:zoom-out-95 ui-expanded:zoom-in-95 z-50 min-w-[8rem] max-w-[var(--kb-popper-content-available-width)] origin-[var(--kb-menu-content-transform-origin)] rounded-md border p-1 shadow-md",
      props.class,
    )}
    {...rest}
  />
}

export function DropdownMenuItem(props: ComponentProps<typeof DropdownMenu.Item> & { inset?: boolean }) {
  const [_, rest] = splitProps(props, ["class", "inset"])
  return <DropdownMenu.Item
    class={cn(
      "focus:bg-accent focus:text-accent-foreground ui-disabled:pointer-events-none ui-disabled:opacity-50 relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors",
      props.inset && "pl-8",
      props.class
    )}
    {...rest}
  />
}

export function DropdownMenuCheckboxItem(props: ComponentProps<typeof DropdownMenu.CheckboxItem> & { inset?: boolean }) {
  const [_, rest] = splitProps(props, ["class", "children", "checked"])
  return <DropdownMenu.CheckboxItem
    class={cn(
      "focus:bg-accent focus:text-accent-foreground ui-disabled:pointer-events-none ui-disabled:opacity-50 relative flex cursor-default select-none items-center rounded-sm py-1.5 pl-8 pr-2 text-sm outline-none transition-colors",
      props.class
    )}
    checked={props.checked}
    {...rest}
  >
    {props.children}
  </DropdownMenu.CheckboxItem>
}

export function DropdownMenuCheckboxItemIndicator() {
  return <span class="absolute left-2 flex h-3.5 w-3.5 items-center justify-center">
    <DropdownMenu.ItemIndicator>
      <RiSystemCheckLine class="h-4 w-4" />
    </DropdownMenu.ItemIndicator>
  </span>
}

export function DropdownMenuRadioItem(props: ComponentProps<typeof DropdownMenu.RadioItem> & { inset?: boolean }) {
  const [_, rest] = splitProps(props, ["class", "children"])
  return <DropdownMenu.RadioItem
    class={cn(
      "focus:bg-accent focus:text-accent-foreground ui-disabled:pointer-events-none ui-disabled:opacity-50 relative flex cursor-default select-none items-center rounded-sm py-1.5 pl-8 pr-2 text-sm outline-none transition-colors",
      props.class
    )}
    {...rest}
  >
    {props.children}
  </DropdownMenu.RadioItem>
}

export function DropdownMenuRadioItemIndicator() {
  return <span class="absolute left-2 flex h-3.5 w-3.5 items-center justify-center">
    <DropdownMenu.ItemIndicator>
      <RiSystemCheckboxBlankCircleFill class="h-2 w-2 fill-current" />
    </DropdownMenu.ItemIndicator>
  </span>
}

export function DropdownMenuGroupLabel(props: ComponentProps<typeof DropdownMenu.GroupLabel> & { inset?: boolean }) {
  const [_, rest] = splitProps(props, ["class", "inset"])
  return <DropdownMenu.GroupLabel
    class={cn(
      "px-2 py-1.5 text-sm font-semibold",
      props.inset && "pl-8",
      props.class
    )}
    {...rest}
  />
}

export function DropdownMenuSeparator(props: ComponentProps<typeof DropdownMenu.Separator>) {
  const [_, rest] = splitProps(props, ["class"])
  return <DropdownMenu.Separator
    class={cn("bg-muted -mx-1 my-1 h-px", props.class)}
    {...rest}
  />
}

export function DropdownMenuShortcut(props: JSX.HTMLAttributes<HTMLSpanElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <span
    class={cn("ml-auto pl-2 text-xs tracking-widest opacity-60", props.class)}
    {...rest}
  />
}

