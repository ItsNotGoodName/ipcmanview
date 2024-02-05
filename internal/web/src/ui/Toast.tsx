// # Changes
// - Region is top right
//
// # URLs
// https://kobalte.dev/docs/core/components/toast
// https://ui.shadcn.com/docs/components/toast
import { Toast, toaster } from "@kobalte/core";
import { cva, type VariantProps } from "class-variance-authority"
import { RiSystemCloseLine } from "solid-icons/ri";
import { ComponentProps, JSX, splitProps } from "solid-js";

import { cn } from "~/lib/utils"

export function ToastRegion(props: ComponentProps<typeof Toast.Region>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Toast.Region
    class={cn(
      "fixed top-0 z-[100] flex max-h-screen w-full flex-col-reverse p-4 sm:bottom-auto sm:right-0 sm:top-0 sm:flex-col md:max-w-[420px]",
      props.class
    )}
    {...rest}
  />
}

export function ToastList(props: ComponentProps<typeof Toast.List>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Toast.List
    class={cn(
      "left-0 right-0 top-0 z-50 m-0 flex max-h-screen w-full flex-col-reverse gap-2 p-4 sm:left-auto md:max-w-md",
      props.class
    )}
    {...rest}
  />
}

const toastVariants = cva(
  "ui-opened:animate-in ui-closed:animate-out data-[swipe=end]:animate-out ui-closed:fade-out-80 ui-closed:slide-out-to-right-full ui-opened:slide-in-from-top-full ui-opened:sm:slide-in-from-top-full group pointer-events-auto relative flex w-full items-center justify-between space-x-4 overflow-hidden rounded-md border p-6 pr-8 shadow-lg transition-all data-[swipe=cancel]:translate-x-0 data-[swipe=end]:translate-x-[var(--kb-toast-swipe-end-x)] data-[swipe=move]:translate-x-[var(--kb-toast-swipe-move-x)] data-[swipe=move]:transition-none",
  {
    variants: {
      variant: {
        default: "bg-background text-foreground border",
        destructive:
          "destructive border-destructive bg-destructive text-destructive-foreground group",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
)

export function ToastRoot(props: ComponentProps<typeof Toast.Root> & VariantProps<typeof toastVariants>) {
  const [_, rest] = splitProps(props, ["class", "variant"])
  return <Toast.Root
    class={cn(toastVariants({ variant: props.variant }), props.class)}
    {...rest}
  />
}

export function ToastContent(props: JSX.HTMLAttributes<HTMLDivElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <div
    class={cn("flex w-full flex-col gap-2")}
    {...rest}
  />
}

export function ToastCloseButton(props: ComponentProps<typeof Toast.CloseButton>) {
  const [_, rest] = splitProps(props, ["class"])
  return (
    <Toast.CloseButton
      title="Close"
      class={cn(
        "text-foreground/50 hover:text-foreground absolute right-2 top-2 rounded-md p-1 opacity-0 transition-opacity focus:opacity-100 focus:outline-none focus:ring-2 group-hover:opacity-100 group-[.destructive]:text-red-300 group-[.destructive]:hover:text-red-50 group-[.destructive]:focus:ring-red-400 group-[.destructive]:focus:ring-offset-red-600",
        props.class
      )}

      {...rest}>
      <RiSystemCloseLine class="h-4 w-4" />
    </Toast.CloseButton>
  )
}

export function ToastTitle(props: ComponentProps<typeof Toast.Title>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Toast.Title
    class={cn("text-sm font-semibold", props.class)}
    {...rest}
  />
}

export function ToastDescription(props: ComponentProps<typeof Toast.Description>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Toast.Description
    class={cn("text-sm opacity-90", props.class)}
    {...rest}
  />
}

export function ToastProgressTrack(props: ComponentProps<typeof Toast.ProgressTrack>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Toast.ProgressTrack
    class={cn("bg-primary-foreground h-2 w-full rounded", props.class)}
    {...rest}
  />
}

export function ToastProgressFill(props: ComponentProps<typeof Toast.ProgressFill>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Toast.ProgressFill
    style={{ width: "var(--kb-toast-progress-fill-width)" }}
    class={cn("bg-primary h-full rounded transition-all", props.class)}
    {...rest}
  />
}

function show(message: string) {
  return toaster.show(props => (
    <ToastRoot toastId={props.toastId}>
      <ToastContent>
        <ToastCloseButton />
        {message}
      </ToastContent>
    </ToastRoot>
  ));
}

function success(message: string) {
  return toaster.show(props => (
    <ToastRoot toastId={props.toastId}>
      <ToastContent>
        <ToastCloseButton />
        {message}
      </ToastContent>
    </ToastRoot>
  ));
}

function error(title: string, message: string): number {
  return toaster.show(props => (
    <ToastRoot toastId={props.toastId} variant="destructive">
      <ToastContent>
        <ToastCloseButton />
        <ToastTitle>{title}</ToastTitle>
        <ToastDescription>
          {message}
        </ToastDescription>
      </ToastContent>
    </ToastRoot>
  ));
}

function custom(ele: () => JSX.Element, rootProps?: Omit<ComponentProps<typeof ToastRoot>, "toastId">) {
  return toaster.show(props => <ToastRoot toastId={props.toastId} {...rootProps}>{ele as any}</ToastRoot>);
}

function dismiss(id: number) {
  return toaster.dismiss(id);
}

export const toast = {
  show,
  success,
  error,
  custom,
  dismiss,
};
