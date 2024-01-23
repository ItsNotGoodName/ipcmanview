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

// const Progress = React.forwardRef<
//   React.ElementRef<typeof ProgressPrimitive.Root>,
//   React.ComponentPropsWithoutRef<typeof ProgressPrimitive.Root>
// >(({ className, value, ...props }, ref) => (
//   <ProgressPrimitive.Root
//     ref={ref}
//     className={cn(
//       "",
//       className
//     )}
//     {...props}
//   >
//     <ProgressPrimitive.Indicator
//       className="h-full w-full flex-1 bg-primary transition-all"
//       style={{ transform: `translateX(-${100 - (value || 0)}%)` }}
//     />
//   </ProgressPrimitive.Root>
// ))
// Progress.displayName = ProgressPrimitive.Root.displayName
//
// export { Progress }
