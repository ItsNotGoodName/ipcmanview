import { Image } from "@kobalte/core"
import { ComponentProps, splitProps } from "solid-js"

import { cn } from "~/lib/utils"

export function AvatarRoot(props: ComponentProps<typeof Image.Root>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Image.Root
    class={cn(
      "relative flex h-10 w-10 shrink-0 overflow-hidden rounded-full",
      props.class
    )}
    {...rest}
  />
}

export function AvatarImage(props: ComponentProps<typeof Image.Img>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Image.Img
    class={cn("aspect-square h-full w-full", props.class)}
    {...rest}
  />
}

export function AvatarFallback(props: ComponentProps<typeof Image.Fallback>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Image.Fallback
    class={cn(
      "bg-muted flex h-full w-full items-center justify-center rounded-full",
      props.class
    )}
    {...rest}
  />
}

