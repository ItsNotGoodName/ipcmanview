// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/hover-card
// https://ui.shadcn.com/docs/components/hover-card
import { splitProps } from "solid-js"
import { HoverCard } from "@kobalte/core"

import { cn } from "~/lib/utils"

export const HoverCardRoot = HoverCard.Root
export const HoverCardTrigger = HoverCard.Trigger
export const HoverCardArrow = HoverCard.Arrow

export function HoverCardContent(props: HoverCard.HoverCardContentProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <HoverCard.Portal>
    <HoverCard.Content
      class={cn(
        "bg-popover text-popover-foreground ui-expanded:animate-in ui-not-expanded:animate-out ui-not-expanded:fade-out-0 ui-expanded:fade-in-0 ui-not-expanded:zoom-out-95 ui-expanded:zoom-in-95 z-50 w-64 max-w-[var(--kb-popper-content-available-width)] origin-[var(--kb-hovercard-content-transform-origin)] rounded-md border p-4 shadow-md outline-none",
        props.class
      )}
      {...rest}
    />
  </HoverCard.Portal>
}
