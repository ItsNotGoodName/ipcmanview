// # Changes
// N/A
//
// # URLs
// https://ui.shadcn.com/docs/components/card
import { JSX, splitProps } from "solid-js"

import { cn } from "~/lib/utils"

export function CardRoot(props: JSX.HTMLAttributes<HTMLDivElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <div
    class={cn("bg-card text-card-foreground rounded-lg border shadow-sm", props.class)}
    {...rest}
  />
}

export function CardHeader(props: JSX.HTMLAttributes<HTMLDivElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <div
    class={cn("flex flex-col space-y-1.5 p-6", props.class)}
    {...rest}
  />
}

export function CardTitle(props: JSX.HTMLAttributes<HTMLHeadingElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <h3
    class={cn("text-2xl font-semibold leading-none tracking-tight", props.class)}
    {...rest}
  />
}

export function CardDescription(props: JSX.HTMLAttributes<HTMLParagraphElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <p
    class={cn("text-muted-foreground text-sm", props.class)}
    {...rest}
  />
}

export function CardContent(props: JSX.HTMLAttributes<HTMLDivElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <div class={cn("p-6 pt-0", props.class)} {...rest} />
}

export function CardFooter(props: JSX.HTMLAttributes<HTMLDivElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <div
    class={cn("flex items-center p-6 pt-0", props.class)}
    {...rest}
  />
}
