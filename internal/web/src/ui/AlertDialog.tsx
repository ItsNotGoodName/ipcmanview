import { ComponentProps, JSX, splitProps } from "solid-js"
import { AlertDialog } from "@kobalte/core"

import { cn } from "~/lib/utils"
import { Button, buttonVariants } from "~/ui/Button"

export const AlertDialogRoot = AlertDialog.Root

export const AlertDialogTrigger = AlertDialog.Trigger

const AlertDialogPortal = AlertDialog.Portal

export function AlertDialogOverlay(props: ComponentProps<typeof AlertDialog.Overlay>) {
  const [_, rest] = splitProps(props, ["class"])
  return <AlertDialog.Overlay
    class={cn(
      "ui-expanded:animate-in ui-not-expanded:animate-out ui-not-expanded:fade-out-0 ui-expanded:fade-in-0 fixed inset-0 z-50 bg-black/80",
      props.class
    )}
    {...rest}
  />
}

export function AlertDialogModal(props: ComponentProps<typeof AlertDialog.Content>) {
  const [_, rest] = splitProps(props, ["class"])
  return <AlertDialogPortal>
    <AlertDialogOverlay />
    <AlertDialog.Content
      class={cn(
        "bg-background ui-expanded:animate-in ui-not-expanded:animate-out ui-not-expanded:fade-out-0 ui-expanded:fade-in-0 ui-not-expanded:zoom-out-95 ui-expanded:zoom-in-95 ui-not-expanded:slide-out-to-left-1/2 ui-not-expanded:slide-out-to-top-[48%] ui-expanded:slide-in-from-left-1/2 ui-expanded:slide-in-from-top-[48%] fixed left-[50%] top-[50%] z-50 flex max-h-screen w-full max-w-lg translate-x-[-50%] translate-y-[-50%] flex-col gap-4 border p-4 shadow-lg duration-200 sm:rounded-lg",
        props.class
      )}
      {...rest}
    />
  </AlertDialogPortal>
}

export function AlertDialogHeader(props: JSX.HTMLAttributes<HTMLDivElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <div
    class={cn(
      "flex flex-col space-y-2 overflow-y-hidden px-2 text-center sm:text-left",
      props.class
    )}
    {...rest}
  />
}

export function AlertDialogFooter(props: JSX.HTMLAttributes<HTMLDivElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <div
    class={cn(
      "flex flex-col-reverse px-2 sm:flex-row sm:justify-end sm:space-x-2",
      props.class
    )}
    {...rest}
  />
}

export function AlertDialogTitle(props: ComponentProps<typeof AlertDialog.Title>) {
  const [_, rest] = splitProps(props, ["class"])
  return <AlertDialog.Title
    class={cn("text-lg font-semibold", props.class)}
    {...rest}
  />
}

export function AlertDialogDescription(props: ComponentProps<typeof AlertDialog.Description>) {
  const [_, rest] = splitProps(props, ["class"])
  return <AlertDialog.Description
    class={cn("text-muted-foreground overflow-y-auto text-sm", props.class)}
    {...rest}
  />
}

export const AlertDialogAction = Button

export function AlertDialogCancel(props: JSX.ButtonHTMLAttributes<HTMLButtonElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <AlertDialog.CloseButton
    class={cn(
      buttonVariants({ variant: "outline" }),
      "mt-2 sm:mt-0",
      props.class
    )}
    {...rest}
  />
}
