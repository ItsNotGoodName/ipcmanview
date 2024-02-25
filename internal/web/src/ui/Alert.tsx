// # Changes
// - Destructive has a red background
//
// # URLs
// https://kobalte.dev/docs/core/components/alert
// https://ui.shadcn.com/docs/components/alert
import { cva, type VariantProps } from "class-variance-authority"
import { Alert } from "@kobalte/core"
import { ComponentProps, JSX, splitProps } from "solid-js"

import { cn } from "~/lib/utils"

const alertVariants = cva(
  "relative w-full rounded-lg border p-4 [&>svg+div]:translate-y-[-3px] [&>svg]:absolute [&>svg]:left-4 [&>svg]:top-4 [&>svg~*]:pl-7",
  {
    variants: {
      variant: {
        default: "bg-background text-foreground [&>svg]:text-foreground",
        destructive:
          "border-destructive/50 bg-destructive text-destructive-foreground dark:border-destructive [&>svg]:text-destructive-foreground",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
)

export function AlertRoot(props: ComponentProps<typeof Alert.Root> & VariantProps<typeof alertVariants>) {
  const [_, rest] = splitProps(props, ["class", "variant"])
  return <Alert.Root
    class={cn(alertVariants({ variant: props.variant }), props.class)}
    {...rest}
  />
}

export function AlertTitle(props: JSX.HTMLAttributes<HTMLHeadingElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <h5
    class={cn("mb-1 font-medium leading-none tracking-tight", props.class)}
    {...rest}
  />
}

export function AlertDescription(props: JSX.HTMLAttributes<HTMLParagraphElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <div
    class={cn("text-sm [&_p]:leading-relaxed", props.class)}
    {...rest}
  />
}

