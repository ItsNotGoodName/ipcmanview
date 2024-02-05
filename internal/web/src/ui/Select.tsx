// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/select
// https://ui.shadcn.com/docs/components/select
import { As, Select } from "@kobalte/core"
import { RiArrowsArrowDownSLine, RiSystemCheckLine } from "solid-icons/ri"
import { ComponentProps, JSX, splitProps } from "solid-js"

import { cn } from "~/lib/utils"

export function SelectHTML(props: JSX.SelectHTMLAttributes<HTMLSelectElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <select
    class={cn(
      "border-input bg-background ring-offset-background placeholder:text-muted-foreground focus:ring-ring flex h-10 w-full items-center justify-between rounded-md border px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
      props.class
    )}
    {...rest}
  >
    {props.children}
  </select>
}

export const SelectRoot = Select.Root

export const SelectValue = Select.Value

export const SelectHiddenSelect = Select.HiddenSelect

export const SelectDescription = Select.Description

export function SelectTrigger(props: ComponentProps<typeof Select.Trigger>) {
  const [_, rest] = splitProps(props, ["class", "children"])
  return <Select.Trigger
    class={cn(
      "border-input bg-background ring-offset-background placeholder:text-muted-foreground focus:ring-ring flex h-10 w-full items-center justify-between rounded-md border px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 [&>span]:line-clamp-1",
      props.class
    )}
    {...rest}
  >
    {props.children}
    <Select.Icon asChild>
      <As component={RiArrowsArrowDownSLine} class="h-4 w-4 opacity-50" />
    </Select.Icon>
  </Select.Trigger>
}

export function SelectContent(props: Omit<Select.SelectContentProps, "style">) {
  const [_, rest] = splitProps(props, ["class"])
  return <Select.Portal>
    <Select.Content
      style={{ "max-width": "var(--kb-popper-content-available-width)", "transform-origin": "var(--kb-menu-content-transform-origin)" }}
      class={cn(
        "bg-popover text-popover-foreground ui-expanded:animate-in ui-not-expanded:animate-out ui-not-expanded:fade-out-0 ui-expanded:fade-in-0 ui-not-expanded:zoom-out-95 ui-expanded:zoom-in-95 relative z-50 min-w-[8rem] overflow-hidden rounded-md border shadow-md",
        props.class
      )}
      {...rest}
    >
    </Select.Content>
  </Select.Portal>
}

export function SelectListbox() {
  return <Select.Listbox class="max-h-96 overflow-y-auto p-1" />
}

export function SelectLabel(props: ComponentProps<typeof Select.Label>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Select.Label
    class={cn("py-1.5 pl-8 pr-2 text-sm font-semibold", props.class)}
    {...rest}
  />
}

export function SelectItem(props: ComponentProps<typeof Select.Item>) {
  const [_, rest] = splitProps(props, ["class", "children"])
  return <Select.Item
    class={cn(
      "focus:bg-accent focus:text-accent-foreground ui-disabled:pointer-events-none ui-disabled:opacity-50 relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-8 pr-2 text-sm outline-none",
      props.class
    )}
    {...rest}
  >
    <span class="absolute left-2 flex h-3.5 w-3.5 items-center justify-center">
      <Select.ItemIndicator>
        <RiSystemCheckLine class="h-4 w-4" />
      </Select.ItemIndicator>
    </span>
    <Select.ItemLabel>{props.children}</Select.ItemLabel>
  </Select.Item>
}

export function SelectSeparator(props: JSX.HTMLAttributes<HTMLDivElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <div
    class={cn("-mx-1 my-1 h-px bg-muted", props.class)}
    {...rest}
  />
}
