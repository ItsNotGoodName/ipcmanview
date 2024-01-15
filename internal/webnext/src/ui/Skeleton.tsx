import { JSX, splitProps } from "solid-js"
import { cn } from "~/lib/utils"

export function Skeleton(props: JSX.HTMLAttributes<HTMLDivElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <div
    class={cn("bg-muted animate-pulse rounded-md", props.class)}
    {...rest}
  />
}
