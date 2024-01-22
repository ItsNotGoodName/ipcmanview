import { RiSystemLoader4Line } from "solid-icons/ri"
import { ComponentProps, JSX, splitProps } from "solid-js"

import { cn } from "~/lib/utils"
import { AlertDescription, AlertRoot, AlertTitle } from "./Alert"
import { Alert } from "@kobalte/core"

export function PageError(props: ComponentProps<typeof Alert.Root> & { error: Error }) {
  const [_, rest] = splitProps(props, ["error"])
  return (
    <AlertRoot variant="destructive" {...rest}>
      <AlertTitle>Error</AlertTitle>
      <AlertDescription>{props.error.message}</AlertDescription>
    </AlertRoot>
  )
}

export function PageLoading(props: JSX.HTMLAttributes<HTMLDivElement>) {
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
