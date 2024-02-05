// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/popover
// https://ui.shadcn.com/docs/components/popover
import { Popover } from "@kobalte/core";
import { RiSystemCloseLine } from "solid-icons/ri";
import { splitProps } from "solid-js";

import { cn } from "~/lib/utils";

export const PopoverRoot = Popover.Root
export const PopoverTrigger = Popover.Trigger
export const PopoverAnchor = Popover.Anchor
export const PopoverPortal = Popover.Portal

export function PopoverContent(props: Omit<Popover.PopoverContentProps, "style">) {
  const [_, rest] = splitProps(props, ["class"])
  return <Popover.Content
    style={{ "max-width": "var(--kb-popper-content-available-width)", "transform-origin": "var(--kb-menu-content-transform-origin)" }}
    class={cn(
      "bg-popover text-popover-foreground ui-expanded:animate-in ui-not-expanded:animate-out ui-not-expanded:fade-out-0 ui-expanded:fade-in-0 ui-not-expanded:zoom-out-95 ui-expanded:zoom-in-95 z-50 w-72 rounded-md border p-4 shadow-md outline-none",
      props.class
    )}
    {...rest}
  />
}

export const PopoverArrow = Popover.Arrow
export const PopoverCloseButton = Popover.CloseButton
export const PopoverCloseIcon = RiSystemCloseLine
export const PopoverTitle = Popover.Title
export const PopoverDescription = Popover.Description
