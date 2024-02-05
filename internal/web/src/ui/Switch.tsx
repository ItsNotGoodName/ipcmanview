// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/switch
// https://ui.shadcn.com/docs/components/switch
import { Switch } from "@kobalte/core"
import { ComponentProps, splitProps } from "solid-js"

import { cn } from "~/lib/utils"

export const SwitchRoot = Switch.Root
export const SwitchLabel = Switch.Label
export const SwitchDescription = Switch.Description
export const SwitchErrorMessage = Switch.ErrorMessage
export const SwitchInput = Switch.Input

export function SwitchControl(props: ComponentProps<typeof Switch.Control>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Switch.Control
    class={cn(
      "focus-visible:ring-ring focus-visible:ring-offset-background ui-checked:bg-primary ui-not-checked:bg-input peer inline-flex h-6 w-11 shrink-0 cursor-pointer items-center rounded-full border-2 border-transparent transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
      props.class
    )}
    {...rest}
  >
    <Switch.Thumb
      class={cn(
        "bg-background ui-checked:translate-x-5 ui-not-checked:translate-x-0 pointer-events-none block h-5 w-5 rounded-full shadow-lg ring-0 transition-transform"
      )}
    />
  </Switch.Control>
}

