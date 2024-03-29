// # Changes
// N/A
//
// # URLs
// https://kobalte.dev/docs/core/components/breadcrumbs
import { Breadcrumbs } from "@kobalte/core";
import { ParentProps, splitProps } from "solid-js";

import { cn } from "~/lib/utils";

export function BreadcrumbsRoot(props: Breadcrumbs.BreadcrumbsRootProps) {
  const [_, rest] = splitProps(props, ["children"])
  return <Breadcrumbs.Root {...rest}>
    <ol class="inline-flex items-center" >
      {props.children}
    </ol>
  </Breadcrumbs.Root>
}

export function BreadcrumbsItem(props: ParentProps) {
  return <li class="inline-flex items-center" >
    {props.children}
  </li>
}

export function BreadcrumbsLink(props: Breadcrumbs.BreadcrumbsLinkProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <Breadcrumbs.Link class={cn("hover:text-sky-500", props.class)} {...rest} />
}

export function BreadcrumbsSeparator(props: Breadcrumbs.BreadcrumbsSeparatorProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <Breadcrumbs.Separator class={cn("px-2", props.class)} {...rest} />
}

