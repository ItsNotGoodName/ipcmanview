import { Checkbox } from "@kobalte/core";
import { RiSystemCheckLine } from "solid-icons/ri";
import { ComponentProps, splitProps } from "solid-js";

import { cn } from "~/lib/utils"

export const CheckboxRoot = Checkbox.Root
export const CheckboxInput = Checkbox.Input

export function CheckboxControl(props: ComponentProps<typeof Checkbox.Control>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Checkbox.Control
    class={cn(
      "border-primary ring-offset-background focus-visible:ring-ring ui-checked:bg-primary ui-checked:text-primary-foreground peer h-4 w-4 shrink-0 rounded-sm border focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
      props.class
    )}
    {...rest}
  />
}

export function CheckboxIndicator(props: Omit<ComponentProps<typeof Checkbox.Indicator>, "class">) {
  return <Checkbox.Indicator
    class={cn("flex items-center justify-center text-current")}
    {...props}
  />
}

export const CheckboxIcon = RiSystemCheckLine
export const CheckboxLabel = Checkbox.Label
export const CheckboxDescription = Checkbox.Description
export const CheckboxErrorMessage = Checkbox.ErrorMessage
