// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/checkbox
// https://ui.shadcn.com/docs/components/checkbox
import { Checkbox } from "@kobalte/core";
import { RiSystemCheckLine } from "solid-icons/ri";
import { ComponentProps, splitProps } from "solid-js";

import { cn } from "~/lib/utils"
import { labelVariants } from "./Label";

export const CheckboxRoot = Checkbox.Root

export function CheckboxControl(props: Omit<Checkbox.CheckboxControlProps, "children">) {
  const [_, rest] = splitProps(props, ["class"])
  return <>
    <Checkbox.Input class="peer" />
    <Checkbox.Control
      class={cn(
        "border-primary peer-focus-visible:ring-ring ui-checked:bg-primary ui-checked:text-primary-foreground ui-disabled:cursor-not-allowed ui-disabled:opacity-50 peer h-4 w-4 shrink-0 cursor-pointer rounded-sm border shadow peer-focus-visible:outline-none peer-focus-visible:ring-1",
        props.class
      )}
      {...rest}
    >
      <Checkbox.Indicator class="flex items-center justify-center text-current">
        <RiSystemCheckLine class="h-4 w-4" />
      </Checkbox.Indicator>
    </Checkbox.Control>
  </>
}

export function CheckboxLabel(props: Checkbox.CheckboxLabelProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <Checkbox.Label
    class={cn(labelVariants(), props.class)}
    {...rest}
  />
}

export function CheckboxDescription(props: ComponentProps<typeof Checkbox.Description>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Checkbox.Description
    class={cn("text-muted-foreground text-sm")}
    {...rest}
  />
}

export function CheckboxErrorMessage(props: ComponentProps<typeof Checkbox.ErrorMessage>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Checkbox.ErrorMessage
    class={cn("text-destructive text-sm font-medium")}
    {...rest}
  />
}
