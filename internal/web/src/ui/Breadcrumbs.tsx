import { Breadcrumbs } from "@kobalte/core";
import { ParentProps, splitProps } from "solid-js";
import { cn } from "~/lib/utils";


export const BreadcrumbsRoot = Breadcrumbs.Root

export function BreadcrumbsList(props: ParentProps) {
  return <ol class="inline-flex items-center" >
    {props.children}
  </ol>
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

