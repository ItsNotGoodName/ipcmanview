import { Dialog } from "@kobalte/core"
import { RiSystemCloseLine } from "solid-icons/ri";
import { cva, type VariantProps } from "class-variance-authority"
import { ComponentProps, JSX, splitProps } from "solid-js";

import { cn } from "~/lib/utils"
import { DialogOverlay, DialogPortal } from "./Dialog";

export const SheetRoot = Dialog.Root
export const SheetTrigger = Dialog.Trigger
export const SheetPortal = Dialog.Portal
export const SheetCloseButton = Dialog.CloseButton

export function SheetOverlay(props: ComponentProps<typeof Dialog.Overlay>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Dialog.Overlay
    class={cn(
      "ui-expanded:animate-in ui-not-expanded:animate-out ui-not-expanded:fade-out-0 ui-expanded:fade-in-0 fixed inset-0 z-50 bg-black/80",
      props.class
    )}
    {...rest}
  />
}

const sheetVariants = cva(
  "bg-background ui-expanded:animate-in ui-not-expanded:animate-out ui-not-expanded:duration-300 ui-expanded:duration-500 fixed z-50 gap-4 p-6 shadow-lg transition ease-in-out",
  {
    variants: {
      side: {
        top: "ui-not-expanded:slide-out-to-top ui-expanded:slide-in-from-top inset-x-0 top-0 border-b",
        bottom:
          "ui-not-expanded:slide-out-to-bottom ui-expanded:slide-in-from-bottom inset-x-0 bottom-0 border-t",
        left: "ui-not-expanded:slide-out-to-left ui-expanded:slide-in-from-left inset-y-0 left-0 h-full w-3/4 border-r sm:max-w-sm",
        right:
          "ui-not-expanded:slide-out-to-right ui-expanded:slide-in-from-right inset-y-0 right-0  h-full w-3/4 border-l sm:max-w-sm",
      },
    },
    defaultVariants: {
      side: "right",
    },
  }
)

type SheetContentProps = ComponentProps<typeof Dialog.Content> & VariantProps<typeof sheetVariants>

export function SheetContent(props: SheetContentProps) {
  const [_, rest] = splitProps(props, ["class", "side", "children"])
  return <DialogPortal>
    <DialogOverlay />
    <Dialog.Content
      class={cn(sheetVariants({ side: props.side }), props.class)}
      {...rest}
    >
      {props.children}
      <Dialog.CloseButton class="ring-offset-background focus:ring-ring ui-expanded:bg-secondary absolute right-4 top-4 rounded-sm opacity-70 transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:pointer-events-none">
        <RiSystemCloseLine class="h-4 w-4" />
        <span class="sr-only">Close</span>
      </Dialog.CloseButton>
    </Dialog.Content>
  </DialogPortal>
}

export function SheetHeader(props: JSX.HTMLAttributes<HTMLDivElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <div
    class={cn(
      "flex flex-col space-y-1.5 text-center sm:text-left",
      props.class
    )}
    {...rest}
  />
}

export function SheetFooter(props: JSX.HTMLAttributes<HTMLDivElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <div
    class={cn(
      "flex flex-col-reverse sm:flex-row sm:justify-end sm:space-x-2",
      props.class
    )}
    {...rest}
  />
}

export function SheetTitle(props: ComponentProps<typeof Dialog.Title>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Dialog.Title
    class={cn(
      "text-foreground text-lg font-semibold",
      props.class
    )}
    {...rest}
  />
}

export function SheetDescription(props: ComponentProps<typeof Dialog.Description>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Dialog.Description
    class={cn("text-muted-foreground text-sm", props.class)}
    {...rest}
  />
}

