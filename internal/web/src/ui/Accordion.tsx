import { Accordion } from "@kobalte/core"
import { RiArrowsArrowDownSFill, RiArrowsArrowDownSLine } from "solid-icons/ri"
import { splitProps } from "solid-js"

import { cn } from "~/lib/utils"

export const AccordionRoot = Accordion.Root

export function AccordionItem(props: Accordion.AccordionItemProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <Accordion.Item
    class={cn("border-b", props.class)}
    {...rest}
  />
}

export function AccordionTrigger(props: Accordion.AccordionTriggerProps) {
  const [_, rest] = splitProps(props, ["class", "children"])
  return <Accordion.Header class="flex">
    <Accordion.Trigger
      class={cn(
        "flex flex-1 items-center justify-between py-4 font-medium transition-all hover:underline [&[data-expanded]>svg]:rotate-180",
        props.class
      )}
      {...rest}
    >
      {props.children}
      <RiArrowsArrowDownSLine class="h-4 w-4 shrink-0 transition-transform duration-200" />
    </Accordion.Trigger>
  </Accordion.Header>
}


export function AccordionContent(props: Accordion.AccordionContentProps) {
  const [_, rest] = splitProps(props, ["class", "children"])
  return <Accordion.Content
    class="overflow-hidden text-sm transition-all ui-not-expanded:animate-accordion-up ui-expanded:animate-accordion-down"
    {...rest}
  >
    <div class={cn("pb-4 pt-0", props.class)}>{props.children}</div>
  </Accordion.Content>
}
