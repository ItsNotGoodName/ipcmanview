// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/text-field
// https://ui.shadcn.com/docs/components/input
// https://ui.shadcn.com/docs/components/textarea
import { TextField } from "@kobalte/core"
import { ComponentProps, JSX, splitProps } from "solid-js"

import { cn } from "~/lib/utils"
import { labelVariants } from "./Label"

export function InputHTML(props: JSX.InputHTMLAttributes<HTMLInputElement>) {
  const [_, rest] = splitProps(props, ["class", "type"])
  return (
    <input
      type={props.type ?? "text"}
      class={cn(
        "border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring flex h-10 w-full rounded-md border px-3 py-2 text-sm file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
        props.class
      )}
      {...rest}
    />
  )
}

export function TextFieldInput(props: TextField.TextFieldInputProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <TextField.Input
    class={cn(
      "border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring flex h-10 w-full rounded-md border px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
      props.class
    )}
    {...rest}
  />
}

export function TextareaHTML(props: JSX.TextareaHTMLAttributes<HTMLTextAreaElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <textarea
    class={cn(
      props.class
    )}
    {...rest}
  />
}

export function TextFieldTextArea(props: TextField.TextFieldTextAreaProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <TextField.TextArea
    class={cn(
      "border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring flex min-h-[80px] w-full rounded-md border px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
      props.class
    )}
    {...rest}
  />
}

export const TextFieldRoot = TextField.Root

export function TextFieldLabel(props: TextField.TextFieldLabelProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <TextField.Label
    class={cn(labelVariants(), props.class)}
    {...rest}
  />
}

export function TextFieldDescription(props: TextField.TextFieldDescriptionProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <TextField.Description
    class={cn("text-muted-foreground text-sm", props.class)}
    {...rest}
  />
}

export function TextFieldErrorMessage(props: ComponentProps<typeof TextField.ErrorMessage>) {
  const [_, rest] = splitProps(props, ["class"])
  return <TextField.ErrorMessage
    class={cn("text-destructive text-sm font-medium")}
    {...rest}
  />
}
