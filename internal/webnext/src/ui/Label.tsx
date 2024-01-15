import { cva, type VariantProps } from "class-variance-authority"
import { JSX, splitProps } from "solid-js"

import { cn } from "~/lib/utils"

const labelVariants = cva(
  "text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
)

export function Label(props: JSX.LabelHTMLAttributes<HTMLLabelElement> & VariantProps<typeof labelVariants>) {
  const [_, rest] = splitProps(props, ["class"])
  return <label
    class={cn(labelVariants(), props.class)}
    {...rest}
  />
}
