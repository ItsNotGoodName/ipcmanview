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
import { labelVariants } from "./Label"

export function SelectHTML(props: JSX.SelectHTMLAttributes<HTMLSelectElement>) {
  const [_, rest] = splitProps(props, ["class"])
  // TODO: the dropdown arrow should be opacity-50 but the text should be normal
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

export function SelectLabel(props: Select.SelectLabelProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <Select.Label
    class={cn(labelVariants(), props.class)}
    {...rest}
  />
}

export function SelectDescription(props: Select.SelectDescriptionProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <Select.Description
    class={cn("text-muted-foreground text-sm", props.class)}
    {...rest}
  />
}

export function SelectErrorMessage(props: ComponentProps<typeof Select.ErrorMessage>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Select.ErrorMessage
    class={cn("text-destructive text-sm font-medium")}
    {...rest}
  />
}

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

export function SelectContent(props: Select.SelectContentProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <Select.Portal>
    <Select.Content
      class={cn(
        "bg-popover text-popover-foreground ui-expanded:animate-in ui-not-expanded:animate-out ui-not-expanded:fade-out-0 ui-expanded:fade-in-0 ui-not-expanded:zoom-out-95 ui-expanded:zoom-in-95 relative z-50 min-w-[8rem] max-w-[var(--kb-popper-content-available-width)] origin-[var(--kb-select-content-transform-origin)] overflow-hidden rounded-md border shadow-md",
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

export function SelectItem(props: Omit<Select.SelectItemProps, "class">) {
  const [_, rest] = splitProps(props, ["children"])
  return <Select.Item
    class="focus:bg-accent focus:text-accent-foreground ui-disabled:pointer-events-none ui-disabled:opacity-50 relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-8 pr-2 text-sm outline-none"
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
    class={cn("bg-muted -mx-1 my-1 h-px", props.class)}
    {...rest}
  />
}
