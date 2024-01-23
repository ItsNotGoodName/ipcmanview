import { Tabs } from "@kobalte/core"
import { ComponentProps, splitProps } from "solid-js"

import { cn } from "~/lib/utils"

export const TabsRoot = Tabs.Root

export function TabsList(props: ComponentProps<typeof Tabs.List>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Tabs.List
    class={cn(
      "bg-muted text-muted-foreground inline-flex h-10 items-center justify-center rounded-md p-1",
      props.class
    )}
    {...rest}
  />
}

export function TabsTrigger(props: ComponentProps<typeof Tabs.Trigger>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Tabs.Trigger
    class={cn(
      "ring-offset-background focus-visible:ring-ring ui-selected:bg-background ui-selected:text-foreground inline-flex items-center justify-center whitespace-nowrap rounded-sm px-3 py-1.5 text-sm font-medium transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 ui-selected::shadow-sm",
      props.class
    )}
    {...rest}
  />
}

export function TabsContent(props: ComponentProps<typeof Tabs.Content>) {
  const [_, rest] = splitProps(props, ["class"])
  return <Tabs.Content
    class={cn(
      "ring-offset-background focus-visible:ring-ring mt-2 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2",
      props.class
    )}
    {...rest}
  />
}

