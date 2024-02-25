// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/progress
// https://ui.shadcn.com/docs/components/progress
import { Progress } from "@kobalte/core"
import { ComponentProps, splitProps } from "solid-js"

import { cn } from "~/lib/utils"


export const ProgressRoot = Progress.Root;
export const ProgressLabel = Progress.Label;
export const ProgressValueLabel = Progress.ValueLabel;

export function ProgressTrack(props: ComponentProps<typeof Progress.Track>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Progress.Track class={cn("bg-secondary relative h-4 w-full overflow-hidden rounded-full", props.class)} {...rest} />
}

export function ProgressFill(props: ComponentProps<typeof Progress.Fill>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Progress.Fill
    class={cn("bg-primary h-full w-[var(--kb-progress-fill-width)] flex-1 transition-all", props.class)} {...rest} />
}
