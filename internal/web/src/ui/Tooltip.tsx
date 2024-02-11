// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/tooltip
// https://ui.shadcn.com/docs/components/tooltip
import { splitProps } from "solid-js"
import { Tooltip } from "@kobalte/core"

import { cn } from "~/lib/utils"


export const TooltipRoot = Tooltip.Root
export const TooltipTrigger = Tooltip.Trigger
export const TooltipArrow = Tooltip.Arrow

export function TooltipContent(props: Tooltip.TooltipContentProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <Tooltip.Portal>
    <Tooltip.Content
      class={cn(
        "bg-popover text-popover-foreground animate-in fade-in-0 zoom-in-95 ui-not-expanded:animate-out ui-not-expanded:fade-out-0 ui-not-expanded:zoom-out-95 z-50 max-w-[var(--kb-popper-content-available-width)] origin-[var(--kb-tooltip-content-transform-origin)] rounded-md border px-3 py-1.5 text-sm shadow-md",
        props.class
      )}
      {...rest}
    />
  </Tooltip.Portal>
}
