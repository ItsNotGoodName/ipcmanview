import { RiSystemLoader4Line } from "solid-icons/ri"
import { JSX, splitProps } from "solid-js"

import { cn } from "~/lib/utils"

export function Loading(props: JSX.HTMLAttributes<HTMLDivElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return (
    <div class={cn("flex justify-center", props.class)} {...rest}>
      <div class="flex flex-col items-center gap-2">
        <RiSystemLoader4Line class="h-12 w-12 animate-spin" />
        <div class="text-xl">Loading...</div>
      </div>
    </div>
  )
}
