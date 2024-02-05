// # Changes
// N/A
//
// # URLs
// https://ui.shadcn.com/docs/components/label
import { cva, type VariantProps } from "class-variance-authority"
import { JSX, splitProps } from "solid-js"

import { cn } from "~/lib/utils"

export const labelVariants = cva(
  "text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
)

export type LabelProps = JSX.LabelHTMLAttributes<HTMLLabelElement> & VariantProps<typeof labelVariants>

export function Label(props: LabelProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <label
    class={cn(labelVariants(), props.class)}
    {...rest}
  />
}
